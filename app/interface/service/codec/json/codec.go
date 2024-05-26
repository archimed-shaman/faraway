package json

import (
	"faraway/wow/pkg/protocol"
	"io"

	"github.com/goccy/go-json"
	pkgerr "github.com/pkg/errors"
)

type Codec struct{}

func NewCodec() *Codec {
	return &Codec{}
}

func (s *Codec) GetRaw(r io.Reader, buff []byte) (*protocol.Package, error) {
	var pkg protocol.Package

	d := json.NewDecoder(r)
	if err := d.Decode(&pkg); err != nil {
		return nil, pkgerr.Wrap(err, "json codec failed to get raw")
	}

	return &pkg, nil
}

func (s *Codec) Marshal(v any) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, pkgerr.Wrap(err, "json codec failed to marshal package")
	}

	return data, nil
}

func (s *Codec) Unmarshal(data []byte, v any) error {
	err := json.Unmarshal(data, v)
	return pkgerr.Wrap(err, "json codec failed to unmarshal payload")
}
