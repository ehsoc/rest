package httputil

import (
	"mime"
	"strings"
)

// ParseMediaTypes will parse one or multiple media type directives (i.e Accept or Content-Type headers)
// and return an array of MediaType representing the MIME type name and optional parameters.
func ParseMediaTypes(accept string) []MediaType {
	mediatypes := []MediaType{}
	types := strings.Split(accept, ",")
	for _, ct := range types {
		name, params, _ := mime.ParseMediaType(ct)
		mediatypes = append(mediatypes, MediaType{name, params})
	}
	return mediatypes
}
