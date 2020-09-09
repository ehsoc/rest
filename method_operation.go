package resource

//MethodOperation is a logic operation that is performed in an entity
type MethodOperation struct {
	//entity is the entity
	entity interface{}
	//Logic of the method
	Operation
	//Response if Operation Execute method function returns no error
	successResponse Response
	//Response if Operation Execute method function returns error
	failResponse Response
	//Return the entity (interface{}) returned by Execute method,
	//in the Body of the response instead the Body of successResponse.
	//only the Code of successResponse will be returned in the http response.
	returnEntityOnBodySuccess bool
	//The entity is expected to be send in the request body
	entityOnRequestBody bool
}

//NewMethodOperation returns a new MethodOperation instance
func NewMethodOperation(entity interface{}, operation Operation, successResponse, failedResponse Response, returnEntityOnBodySuccess bool, entityOnRequestBody bool) MethodOperation {
	return MethodOperation{entity, operation, successResponse, failedResponse, returnEntityOnBodySuccess, entityOnRequestBody}

}
