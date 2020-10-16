package resource

import (
	"net/http"

	"github.com/ehsoc/resource/encdec"
)

// Negotiator interface describe the necessary methods to process a HTTP content negotiation logic.
type Negotiator interface {
	NegotiateEncoder(*http.Request, *HTTPContentTypeSelector) (MIMEType string, encoder encdec.Encoder, err error)
	NegotiateDecoder(*http.Request, *HTTPContentTypeSelector) (MIMEType string, decoder encdec.Decoder, err error)
}
