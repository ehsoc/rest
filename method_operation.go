package resource

import "net/http"

//MethodOperation is a logic operation that is performed in a entity
type MethodOperation struct {
	entity interface{}
	Operation
	successResponse       Response
	failResponse          Response
	GetIdURLParam         func(r *http.Request) string
	returnEntityOnSuccess bool
}

//NewMethodOperation returns a new MethodOperation instance
func NewMethodOperation(entity interface{}, operation Operation, successResponse, failedResponse Response, getIdFunc func(r *http.Request) string, returnEntityOnSuccess bool) MethodOperation {
	return MethodOperation{entity, operation, successResponse, failedResponse, getIdFunc, returnEntityOnSuccess}

}
