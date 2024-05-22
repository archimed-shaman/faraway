package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	pkgerr "github.com/pkg/errors"
	"go.uber.org/zap"
)

type (
	LogicAllocator func(askStop func(), send func(ctx context.Context, data []byte) error) Logic
	Logic          interface {
		Handle(ctx context.Context, req []byte) error
	}
)

type Server struct {
	maxConnections int32
	bufferPool     sync.Pool
	timeout        time.Duration
	allocLogic     LogicAllocator
	count          atomic.Int32
	wg             sync.WaitGroup
}

func New(buffSize int, maxConnections int32, timeout time.Duration, allocLogic LogicAllocator) *Server {
	return &Server{
		maxConnections: maxConnections,
		bufferPool:     sync.Pool{New: func() any { return make([]byte, buffSize) }},
		timeout:        timeout,
		allocLogic:     allocLogic,
		count:          atomic.Int32{},
		wg:             sync.WaitGroup{},
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

	s.count.Store(0)

	for {
		// Accept a new connection
		conn, err := listener.Accept()
		if err != nil {
			var opErr *net.OpError

			// Weird type to check error, but could not find an error code, and 'net' tests do the same
			if errors.As(err, &opErr) && opErr.Err.Error() == "use of closed network connection" {
				break
			}

			zap.L().Debug("Failed to accept client connection", zap.Error(err))

			return pkgerr.Wrap(err, "Server failed to accept client")
		}

		s.wg.Add(1)
		count := s.count.Add(1)

		// It is assumed that Go routines are efficient enough to run one for each client.
		// However, if necessary for further optimization, a pre-allocated reusable pool of goroutines
		// may be used here.
		go s.accept(ctx, conn, count)
	}

	s.wg.Wait() // Wait for graceful shutdown

	return nil
}

func (s *Server) accept(ctx context.Context, conn net.Conn, count int32) {
	defer conn.Close()
	defer s.wg.Done()
	defer recoverFromPanic()
	defer s.count.Add(-1)

	if count > s.maxConnections {
		// TODO: some verbose response can be send before disconnect
		return
	}

	buff := s.getBuffer()
	defer s.putBuffer(buff)

	clientCtx, cancel := context.WithCancel(ctx)
	logic := s.allocLogic(cancel, getSender(conn, s.timeout))

	zap.L().Info("New client connected", zap.String("addr", conn.RemoteAddr().String()))
	handleClient(clientCtx, conn, s.timeout, buff, logic.Handle)
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

func getSender(conn net.Conn, timeout time.Duration) func(ctx context.Context, data []byte) error {
	return func(ctx context.Context, data []byte) error {
		if err := conn.SetWriteDeadline(time.Now().Add(timeout)); err != nil {
			return pkgerr.Wrap(err, "failed to set write deadline for the connection")
		}

		// Assuming all data is sent on success. Anyway more complicated logic can be implemented
		if _, err := conn.Write(data); err != nil {
			return pkgerr.Wrap(err, "failed to write data")
		}

		return nil
	}
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
