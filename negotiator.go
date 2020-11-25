package rest

import (
	"net/http"

	"github.com/ehsoc/rest/encdec"
)

// Negotiator interface describe the necessary methods to process the HTTP content negotiation logic.
type Negotiator interface {
	NegotiateEncoder(*http.Request, *ContentTypes) (MIMEType string, encoder encdec.Encoder, err error)
	NegotiateDecoder(*http.Request, *ContentTypes) (MIMEType string, decoder encdec.Decoder, err error)
}
