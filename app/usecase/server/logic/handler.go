package logic

import "context"

type Handler struct {
	askStop func()
	send    func(ctx context.Context, data []byte) error
}

func New(askStop func(), send func(ctx context.Context, data []byte) error) *Handler {
	return &Handler{
		askStop: askStop,
		send:    send,
	}
}

func (h *Handler) Handle(ctx context.Context, req []byte) error {
	h.send(ctx, req)
	h.askStop()
	return nil
}
