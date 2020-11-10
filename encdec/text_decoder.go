package encdec

import (
	"io"
	"io/ioutil"
	"reflect"
)

// TextDecoder implements Decoder to encode on text format
type TextDecoder struct{}

// Decode implements Decode method of interface Decoder
func (t TextDecoder) Decode(r io.Reader, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return ErrorTextDecoderNoValidPointer
	}

	if rv.Elem().Type().Kind() != reflect.String {
		return ErrorTextDecoderNoString
	}

	b, err := ioutil.ReadAll(r)

	if err != nil {
		return err
	}

	rv.Elem().SetString(string(b))

	return nil
}
