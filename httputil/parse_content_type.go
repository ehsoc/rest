package httputil

import (
	"mime"
	"strings"
)

//ParseContentType will parse a Accept header and return an array of MediaType type
//representing the media type and options parameters
func ParseContentType(accept string) []MediaType {
	mediatypes := []MediaType{}
	types := strings.Split(accept, ",")
	for _, ct := range types {
		name, params, _ := mime.ParseMediaType(ct)
		mediatypes = append(mediatypes, MediaType{name, params})
	}
	return mediatypes
}
