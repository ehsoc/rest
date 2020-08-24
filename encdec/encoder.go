package encdec

import "io"

//Encoder interface purpose is to compose different encoders based on
//libraries like encoding/json or encoding/xml
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
