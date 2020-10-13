package resource

// MethodOperation contains the operation, and its responses in case of success or failure.
// Operation Execute method will return a body (interface{}), success (bool), and err (error).
// The error case will not be manage by MethodOperation.
type MethodOperation struct {
	// Logic operation of the method
	Operation
	// Response if Operation Execute method function returns a success with true value.
	successResponse Response
	// Response if Operation Execute method function returns success with false value.
	// Failure is not an error.
	// response code == 0 ,indicate that there is no response in case of Operation failure.
	failResponse Response
}

// NewMethodOperation returns a new MethodOperation instance.
// successResponse code property cannot be 0. 0 code means a nil response, and this parameter is required.
// This response will be returned if the operation success return value is true and error is nil.
// failedResponse response property with code 0, means that there is no response in case of Operation failure.
// This response will be returned if the operation success value is false.
// Please check if your operation returns a success false, if you don't define a failure response (a response with code 0),
// and your operation returns a success false, the HTTP Server could return a panic.
func NewMethodOperation(operation Operation, successResponse, failedResponse Response) MethodOperation {
	if successResponse.code == 0 {
		panic(ErrorNilCodeSuccessResponse)
	}
	return MethodOperation{operation, successResponse, failedResponse}
}
