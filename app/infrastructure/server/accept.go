package server

import (
	"context"
	"errors"
	"faraway/wow/app/infrastructure/server/types"
	"io"
	"net"
	"time"

	"go.uber.org/zap"
)

func handleClient(ctx context.Context, conn net.Conn, timeout time.Duration, handler Service) {
	writer := types.NewWriter(conn, timeout)
	reader := types.NewReader(conn)

	err := handler.OnConnect(ctx, writer)

	switch {
	case err == nil:
	case errors.Is(err, io.EOF):
		return
	default:
		zap.L().Debug("Error handling OnConnect", zap.Error(err))
		return
	}

	defer handler.OnDisconnect(ctx)

	for {
		if err := conn.SetDeadline(time.Now().Add(timeout)); err != nil {
			zap.L().Debug("Error setting read deadline for the connection", zap.Error(err))
			return
		}

		select {
		case <-ctx.Done():
			return

		default:
			err := handler.OnData(ctx, reader, writer)

			switch {
			case err == nil: // Everything is ok, just read data
			case errors.Is(err, io.EOF):
				return
			default:
				zap.L().Debug("Error processing data from connection", zap.Error(err))
				return
			}
		}
	}
}
