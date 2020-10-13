package resource

// Operation defines a resource operation.
// Execute function will execute the operation.
// Return values:
// Body: can be nil, and it will be the body to be returned if it sets to do so.
// Success: will be true if the operation did what the client was expecting in the most positive outcome.
// Success false means that the operation has no errors, but the positive outcome was not achieved (something was not found in a database)
// Error: means that an error was risen and the operation was not executed because an internal error in the API.
// Error should cause a 500's error code.
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
