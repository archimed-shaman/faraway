package logic

import "context"

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Exchange(ctx context.Context, req []byte) ([]byte, error) {
	return []byte("Hello"), nil
}
