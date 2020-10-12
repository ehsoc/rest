package resource

type MutableResponseBody interface {
	Mutate(operationResultBody interface{}, success bool, err error)
}

// The MutateFunc type is an adapter to allow the use of
// ordinary functions as ResponseBody. If f is a function
// with the appropriate signature, MutateFunc(f) is a
// Mutate that calls f.
type MutateFunc func(operationResultBody interface{}, success bool, err error)

//Execute calls f(body, parameters, decoder)
func (f MutateFunc) Mutate(o interface{}, s bool, e error) {
	f(o, s, e)
}
