package resource

type Request struct {
	Description string
	body        interface{}
}

func (r Request) GetBody() interface{} {
	return r.body
}
