package acceptor

import (
	"context"
	"errors"
	"io"
	"net"
	"os"
	"time"

	"go.uber.org/zap"
)

type Handler interface {
	Exchange(ctx context.Context, req []byte) ([]byte, error)
}

type Acceptor struct {
	buffer  []byte
	timeout time.Duration
	handler Handler
}

func New(handler Handler, buffSize int, timeout time.Duration) *Acceptor {
	return &Acceptor{
		buffer:  make([]byte, buffSize),
		timeout: timeout,
		handler: handler,
	}
}

func (c *Acceptor) Accept(ctx context.Context, conn net.Conn) {
	// Set read deadline
	if err := conn.SetReadDeadline(time.Now().Add(c.timeout)); err != nil {
		zap.L().Debug("Error setting read deadline for the connection", zap.Error(err))
		return
	}

	for {
		select {
		case <-ctx.Done():
			return

		default:
			n, err := conn.Read(c.buffer)
			if err != nil {
				c.handleReadError(err, conn.RemoteAddr().String())
				return
			}

			// No data received, wait until the next call, timeout or context shutdown
			if n <= 0 {
				continue
			}

			resp, err := c.handler.Exchange(ctx, c.buffer[:n])
			if err != nil {
				// Some verbose response may be sent here
				return
			}

			if err := conn.SetWriteDeadline(time.Now().Add(c.timeout)); err != nil {
				zap.L().Debug("Error setting write deadline for the connection", zap.Error(err))
				return
			}

			if _, err := conn.Write(resp); err != nil {
				c.handleWriteError(err, conn.RemoteAddr().String())
			}

			// We are done here, just stop the logic and close the connection
			return
		}
	}
}

func (c *Acceptor) handleReadError(err error, addr string) {
	switch {
	case errors.Is(err, io.EOF):
		zap.L().Info("Connection closed by client", zap.String("addr", addr))
	case errors.Is(err, os.ErrDeadlineExceeded):
		// Just close connection. In other cases new deadline may be set.
	default:
		zap.L().Error("Error reading from connection", zap.Error(err))
	}
}

func (c *Acceptor) handleWriteError(err error, addr string) {
	switch {
	case errors.Is(err, io.EOF):
		zap.L().Info("Connection closed by client", zap.String("addr", addr))
	case errors.Is(err, os.ErrDeadlineExceeded):
		// Do nothing, just close connection
	default:
		zap.L().Error("Error writing to connection", zap.Error(err))
	}
}
