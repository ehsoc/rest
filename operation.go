package resource

import (
	"net/url"
)

//Operation defines an operation over a data entity
//Execute function will execute the operation.
type Operation interface {
	Execute(id string, query url.Values, entity interface{}) (interface{}, error)
}
