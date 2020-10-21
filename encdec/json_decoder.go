package encdec

import (
	"encoding/json"
	"io"
)

// JSONDecoder implements Decoder interface to decode JSON format
type JSONDecoder struct{}

// Decode is a wrapper around encoding/json JSON Decoder, that read from r and store it in v
func (j JSONDecoder) Decode(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)

}
