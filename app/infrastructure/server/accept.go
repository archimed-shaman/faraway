package server

import (
	"context"
	"errors"
	"io"
	"net"
	"time"

	pkgerr "github.com/pkg/errors"
	"go.uber.org/zap"
)

func (w *Writer) Write(ctx context.Context, data []byte) (int, error) {
	if err := w.conn.SetWriteDeadline(time.Now().Add(w.timeout)); err != nil {
		return 0, pkgerr.Wrap(err, "failed to set write timeout")
	}

	n, err := w.conn.Write(data)
	if err != nil {
		return n, pkgerr.Wrap(err, "response writer failed to write data")
	}

	return n, nil
}

func handleClient(ctx context.Context, conn net.Conn, timeout time.Duration, handler Service) {
	writer := NewWriter(conn, timeout)
	reader := NewReader(conn)

	if err := handler.OnConnect(ctx, writer); err != nil {
		switch {
		case err == nil:
		case errors.Is(err, io.EOF):
			return
		default:
			zap.L().Debug("Error handling OnConnect", zap.Error(err))
			return
		}
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
			// case errors.Is(err, os.ErrDeadlineExceeded): // No data received
			case errors.Is(err, io.EOF):
				// 	zap.L().Info("Connection closed by client", zap.String("addr", conn.RemoteAddr().String()))
				return
			default:
				zap.L().Debug("Error processing data from connection", zap.Error(err))
				return
				// 	zap.L().Error("Error reading from connection", zap.Error(err))
				// 	return
				// }

				// if err := handler(ctx, buffer[:n]); err != nil {
				// 	return
			}
		}
	}
}
