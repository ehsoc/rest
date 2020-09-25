package resource

//Response represents a HTTP response.
//This type is meant to be immutable.
type Response struct {
	code        int
	body        interface{}
	description string
}

//NewResponse resturns a Response
func NewResponse(code int) Response {
	r := Response{}
	r.code = code
	return r
}

//WithBody chain method will set body property.
//This will not work outside a chain.
func (r Response) WithBody(body interface{}) Response {
	r.body = body
	return r
}

//WithDescription chain method will set description property.
//This will not work outside a chain.
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
	return r.body
}
