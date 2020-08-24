package resource

import (
	"errors"
	"net/http"

	"github.com/ehsoc/resource/encdec"
)

//HTTPContentTypeSelector contains all available content-types mimetypes,
//associated with their respective encoder-decoders. Implements Negotiator interface through
//embedding a Negotiator.
type HTTPContentTypeSelector struct {
	contentTypes          map[string]encdec.EncoderDecoder
	defaultEncoderDecoder string
	Negotiator
}

//NewHTTPContentTypeSelector will return a HTTPContentTypeSelector with an empty content-type
//map and the Default Negotiator.
//The Negotiator is a content-type negotiator, that can be replaced by a custom Negotiator implementation.
func NewHTTPContentTypeSelector() HTTPContentTypeSelector {
	var contentTypes = make(map[string]encdec.EncoderDecoder)
	return HTTPContentTypeSelector{contentTypes, "", DefaultNegotiator{}}
}

//Add will add a new content-type encoder-decoder. defaultencdec parameter will set the default one.
func (h *HTTPContentTypeSelector) Add(contentType string, ed encdec.EncoderDecoder, defaultencdec bool) {
	h.contentTypes[contentType] = ed
	if defaultencdec {
		h.defaultEncoderDecoder = contentType
	}
}

func (h *HTTPContentTypeSelector) GetEncoderDecoder(contentType string) (encdec.EncoderDecoder, error) {
	if ed, ok := h.contentTypes[contentType]; ok {
		return ed, nil
	}
	return nil, errors.New("content type not found")
}

func (h *HTTPContentTypeSelector) GetDefaultEncoderDecoder() (string, encdec.EncoderDecoder, error) {
	if ed, ok := h.contentTypes[h.defaultEncoderDecoder]; ok {
		return h.defaultEncoderDecoder, ed, nil
	}
	return "", nil, errors.New("no default content-type is set")
}

func (h *HTTPContentTypeSelector) NegotiateEncoder(r *http.Request) (string, encdec.Encoder, error) {
	return h.Negotiator.NegotiateEncoder(r, h)
}

func (h *HTTPContentTypeSelector) NegotiateDecoder(r *http.Request) (string, encdec.Decoder, error) {
	return h.Negotiator.NegotiateDecoder(r, h)
}
