package client

import (
	"context"
	"errors"
	"faraway/wow/app/infrastructure/server"
	"io"
	"os"

	pkgerr "github.com/pkg/errors"
	"go.uber.org/zap"
)

type Service struct {
	buff []byte
}

func NewService(buffSize int) *Service {
	return &Service{
		buff: make([]byte, buffSize),
	}
}

func (s *Service) OnConnect(ctx context.Context, w server.ResponseWriter) error {
	_, err := w.Write(ctx, []byte("Welcome, to the real world!"))

	switch {
	case err == nil:
	case errors.Is(err, io.EOF):
		zap.L().Info("Connection closed by client")
		return io.EOF
	default:
		return pkgerr.Wrap(err, "client service failed send initial data to client")
	}

	return nil
}

func (s *Service) OnData(ctx context.Context, r server.ResponseReader, w server.ResponseWriter) error {
	n, err := r.Data().Read(s.buff)

	switch {
	case err == nil: // Everything is ok, just read data
	case errors.Is(err, os.ErrDeadlineExceeded) || n == 0: // No data received
	case errors.Is(err, io.EOF):
		zap.L().Info("Connection closed by client")
	default:
		zap.L().Debug("Error reading data from connection", zap.Error(err))
		return pkgerr.Wrap(err, "client service failed read data from the connection")
	}

	if _, err := w.Write(ctx, s.buff[:n]); err != nil {
		return pkgerr.Wrap(err, "client service failed write data to the connection")
	}

	return io.EOF
}

func (s *Service) OnDisconnect(ctx context.Context) {
	// Internal state may be cleared here
}
