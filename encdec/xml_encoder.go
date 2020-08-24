package encdec

import (
	"encoding/xml"
	"io"
)

//XMLEncoder implements Encoder to encode on xml format
type XMLEncoder struct{}

//Encode implements method of Encoder interface
func (x XMLEncoder) Encode(w io.Writer, v interface{}) error {
	w.Write([]byte(xml.Header))
	encoder := xml.NewEncoder(w)
	return encoder.Encode(v)
}
