package types

import (
	"context"
	"io"
	"net"
	"time"

	pkgerr "github.com/pkg/errors"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE

type ResponseWriter interface {
	// os.ErrDeadlineExceeded may be returned on write operations
	Write(ctx context.Context, data []byte) (int, error)
}

type ResponseReader interface {
	// os.ErrDeadlineExceeded may be returned on read operations
	Data() io.Reader
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
