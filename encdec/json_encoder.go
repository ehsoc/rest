package encdec

import (
	"encoding/json"
	"io"
)

// JSONEncoder implements Encoder interface to encode JSON format
type JSONEncoder struct{}

// Encode writes to w the JSON encoding of v
func (j JSONEncoder) Encode(w io.Writer, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}
