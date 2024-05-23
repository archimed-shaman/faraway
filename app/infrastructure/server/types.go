package server

import (
	"context"
	"io"
	"net"
	"time"
)

type ResponseWriter interface {
	// os.ErrDeadlineExceeded may be returned on write operations
	Write(ctx context.Context, data []byte) (int, error)
}

type ResponseReader interface {
	// os.ErrDeadlineExceeded may be returned on read operations
	Data() io.Reader
}

type Service interface {
	// io.EOF is expected as the signal of end of processing.
	OnConnect(ctx context.Context, w ResponseWriter) error

	// Connection will be served until err is nil.
	// io.EOF is expected as the signal of end of processing.
	OnData(ctx context.Context, r ResponseReader, w ResponseWriter) error

	OnDisconnect(ctx context.Context)
}

type Reader struct {
	conn net.Conn
}

func NewReader(conn net.Conn) *Reader {
	return &Reader{
		conn: conn,
	}
}

func (r *Reader) Data() io.Reader {
	return r.conn
}

type Writer struct {
	timeout time.Duration
	conn    net.Conn
}

func NewWriter(conn net.Conn, timeout time.Duration) *Writer {
	return &Writer{
		timeout: timeout,
		conn:    conn,
	}
}
