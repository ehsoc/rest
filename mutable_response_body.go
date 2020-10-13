package resource

type MutableResponseBody interface {
	Mutate(operationResultBody interface{}, success bool, err error)
}
