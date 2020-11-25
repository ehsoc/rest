package rest

import (
	"github.com/ehsoc/rest/encdec"
)

// ContentTypes contains all available MIME types,
// associated with their respective Encoder/Decoder.
type ContentTypes struct {
	encoderContentTypes          map[string]encdec.Encoder
	decoderContentTypes          map[string]encdec.Decoder
	defaultEncoder               string
	defaultDecoder               string
	UnsupportedMediaTypeResponse Response
}

// NewContentTypes will create a new ContentTypes instance
func NewContentTypes() ContentTypes {
	var encoderContentTypes = make(map[string]encdec.Encoder)

	var decoderContentTypes = make(map[string]encdec.Decoder)

	return ContentTypes{encoderContentTypes, decoderContentTypes, "", "", NewResponse(415)}
}

// Add adds a EncoderDecoder.
// isDefault parameter will set the default encoder and decoder.
func (h *ContentTypes) Add(mimeType string, ed encdec.EncoderDecoder, isDefault bool) {
	h.AddEncoder(mimeType, ed, isDefault)
	h.AddDecoder(mimeType, ed, isDefault)
}

// AddEncoder adds a encoder.
// isDefault parameter sets this encoder as default.
func (h *ContentTypes) AddEncoder(mimeType string, encoder encdec.Encoder, isDefault bool) {
	h.checkNilEncoderMap()
	h.encoderContentTypes[mimeType] = encoder

	if isDefault {
		h.defaultEncoder = mimeType
	}
}

// AddDecoder adds a decoder.
// isDefault parameter will set this decoder as default.
func (h *ContentTypes) AddDecoder(mimeType string, decoder encdec.Decoder, isDefault bool) {
	h.checkNilDecoderMap()
	h.decoderContentTypes[mimeType] = decoder

	if isDefault {
		h.defaultDecoder = mimeType
	}
}

// GetEncoder gets the encoder with the provided mimeType as key
func (h *ContentTypes) GetEncoder(mimeType string) (encdec.Encoder, error) {
	if ed, ok := h.encoderContentTypes[mimeType]; ok {
		return ed, nil
	}
	return nil, ErrorNoDefaultContentTypeIsSet
}

// GetDecoder gets the decoder with the provided mimeType as key
func (h *ContentTypes) GetDecoder(mimeType string) (encdec.Decoder, error) {
	if ed, ok := h.decoderContentTypes[mimeType]; ok {
		return ed, nil
	}
	return nil, ErrorNoDefaultContentTypeIsSet
}

// GetDefaultEncoder gets default encoder.
func (h *ContentTypes) GetDefaultEncoder() (string, encdec.Encoder, error) {
	if ed, ok := h.encoderContentTypes[h.defaultEncoder]; ok {
		return h.defaultEncoder, ed, nil
	}
	return "", nil, ErrorNoDefaultContentTypeIsSet
}

// GetDefaultDecoder gets the default decoder.
func (h *ContentTypes) GetDefaultDecoder() (string, encdec.Decoder, error) {
	if ed, ok := h.decoderContentTypes[h.defaultEncoder]; ok {
		return h.defaultDecoder, ed, nil
	}
	return "", nil, ErrorNoDefaultContentTypeIsSet
}

// checkDecoderMap initialize the internal map if is nil
func (h *ContentTypes) checkNilDecoderMap() {
	if h.decoderContentTypes == nil {
		h.decoderContentTypes = make(map[string]encdec.Decoder)
	}
}

// checkNilEncoderMap initialize the internal map if is nil
func (h *ContentTypes) checkNilEncoderMap() {
	if h.encoderContentTypes == nil {
		h.encoderContentTypes = make(map[string]encdec.Encoder)
	}
}
