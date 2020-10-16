package resource

// Response represents a HTTP response.
// MutableResponseBody is an interface that represents the Http body response,
// that can mutate after a validation or an operation, taking the outputs of this methods as inputs.
type Response struct {
	code int
	MutableResponseBody
	description string
	disabled    bool
}

// NewResponse returns a Response with the specified code.
func NewResponse(code int) Response {
	r := Response{}
	r.code = code
	return r
}

type operationResponseBody struct {
	body interface{}
}

func (o *operationResponseBody) Mutate(v interface{}, success bool, err error) {
	o.body = v
}

type staticResponseBody struct {
	body interface{}
}

func (s staticResponseBody) Mutate(v interface{}, success bool, err error) {

}

// WithBody will set a static body property.
// It generates a dummy MutableResponseBody implementation under the hood, that will return the given 'body' parameter without change it.
func (r Response) WithBody(body interface{}) Response {
	r.MutableResponseBody = staticResponseBody{body}
	return r
}

// WithOperationResultBody will set the body property.
// It generates a MutableResponseBody implementation under the hood, that will take the body result of an Operation as the response body.
// The given body parameter will be used for the specification only.
func (r Response) WithOperationResultBody(body interface{}) Response {
	r.MutableResponseBody = &operationResponseBody{body}
	return r
}

// WithMutableBody will set the MutableResponseBody implementation.
// This will be used to mutate the result body, using the result values of validation or an operation as inputs.
// The struct of MutableResponseBody implementation will be used as response body for the specification.
func (r Response) WithMutableBody(mutableResponseBody MutableResponseBody) Response {
	r.MutableResponseBody = mutableResponseBody
	return r
}

// WithDescription will set description property.
func (r Response) WithDescription(description string) Response {
	r.description = description
	return r
}

// Code returns the code property
func (r Response) Code() int {
	return r.code
}

// Description returns the description property
func (r Response) Description() string {
	return r.description
}

// Body returns the MutableResponseBody property.
// In the case of a body that was set with WithBody or WithOperationResultBody methods, it will return the given 'body' parameter.
func (r Response) Body() interface{} {
	if staticResponse, ok := r.MutableResponseBody.(staticResponseBody); ok {
		return staticResponse.body
	}
	if oResponse, ok := r.MutableResponseBody.(*operationResponseBody); ok {
		return oResponse.body
	}
	return r.MutableResponseBody
}
