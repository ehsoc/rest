package resource

import "net/http"

//MethodOperation is a logic operation that is performed in an entity
type MethodOperation struct {
	//Entity type of the resource
	entity interface{}
	//Logic of the method
	Operation
	//Response if Operation Execute method function returns no error
	successResponse Response
	//Response if Operation Execute method function returns error
	failResponse Response
	//Method to get URL parameter value that represents the ID of the entity.
	//This value will be passed to Execute method in the parameter "id"
	GetIdURLParam func(r *http.Request) string
	//Return the entity (interface{}) returned by Execute method,
	//in the Body of the response instead the Body of successResponse.
	//only the Code of successResponse will be returned in the http response.
	returnEntityOnBodySuccess bool
}

//NewMethodOperation returns a new MethodOperation instance
func NewMethodOperation(entity interface{}, operation Operation, successResponse, failedResponse Response, getIdFunc func(r *http.Request) string, returnEntityOnBodySuccess bool) MethodOperation {
	return MethodOperation{entity, operation, successResponse, failedResponse, getIdFunc, returnEntityOnBodySuccess}

}
