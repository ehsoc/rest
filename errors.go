package resource

import (
	"errors"
	"fmt"
)

// ErrorNoDefaultContentTypeIsSet error when no default renderer is set
var ErrorNoDefaultContentTypeIsSet = errors.New("no default renderer is set")
var ErrorResourceURIParamNoParamFound = errors.New("path must include a parameter name in brackets, like {myParamId}")
var ErrorResourceURIParamMoreThanOne = errors.New("path just must include one parameter name in brackets")
var ErrorRequestBodyNotDefined = errors.New("resource: a request body was not defined")
var ErrorNilCodeSuccessResponse = errors.New("resource: successResponse with code 0 is considered a nil response. A not nil successResponse value is required")
var ErrorNilCodeValidationResponse = errors.New("resource: validation Response with code 0 is considered a nil response. A not nil value is required")

var messageErrResourceSlashesNotAllowed = "resource: slash found on resource name '%s', slashes are not allowed"
var messageErrParameterNotDefined = "resource: parameter '%s' not defined"
var messageErrGetURIParamFunctionNotDefined = "resource: no get uri parameter function is defined in context value InputContextKey(\"uriparamfunc\") for '%v' parameter"
var messageErrRequiredParameterNotFound = "'%s' is required"
var messageErrFailResponseNotDefined = "resource: resource '%s' failedResponse was not defined, but the operation was expecting one"

type TypeErrorResourceSlashesNotAllowed struct {
	errorf
}

type TypeErrorParameterNotDefined struct {
	errorf
}

type TypeErrorRequiredParameterNotFound struct {
	errorf
}

type TypeErrorGetURIParamFunctionNotDefined struct {
	errorf
}

type TypeErrorFailResponseNotDefined struct {
	errorf
}

type errorf struct {
	format string
	Var    interface{}
}

func (e *errorf) Error() string {
	return fmt.Sprintf(e.format, e.Var)
}

// AuthError describes a authentication/authorization error.
// Use the following implementations:
// For an authentication failure use the TypeErrorAuthentication error
// For an authorization failure use TypeErrorAuthorization error
type AuthError interface {
	isAuthorization() bool
	Error() string
}

// TypeErrorAuthorization is an AuthError implementation respresenting an authorization failure
type TypeErrorAuthorization struct {
	errorf
}

func (ia TypeErrorAuthorization) isAuthorization() bool {
	return true
}

// TypeErrorAuthentication is an AuthError implementation respresenting an authentication failure
type TypeErrorAuthentication struct {
	errorf
}

func (ia TypeErrorAuthentication) isAuthorization() bool {
	return false
}
