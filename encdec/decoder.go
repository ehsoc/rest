package encdec

import "io"

//Decoder interface purpose is to compose different decoders based on
//libraries like encoding/json or encoding/xml
type Decoder interface {
	Decode(r io.Reader, v interface{}) error
}

// DecoderFunc type is an adapter to allow the use of
// ordinary functions as Decoders. If f is a function
// with the appropriate signature, DecoderFunc(f) is a
// Decoder that calls f.
type DecoderFunc func(r io.Reader, v interface{}) error

//Decode calls e(w,v)
func (e DecoderFunc) Decode(r io.Reader, v interface{}) error {
	return e(r, v)
}
