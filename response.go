package resource

// Response represents a HTTP response.
// A response with code 0 will be consider a nil response.
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

type StaticResponseBody struct {
	body interface{}
}

func (s StaticResponseBody) Mutate(v interface{}, success bool, err error) {

}

//WithBody chain method will set body property.
func (r Response) WithBody(body interface{}) Response {
	r.MutableResponseBody = StaticResponseBody{body}
	return r
}

//WithMutableBody chain method will set body property.
func (r Response) WithMutableBody(mutableResponseBody MutableResponseBody) Response {
	r.MutableResponseBody = mutableResponseBody
	return r
}

//WithDescription chain method will set description property.
func (r Response) WithDescription(description string) Response {
	r.description = description
	return r
}

//Code returns the code property
func (r Response) Code() int {
	return r.code
}

//Description returns the description property
func (r Response) Description() string {
	return r.description
}

//Body returns the body property
func (r Response) Body() interface{} {
	if staticResponse, ok := r.MutableResponseBody.(StaticResponseBody); ok {
		return staticResponse.body
	}
	return r.MutableResponseBody
}
