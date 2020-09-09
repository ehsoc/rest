package resource

import (
	"io"
	"net/url"

	"github.com/ehsoc/resource/encdec"
)

//Operation defines an operation over a data entity
//Execute function will execute the operation.
type Operation interface {
	Execute(body io.ReadCloser, parameters url.Values, decoder encdec.Decoder) (interface{}, error)
}
