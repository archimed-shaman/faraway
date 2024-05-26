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
	conn net.Conn, timeout time.Duration, buffSize int,
	service Service,
) {
	for {
		if err := conn.SetDeadline(time.Now().Add(timeout)); err != nil {
			zap.L().Debug("Error setting read deadline for the connection", zap.Error(err))
			return
		}

		select {
		case <-ctx.Done():
			return

		default:
			err := service.Handle(ctx, conn, conn)

			switch {
			case err == nil: // Everything is ok, just read data
			case errors.Is(err, os.ErrDeadlineExceeded):
				return
			case errors.Is(err, io.EOF):
				return
			default:
				zap.L().Info("Error reading data from connection", zap.Error(err))
				return
			}
		}
	}
}
