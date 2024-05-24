package json

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// Just simple encoder for client-server communication
// Yep, I know, even JSON could be more effective.
// But for the real production we would likely prefer protobuf/msgpack/bson/custom TLV protocol/...

type Codec struct{}

func NewCodec() *Codec {
	return &Codec{}
}

func (s *Codec) Marshal(v any) ([]byte, error) {
	data, err := json.Marshal(v)
	return data, errors.Wrap(err, "json codec failed to marshal")
}

func (s *Codec) Unmarshal(data []byte, v any) error {
	err := json.Unmarshal(data, v)
	return errors.Wrap(err, "json codec failed to unmarshal")
}
