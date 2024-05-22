package logic

import "context"

type Connection interface {
	Send(ctx context.Context, data []byte) error
}

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(ctx context.Context, req []byte) error {
	return nil
}
