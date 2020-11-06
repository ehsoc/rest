package encdec

import "io"

// Decoder purpose is to provide a common interface wraping around decoder libraries
// so that they can easily be passed as arguments values to be executed by other methods.
type Decoder interface {
	Decode(r io.Reader, v interface{}) error
}

// DecoderFunc type is an adapter to allow the use of
// ordinary functions as Decoders. If f is a function
// with the appropriate signature, DecoderFunc(f) is a
// Decoder that calls f.
type DecoderFunc func(r io.Reader, v interface{}) error

// Decode calls e(w,v)
func (e DecoderFunc) Decode(r io.Reader, v interface{}) error {
	return e(r, v)
}
