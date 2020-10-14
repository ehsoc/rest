package resource

import (
	"errors"
	"fmt"
)

var ErrorNoDefaultContentTypeIsSet = errors.New("no default content-type is set")
var ErrorResourceMethodCollition = errors.New("method is already define for this resource")
var ErrorResourceURIParamNoParamFound = errors.New("path must include a parameter name in brackets, like {myParamId}")
var ErrorResourceURIParamMoreThanOne = errors.New("path just must include one parameter name in brackets")
var ErrorRequestBodyNotDefined = errors.New("resource: a request body was not defined.")
var ErrorNilCodeSuccessResponse = errors.New("resource: successResponse with code 0 is considered a nil response. A not nil successResponse value is required.")
var ErrorNilCodeValidationResponse = errors.New("resource: validation Response with code 0 is considered a nil response. A not nil value is required.")

var MessageErrResourceSlashesNotAllowed = "resource: slash found on resource name '%s', slashes are not allowed"
var MessageErrParameterNotDefined = "resource: parameter '%s' not defined"
var MessageErrGetURIParamFunctionNotDefined = "resource: no get uri parameter function is defined in context value InputContextKey(\"uriparamfunc\") for '%v' parameter"
var MessageErrRequiredParameterNotFound = "'%s' is required"
var MessageNotExpectedZeroCodeFailedResponse = "resource: resource '%s' failedResponse code is 0, either the failed response was not expected, or 0 code was set by mistake."

type TypeErrorResourceSlashesNotAllowed struct {
	Errorf
}

type TypeErrorParameterNotDefined struct {
	Errorf
}

type TypeErrorRequiredParameterNotFound struct {
	Errorf
}

type TypeErrorGetURIParamFunctionNotDefined struct {
	Errorf
}

type TypeErrorNotExpectedZeroCodeFailedResponse struct {
	Errorf
}

type Errorf struct {
	format string
	Var    interface{}
}

func (e *Errorf) Error() string {
	return fmt.Sprintf(e.format, e.Var)
}
