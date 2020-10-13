package resource_test

import (
	"testing"

	"github.com/ehsoc/resource"
)

func TestNewMethodOperation(t *testing.T) {
	defer assertPanicError(t, resource.ErrorNilCodeSuccessResponse)
	resource.NewMethodOperation(&OperationStub{}, resource.NewResponse(0), resource.NewResponse(0))
}
