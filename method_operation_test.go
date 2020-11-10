package rest_test

import (
	"testing"

	"github.com/ehsoc/rest"
)

func TestNewMethodOperation(t *testing.T) {
	rest.NewMethodOperation(&OperationStub{}, rest.NewResponse(0))
}
