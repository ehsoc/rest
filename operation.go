package resource

import (
	"github.com/ehsoc/resource/encdec"
)

//Operation defines an operation over a data entity
//Execute function will execute the operation.
type Operation interface {
	Execute(i Input, decoder encdec.Decoder) (interface{}, error)
}

// The OperationFunc type is an adapter to allow the use of
// ordinary functions as Operation. If f is a function
// with the appropriate signature, OperationFunc(f) is a
// Operation that calls f.
type OperationFunc func(i Input, decoder encdec.Decoder) (interface{}, error)

//Execute calls f(body, parameters, decoder)
func (f OperationFunc) Execute(i Input, decoder encdec.Decoder) (interface{}, error) {
	return f(i, decoder)
}
