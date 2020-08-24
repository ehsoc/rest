package httputil

//MediaType represents a media type from an Accept header.
//Name property is the media type name (like "application/json").
//Params property contains the media type option parameters.
type MediaType struct {
	Name   string
	Params map[string]string
}
