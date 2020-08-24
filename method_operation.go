package resource

import "net/http"

type MethodOperation struct {
	Operation
	successResponse Response
	failResponse    Response
	GetIdURLParam   func(r *http.Request) string
}

func NewMethodOperation(operation Operation, successResponse, failedResponse Response, getIdFunc func(r *http.Request) string) MethodOperation {
	return MethodOperation{operation, successResponse, failedResponse, getIdFunc}
}
