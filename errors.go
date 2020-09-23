package resource

import (
	"errors"
	"fmt"
)

var ErrorNoDefaultContentTypeIsSet = errors.New("no default content-type is set")
var ErrorResourceMethodCollition = errors.New("method is already define for this resource")
var ErrorResourceURIParamNoParamFound = errors.New("path must include a parameter name in brackets, like {myParamId}")
var ErrorResourceURIParamMoreThanOne = errors.New("path just must include one parameter name in brackets")
var ErrorResourceBracketsNotAllowed = errors.New("brackets are not allowed, if you are trying to define a uri paramameter use NewResourceWithURIParam instead")
var FormatErrorResourceSlashesNotAllowed = "resource: slash found on resource name '%s', slashes are not allowed"

type ErrorTypeResourceSlashesNotAllowed struct {
	Errorf
}

type Errorf struct {
	format string
	Var    interface{}
}

func (e *Errorf) Error() string {
	return fmt.Sprintf(e.format, e.Var)
}
