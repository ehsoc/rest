package resource

// MethodOperation is a logic operation that is performed in an entity
type MethodOperation struct {
	// Logic of the method
	Operation
	// Response if Operation Execute method function returns no error
	successResponse Response
	// Response if Operation Execute method function returns error.
	// If code == 0 indicate that there is no response in case of Operation error.
	failResponse Response
	// Return the entity (interface{}) returned by Execute method,
	// in the Body of the response instead the Body of successResponse.
	// only the Code of successResponse will be returned in the http response.
	operationResultAsBody bool
}

// NewMethodOperation returns a new MethodOperation instance.
// successResponse code property cannot be 0. 0 code means a nil response, and this parameter is required.
// failedResponse property with code 0, means that there is no response in case of Operation failure.
func NewMethodOperation(operation Operation, successResponse, failedResponse Response, operationResultAsBody bool) MethodOperation {
	if successResponse.code == 0 {
		panic("resource: successResponse with code 0 is consider a nil response, and not nil successResponse value is required.")
	}
	return MethodOperation{operation, successResponse, failedResponse, operationResultAsBody}
}
