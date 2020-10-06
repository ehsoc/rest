package resource

type RequestBody struct {
	Description string
	Body        interface{}
	Required    bool
}
