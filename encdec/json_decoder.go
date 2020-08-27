package encdec

import (
	"encoding/json"
	"io"
)

//JSONDecoder implements Decoder to encode on json format
type JSONDecoder struct{}

//Decode implements Decode method of interface Decoder
func (j JSONDecoder) Decode(r io.Reader, v interface{}) error {
	encoder := json.NewDecoder(r)
	return encoder.Decode(v)
}
