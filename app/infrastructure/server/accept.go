package server

import (
	"context"
	"errors"
	"io"
	"net"
	"os"
	"time"

	"go.uber.org/zap"
)

func handleClient(ctx context.Context,
	conn net.Conn, timeout time.Duration, buffer []byte,
	handler func(ctx context.Context, req []byte) error,
) {
	for {
		// Set read deadline
		if err := conn.SetDeadline(time.Now().Add(timeout)); err != nil {
			zap.L().Debug("Error setting read deadline for the connection", zap.Error(err))
			return
		}

		select {
		case <-ctx.Done():
			return

		default:
			n, err := conn.Read(buffer)

			switch {
			case err == nil: // Everything is ok, just read data
			case errors.Is(err, os.ErrDeadlineExceeded): // No data received
			case errors.Is(err, io.EOF):
				zap.L().Info("Connection closed by client", zap.String("addr", conn.RemoteAddr().String()))
				return
			default:
				zap.L().Error("Error reading from connection", zap.Error(err))
				return
			}

			// No data received, wait until the next call, timeout or context shutdown
			if n <= 0 {
				continue
			}

			if err := handler(ctx, buffer[:n]); err != nil {
				return
			}
		}
	}
}
