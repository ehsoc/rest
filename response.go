package resource

// Response represents a HTTP response.
// A response with code 0 will be consider a nil response.
// MutableResponseBody is an interface that represents the Http body response, that can mutate after a validation or an operation, taking the outputs as
// inputs.
type Response struct {
	code int
	MutableResponseBody
	description string
}

// NewResponse returns a Response with the specified code.
// A response with code 0 will be consider a nil response.
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

// WithBody will set the body property.
// It generates a MutableResponseBody implementation under the hood.
func (r Response) WithBody(body interface{}) Response {
	r.MutableResponseBody = staticResponseBody{body}
	return r
}

// WithOperationResultBody will set the body property.
// It will also return the result of operation body as the response body.
// It generates a MutableResponseBody implementation under the hood.
func (r Response) WithOperationResultBody(body interface{}) Response {
	r.MutableResponseBody = &operationResponseBody{body}
	return r
}

// WithMutableBody will set the MutableResponseBody property.
// This will be used to mutates the result body, using the results of validation or an operation as inputs.
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

// Body returns the body property.
func (r Response) Body() interface{} {
	if staticResponse, ok := r.MutableResponseBody.(staticResponseBody); ok {
		return staticResponse.body
	}
	if oResponse, ok := r.MutableResponseBody.(*operationResponseBody); ok {
		return oResponse.body
	}
	return r.MutableResponseBody
}
