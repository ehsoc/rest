package resource_test

import (
	"testing"

	"github.com/ehsoc/resource"
)

func TestNewMethodOperation(t *testing.T) {
	resource.NewMethodOperation(&OperationStub{}, resource.NewResponse(0))
}
