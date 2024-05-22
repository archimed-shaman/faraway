package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"reflect"
	"runtime"
	"sync"
	"time"

	pkgerr "github.com/pkg/errors"
	"go.uber.org/zap"
)

type Logic interface {
	Handle(ctx context.Context, req []byte) error
}

type Server struct {
	bufferPool  sync.Pool
	handlerPool sync.Pool
	timeout     time.Duration
	wg          sync.WaitGroup
}

func New(buffSize int, timeout time.Duration, allocLogic func() Logic) *Server {
	return &Server{
		bufferPool: sync.Pool{
			New: func() any { return make([]byte, buffSize) },
		},
		handlerPool: sync.Pool{
			New: func() any { return allocLogic() },
		},
		timeout: timeout,
		wg:      sync.WaitGroup{},
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
		go s.accept(ctx, conn)
	}

	s.wg.Wait() // Wait for graceful shutdown

	return nil
}

func (s *Server) accept(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	defer recoverFromPanic()
	defer s.wg.Done()

	buff := s.getBuffer()
	defer s.putBuffer(buff)

	handler := s.getHandler()
	defer s.putHandler(handler)

	zap.L().Info("New client connected", zap.String("addr", conn.RemoteAddr().String()))
	handleClient(ctx, conn, s.timeout, buff, handler.Handle)
	zap.L().Info("Client disconnected", zap.String("addr", conn.RemoteAddr().String()))
}

func (s *Server) getBuffer() []byte {
	obj := s.bufferPool.Get()

	buff, ok := obj.([]byte)
	if !ok {
		panic("unexpected object in acceptor pool: " + describeObj(obj))
	}

	return buff
}

func (s *Server) putBuffer(buff []byte) {
	s.bufferPool.Put(buff) //nolint:staticcheck // slice is the reference type
}

func (s *Server) getHandler() Logic { //nolint:ireturn // Concrete type is unknown here
	obj := s.handlerPool.Get()

	handler, ok := obj.(Logic)
	if !ok {
		panic("unexpected object in handler pool: " + describeObj(obj))
	}

	return handler
}

func (s *Server) putHandler(handler Logic) {
	s.handlerPool.Put(handler)
}

func recoverFromPanic() {
	const maxBuffer = 2048

	if err := recover(); err != nil {
		buf := make([]byte, maxBuffer)
		n := runtime.Stack(buf, false)
		zap.L().Error("Server connection panic recovery", zap.Any("err", err), zap.ByteString("stack", buf[:n]))
	}
}

func describeObj(obj any) string {
	if obj == nil {
		return "nil"
	}

	return reflect.TypeOf(obj).String()
}
