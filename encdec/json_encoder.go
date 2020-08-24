package encdec

import (
	"encoding/json"
	"io"
)

//JSONEncoder implements Encoder to encode on json format
type JSONEncoder struct{}

//Encode implements method of Encoder interface
func (j JSONEncoder) Encode(w io.Writer, v interface{}) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(v)
}
