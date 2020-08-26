package resource

import "net/http"

//MethodOperation is a logic operation that is performed in a entity
type MethodOperation struct {
	Operation
	successResponse       Response
	failResponse          Response
	GetIdURLParam         func(r *http.Request) string
	returnEntityOnSuccess bool
}

//NewMethodOperation returns a new MethodOperation instance
func NewMethodOperation(operation Operation, successResponse, failedResponse Response, getIdFunc func(r *http.Request) string, returnEntityOnSuccess bool) MethodOperation {
	return MethodOperation{operation, successResponse, failedResponse, getIdFunc, returnEntityOnSuccess}

}
