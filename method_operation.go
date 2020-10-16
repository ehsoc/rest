package resource

// MethodOperation contains the operation, and its responses in case of success or failure.
// Operation Execute method will return a body (interface{}), success (bool), and err (error).
type MethodOperation struct {
	// Logic operation of the method
	Operation
	// Response if Operation Execute method function returns a success with true value.
	successResponse Response
	// Response if Operation Execute method function returns success with false value.
	// Failure is not an error.
	failResponse Response
}

// NewMethodOperation returns a new MethodOperation instance.
// successResponse:
// This response will be returned if the operation success return value is true, and the error value is nil.
// failedResponse:
// This response will be returned if the operation success value is false.
// Please check if your operation returns a success false, if you don't define a failure response,
// and your operation returns a success false, the HTTP Server could return a panic.
func NewMethodOperation(operation Operation, successResponse Response) MethodOperation {
	return MethodOperation{operation, successResponse, Response{disabled: true}}
}

// WithFailResponse sets the failResponse property
func (m MethodOperation) WithFailResponse(failResponse Response) MethodOperation {
	m.failResponse = failResponse
	return m
}
