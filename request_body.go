package rest

// RequestBody represents a request body parameter
type RequestBody struct {
	Description string
	Body        interface{}
	Required    bool
}
