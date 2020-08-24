package resource

import (
	"net/http"

	"github.com/ehsoc/resource/encdec"
)

//Negotiator takes a response header and contains the content-type negotiation logic.
//Negotiate function will return the mimetype name, the encoder-decoder and error. Error response will
//produce a 415 code ("Status Unsupported Media Type") response by the handler to the client.
//defaultNegotiatorFunc (the default negotiator) is an example of a Negotiate function implementation
type Negotiator interface {
	NegotiateEncoder(*http.Request, *HTTPContentTypeSelector) (mimeType string, encoder encdec.Encoder, err error)
	NegotiateDecoder(*http.Request, *HTTPContentTypeSelector) (mimeType string, decoder encdec.Decoder, err error)
}
