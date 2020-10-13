package resource

import (
	"net/http"

	"github.com/ehsoc/resource/encdec"
)

// Negotiator interface describe the methods to process a HTTP content type negotiation logic.
type Negotiator interface {
	NegotiateEncoder(*http.Request, *HTTPContentTypeSelector) (mimeType string, encoder encdec.Encoder, err error)
	NegotiateDecoder(*http.Request, *HTTPContentTypeSelector) (mimeType string, decoder encdec.Decoder, err error)
}
