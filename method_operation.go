package resource

//MethodOperation is a logic operation that is performed in an entity
type MethodOperation struct {
	//Logic of the method
	Operation
	//Response if Operation Execute method function returns no error
	successResponse Response
	//Response if Operation Execute method function returns error
	failResponse Response
	//Return the entity (interface{}) returned by Execute method,
	//in the Body of the response instead the Body of successResponse.
	//only the Code of successResponse will be returned in the http response.
	operationResultAsBody bool
}

//NewMethodOperation returns a new MethodOperation instance
func NewMethodOperation(operation Operation, successResponse, failedResponse Response, operationResultAsBody bool) MethodOperation {
	return MethodOperation{operation, successResponse, failedResponse, operationResultAsBody}

}
