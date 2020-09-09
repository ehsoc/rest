package resource

import (
	"errors"
	"net/http"
	"strings"

	"github.com/ehsoc/resource/encdec"
	"github.com/ehsoc/resource/httputil"
)

type DefaultNegotiator struct {
}

//NegotiateEncoder is a method implementation of Negotiator interface
func (d DefaultNegotiator) NegotiateEncoder(r *http.Request, cts *HTTPContentTypeSelector) (string, encdec.Encoder, error) {
	accept := r.Header.Get("Accept")
	if strings.Trim(accept, "") != "" {
		mediaTypes := httputil.ParseContentType(accept)
		for _, mediaType := range mediaTypes {
			enc, err := cts.GetEncoder(mediaType.Name)
			if err == nil {
				return mediaType.Name, enc, nil
			}
		}
	}
	return cts.GetDefaultEncoder()
}

//NegotiateDecoder is a default method implementation of Negotiator interface
func (d DefaultNegotiator) NegotiateDecoder(r *http.Request, cts *HTTPContentTypeSelector) (string, encdec.Decoder, error) {
	ct := r.Header.Get("Content-Type")
	//Only if it is a non-empty and not nil Body we will require a Content-Type header
	if r.Body != http.NoBody && r.Body != nil {
		if strings.Trim(ct, "") != "" {
			mediaTypes := httputil.ParseContentType(ct)
			for _, mediaType := range mediaTypes {
				enc, err := cts.GetDecoder(mediaType.Name)
				if err == nil {
					return mediaType.Name, enc, nil
				}
			}
		}
		return "", nil, errors.New("unavailable content-type decoder")
	}
	return cts.GetDefaultDecoder()
}
