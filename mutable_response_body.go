package rest

// MutableResponseBody is an interface that represents the Http body response
// that can mutate after a validation or an operation, taking the outputs of this methods as inputs.
type MutableResponseBody interface {
	Mutate(operationResultBody interface{}, success bool, err error)
}
