package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"runtime"
	"sync"

	pkgerr "github.com/pkg/errors"
	"go.uber.org/zap"
)

type ReusableAcceptor interface {
	Accept(ctx context.Context, conn net.Conn)
}

type Server struct {
	acceptorPool sync.Pool
	wg           sync.WaitGroup
}

func New(allocator func() ReusableAcceptor) *Server {
	return &Server{
		acceptorPool: sync.Pool{
			New: func() any { return allocator() },
		},
		wg: sync.WaitGroup{},
	}
}

func (s *Server) Run(ctx context.Context, host string, port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return pkgerr.Wrap(err, "failed to start TCP listener")
	}

	defer listener.Close()

	go func() {
		<-ctx.Done()

		if err := listener.Close(); err != nil {
			zap.L().Error("Failed to close listener", zap.Error(err))
		}
	}()

	for {
		// Accept a new connection
		conn, err := listener.Accept()
		if err != nil {
			var opErr *net.OpError

			// Weird type to check error, but could not find an error code, and 'net' tests do the same
			if errors.As(err, &opErr) && opErr.Err.Error() == "use of closed network connection" {
				break
			}

			zap.L().Warn("Failed to accept client connection", zap.Error(err))

			return pkgerr.Wrap(err, "Server failed to accept client")
		}

		// It is assumed that Go routines are efficient enough to run one for each client.
		// However, if necessary for further optimization, a pre-allocated reusable pool of goroutines
		// may be used here.
		s.wg.Add(1)
		go s.handleConnection(ctx, conn)
	}

	s.wg.Wait() // Wait for graceful shutdown

	return nil
}

func (s *Server) handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	defer s.recoverFromPanic()
	defer s.wg.Done()

	acceptor := s.getAcceptor()

	defer s.putAcceptor(acceptor)

	zap.L().Info("New client accepted", zap.String("addr", conn.RemoteAddr().String()))
	acceptor.Accept(ctx, conn)
	zap.L().Info("Client disconnected", zap.String("addr", conn.RemoteAddr().String()))
}

func (s *Server) recoverFromPanic() {
	const maxBuffer = 2048

	if err := recover(); err != nil {
		buf := make([]byte, maxBuffer)
		n := runtime.Stack(buf, false)
		zap.L().Error("Server connection panic recovery", zap.Any("err", err), zap.ByteString("stack", buf[:n]))
	}
}

func (s *Server) getAcceptor() ReusableAcceptor { //nolint:ireturn // Concrete type is unknown here
	// A check for the maximum number of connections can be implemented to prevent resource
	// exhaustion during DDoS attacks.
	acceptor, ok := s.acceptorPool.Get().(ReusableAcceptor)
	if !ok {
		panic("unexpected object in client pool")
	}

	return acceptor
}

func (s *Server) putAcceptor(acceptor ReusableAcceptor) {
	s.acceptorPool.Put(acceptor)
}
