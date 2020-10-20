package resource

import (
	"net/http"

	"github.com/ehsoc/resource/encdec"
)

// Negotiator interface describe the necessary methods to process the HTTP content negotiation logic.
type Negotiator interface {
	NegotiateEncoder(*http.Request, *Renderers) (MIMEType string, encoder encdec.Encoder, err error)
	NegotiateDecoder(*http.Request, *Renderers) (MIMEType string, decoder encdec.Decoder, err error)
}
