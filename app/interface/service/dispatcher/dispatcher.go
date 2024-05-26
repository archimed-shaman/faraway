package dispatcher

import (
	"context"
	"errors"
	"fmt"
	"io"

	pkgerr "github.com/pkg/errors"
)

var ErrBadType = errors.New("bad type")

type TypeProcessor interface {
	Handle(ctx context.Context, codec StreamDecoder, payload []byte, w io.Writer) error
	TypeName() string
}

type Dispatcher struct {
	typeH map[string]TypeProcessor
	codec StreamDecoder
}

func New(codec StreamDecoder) *Dispatcher {
	return &Dispatcher{
		typeH: make(map[string]TypeProcessor),
		codec: codec,
	}
}

func (d *Dispatcher) Register(matchers ...TypeProcessor) {
	if d.typeH == nil {
		d.typeH = make(map[string]TypeProcessor)
	}

	for _, matcher := range matchers {
		if _, found := d.typeH[matcher.TypeName()]; found {
			// TODO: return error
			panic(fmt.Sprintf("type '%s' is already registered", matcher.TypeName()))
		}

		d.typeH[matcher.TypeName()] = matcher
	}
}

func (d *Dispatcher) Dispatch(ctx context.Context, r io.Reader, w io.Writer, buff []byte) error {
	pkg, err := d.codec.GetRaw(r, buff)
	if err != nil {
		return pkgerr.Wrap(err, "dispatcher failed to get raw package")
	}

	h, found := d.typeH[pkg.Tag]
	if !found {
		return pkgerr.Wrapf(ErrBadType, "Unknown type '%s'", pkg.Tag)
	}

	if err := h.Handle(ctx, d.codec, pkg.Payload, w); err != nil {
		return pkgerr.Wrap(err, "registered handler returned error")
	}

	return nil
}
