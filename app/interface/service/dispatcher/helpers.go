package dispatcher

import (
	"faraway/wow/pkg/protocol"
	"io"
	"reflect"

	"github.com/pkg/errors"
)

type Encoder interface {
	Marshal(v any) ([]byte, error)
}

type Decoder interface {
	Unmarshal(data []byte, v any) error
}

type StreamDecoder interface {
	GetRaw(r io.Reader, buff []byte) (*protocol.Package, error)
	Decoder
}

func GetTypeName[T any]() string {
	var t [0]T

	if typeT := reflect.TypeOf(t).Elem(); typeT != nil {
		return typeT.Name()
	}

	return ""
}

func EncodePackage[T any](v *T, enc Encoder) ([]byte, error) {
	typeName := GetTypeName[T]()

	payload, err := enc.Marshal(v)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal payload of type '%s'", typeName)
	}

	pkg := protocol.Package{
		Tag:     typeName,
		Payload: payload,
	}

	data, err := enc.Marshal(&pkg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal package")
	}

	return data, nil
}
