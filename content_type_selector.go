package resource

import (
	"net/http"

	"github.com/ehsoc/resource/encdec"
)

// HTTPContentTypeSelector contains all available content-types mimetypes,
// associated with their respective encoder-decoders. Implements Negotiator interface through
// embedding a Negotiator.
type HTTPContentTypeSelector struct {
	encoderContentTypes map[string]encdec.Encoder
	decoderContentTypes map[string]encdec.Decoder
	defaultEncoder      string
	defaultDecoder      string
	Negotiator
	UnsupportedMediaTypeResponse Response
}

var DefaultUnsupportedMediaResponse = NewResponse(415)

// NewHTTPContentTypeSelector will return a HTTPContentTypeSelector with an empty content-type
// map and the Default Negotiator.
// The Negotiator is a content-type negotiator, that can be replaced by a custom Negotiator implementation.
func NewHTTPContentTypeSelector() HTTPContentTypeSelector {
	var encoderContentTypes = make(map[string]encdec.Encoder)
	var decoderContentTypes = make(map[string]encdec.Decoder)
	return HTTPContentTypeSelector{encoderContentTypes, decoderContentTypes, "", "", DefaultNegotiator{}, DefaultUnsupportedMediaResponse}
}

// Add will add a new content-type encoder and decoder. defaultencdec encoder and decoder parameter will set the default for decoder and decoder.
func (h *HTTPContentTypeSelector) Add(contentType string, ed encdec.EncoderDecoder, defaultencdec bool) {
	h.AddEncoder(contentType, ed, defaultencdec)
	h.AddDecoder(contentType, ed, defaultencdec)
}

// AddEncoder will add a new content-type decoder. isDefault parameter will set this decoder as the default one.
func (h *HTTPContentTypeSelector) AddEncoder(contentType string, encoder encdec.Encoder, isDefault bool) {
	h.checkNilEncoderMap()
	h.encoderContentTypes[contentType] = encoder
	if isDefault {
		h.defaultEncoder = contentType
	}
}

// AddDecoder will add a new content-type decoder. isDefault parameter will set this decoder as the default one.
func (h *HTTPContentTypeSelector) AddDecoder(contentType string, decoder encdec.Decoder, isDefault bool) {
	h.checkNilDecoderMap()
	h.decoderContentTypes[contentType] = decoder
	if isDefault {
		h.defaultDecoder = contentType
	}
}

// GetEncoder will return the encoder with the provided contentType as key
func (h *HTTPContentTypeSelector) GetEncoder(contentType string) (encdec.Encoder, error) {
	if ed, ok := h.encoderContentTypes[contentType]; ok {
		return ed, nil
	}
	return nil, ErrorNoDefaultContentTypeIsSet
}

// GetDecoder will return the encoder with the provided contentType as key
func (h *HTTPContentTypeSelector) GetDecoder(contentType string) (encdec.Decoder, error) {
	if ed, ok := h.decoderContentTypes[contentType]; ok {
		return ed, nil
	}
	return nil, ErrorNoDefaultContentTypeIsSet
}

// GetDefaultEncoder returns the default encoder is one is defined
func (h *HTTPContentTypeSelector) GetDefaultEncoder() (string, encdec.Encoder, error) {
	if ed, ok := h.encoderContentTypes[h.defaultEncoder]; ok {
		return h.defaultEncoder, ed, nil
	}
	return "", nil, ErrorNoDefaultContentTypeIsSet
}

// GetDefaultDecoder returns the default encoder is one is defined
func (h *HTTPContentTypeSelector) GetDefaultDecoder() (string, encdec.Decoder, error) {
	if ed, ok := h.decoderContentTypes[h.defaultEncoder]; ok {
		return h.defaultDecoder, ed, nil
	}
	return "", nil, ErrorNoDefaultContentTypeIsSet
}

// NegotiateEncoder execute the NegotiateEncoder function of the Negotiator implementation
func (h *HTTPContentTypeSelector) NegotiateEncoder(r *http.Request) (string, encdec.Encoder, error) {
	return h.Negotiator.NegotiateEncoder(r, h)
}

// NegotiateDecoder execute the NegotiateDecoder function of the Negotiator implementation
func (h *HTTPContentTypeSelector) NegotiateDecoder(r *http.Request) (string, encdec.Decoder, error) {
	return h.Negotiator.NegotiateDecoder(r, h)
}

// checkDecoderMap initialize the internal map if is nil
func (h *HTTPContentTypeSelector) checkNilDecoderMap() {
	if h.decoderContentTypes == nil {
		h.decoderContentTypes = make(map[string]encdec.Decoder)
	}
}

// checkNilEncoderMap initialize the internal map if is nil
func (h *HTTPContentTypeSelector) checkNilEncoderMap() {
	if h.encoderContentTypes == nil {
		h.encoderContentTypes = make(map[string]encdec.Encoder)
	}
}
