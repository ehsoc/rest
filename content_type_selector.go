package resource

import (
	"net/http"

	"github.com/ehsoc/resource/encdec"
)

// HTTPContentTypeSelector contains all available MIME types,
// associated with their respective Encoder/Decoder. Implements Negotiator interface embedding a Negotiator.
type HTTPContentTypeSelector struct {
	encoderContentTypes map[string]encdec.Encoder
	decoderContentTypes map[string]encdec.Decoder
	defaultEncoder      string
	defaultDecoder      string
	Negotiator
	UnsupportedMediaTypeResponse Response
}

// NewHTTPContentTypeSelector will return a HTTPContentTypeSelector.
// Negotiator is a content-type negotiator, that can be replaced by a custom Negotiator implementation.
func NewHTTPContentTypeSelector() HTTPContentTypeSelector {
	var encoderContentTypes = make(map[string]encdec.Encoder)
	var decoderContentTypes = make(map[string]encdec.Decoder)
	return HTTPContentTypeSelector{encoderContentTypes, decoderContentTypes, "", "", DefaultNegotiator{}, NewResponse(415)}
}

// Add adds a content-type EncoderDecoder.
// isDefault parameter will set the default encoder and decoder.
func (h *HTTPContentTypeSelector) Add(contentType string, ed encdec.EncoderDecoder, isDefault bool) {
	h.AddEncoder(contentType, ed, isDefault)
	h.AddDecoder(contentType, ed, isDefault)
}

// AddEncoder adds a content-type encoder.
// isDefault parameter sets this encoder as default.
func (h *HTTPContentTypeSelector) AddEncoder(contentType string, encoder encdec.Encoder, isDefault bool) {
	h.checkNilEncoderMap()
	h.encoderContentTypes[contentType] = encoder
	if isDefault {
		h.defaultEncoder = contentType
	}
}

// AddDecoder adds a new content-type decoder.
// isDefault parameter will set this decoder as default.
func (h *HTTPContentTypeSelector) AddDecoder(contentType string, decoder encdec.Decoder, isDefault bool) {
	h.checkNilDecoderMap()
	h.decoderContentTypes[contentType] = decoder
	if isDefault {
		h.defaultDecoder = contentType
	}
}

// GetEncoder gets the encoder with the provided contentType as key
func (h *HTTPContentTypeSelector) GetEncoder(contentType string) (encdec.Encoder, error) {
	if ed, ok := h.encoderContentTypes[contentType]; ok {
		return ed, nil
	}
	return nil, ErrorNoDefaultContentTypeIsSet
}

// GetDecoder gets the decoder with the provided contentType as key
func (h *HTTPContentTypeSelector) GetDecoder(contentType string) (encdec.Decoder, error) {
	if ed, ok := h.decoderContentTypes[contentType]; ok {
		return ed, nil
	}
	return nil, ErrorNoDefaultContentTypeIsSet
}

// GetDefaultEncoder gets default encoder.
func (h *HTTPContentTypeSelector) GetDefaultEncoder() (string, encdec.Encoder, error) {
	if ed, ok := h.encoderContentTypes[h.defaultEncoder]; ok {
		return h.defaultEncoder, ed, nil
	}
	return "", nil, ErrorNoDefaultContentTypeIsSet
}

// GetDefaultDecoder gets the default decoder.
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
