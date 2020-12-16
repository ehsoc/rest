package encdec

import (
	"io"
	"reflect"
)

// TextEncoder implements Encoder to encode on text format
type TextEncoder struct{}

// Encode implements method of Encoder interface
func (t TextEncoder) Encode(w io.Writer, v interface{}) error {
	rv := reflect.ValueOf(v)

	if rv.Kind() == reflect.Ptr {
		if rv.Elem().Kind() != reflect.String {
			return ErrorTextDecoderNoString
		}
		v = rv.Elem().Interface()
	}

	s, ok := v.(string)
	if !ok {
		return ErrorTextDecoderNoString
	}

	_, err := w.Write([]byte(s))

	return err
}
