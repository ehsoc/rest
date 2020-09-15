package encdec

import (
	"encoding/xml"
	"io"
)

//XMLDecoder implements restapigen.Decoder to decode xml format
type XMLDecoder struct{}

//Decode implements Decode method of interface Decoder
func (x XMLDecoder) Decode(r io.Reader, v interface{}) error {
	encoder := xml.NewDecoder(r)
	return encoder.Decode(v)
}
