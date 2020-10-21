package encdec

import "io"

// Encoder purpose is to provide a common interface wraping around encoder libraries
// so that they can easily be passed as arguments values to be executed by other methods.
type Encoder interface {
	Encode(w io.Writer, v interface{}) error
}

// EncoderFunc type is an adapter to allow the use of
// ordinary functions as Encoders. If f is a function
// with the appropriate signature, EncoderFunc(f) is a
// Encoder that calls f.
type EncoderFunc func(w io.Writer, v interface{}) error

//Encode calls e(w,v)
func (e EncoderFunc) Encode(w io.Writer, v interface{}) error {
	return e(w, v)
}
