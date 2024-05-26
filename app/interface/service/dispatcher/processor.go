package dispatcher

import (
	"context"
	"io"

	pkgerr "github.com/pkg/errors"
)

type Processor[T any] struct {
	typeName string
	handler  func(ctx context.Context, pkg *T, w io.Writer) error
}

func NewProcessor[T any](h func(ctx context.Context, pkg *T, w io.Writer) error) *Processor[T] {
	return &Processor[T]{
		typeName: GetTypeName[T](),
		handler:  h,
	}
}

func (p *Processor[T]) TypeName() string {
	return p.typeName
}

func (p *Processor[T]) Handle(ctx context.Context, codec StreamDecoder, payload []byte, w io.Writer) error {
	inst := new(T)

	if err := codec.Unmarshal(payload, inst); err != nil {
		return pkgerr.Wrapf(err, "handler[%s] failed to decode package payload", p.typeName)
	}

	if err := p.handler(ctx, inst, w); err != nil {
		return pkgerr.Wrapf(err, "handler[%s] got error from underlying logic", p.typeName)
	}

	return nil
}
