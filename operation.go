package rest

// Operation defines a resource operation.
// Execute method will execute the operation.
// Return values:
// `body`: can be nil, and it will be the body to be returned if the operation's success response is set
// with `WithOperationResultBody`.
// `success`: should be true if the operation got the most positive outcome.
// Success with false value means that the operation has no errors, but the positive outcome was not achieved (something was not found in a database).
// `err`: means that an error was risen and the operation was not executed because an internal error in the API.
// Error should cause a 500's http error code.
type Operation interface {
	Execute(i Input) (body interface{}, success bool, err error)
}

// The OperationFunc type is an adapter to allow the use of
// ordinary functions as Operation. If f is a function
// with the appropriate signature, OperationFunc(f) is a
// Operation that calls f.
type OperationFunc func(i Input) (body interface{}, success bool, err error)

// Execute calls f(i)
func (f OperationFunc) Execute(i Input) (body interface{}, success bool, err error) {
	return f(i)
}
