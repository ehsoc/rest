package rest

import (
	"errors"
	"fmt"
)

// ErrorNoDefaultContentTypeIsSet error when no default renderer is set
var ErrorNoDefaultContentTypeIsSet = errors.New("no default renderer is set")

// ErrorRequestBodyNotDefined error describes a specification/parameter check error trying to get the request body,
// but it was not declared as parameter.
var ErrorRequestBodyNotDefined = errors.New("resource: a request body was not defined")

var messageErrResourceSlashesNotAllowed = "resource: slash found on resource name '%s', slashes are not allowed"
var messageErrParameterNotDefined = "resource: parameter '%s' not defined"
var messageErrGetURIParamFunctionNotDefined = "resource: no get uri parameter function is defined in context value InputContextKey(\"uriparamfunc\") for '%v' parameter"
var messageErrFailResponseNotDefined = "resource: resource '%s' failedResponse was not defined, but the operation was expecting one"

// TypeErrorResourceSlashesNotAllowed typed error when a slash character is included in the `name` parameter of a `Resource`.
type TypeErrorResourceSlashesNotAllowed struct {
	errorf
}

// TypeErrorParameterNotDefined typed error describes a specification/parameter check error trying to get a parameter,
// but it was not declared.
type TypeErrorParameterNotDefined struct {
	errorf
}

// TypeErrorGetURIParamFunctionNotDefined typed error describes a
// problem when the method GetURIParam has not been declared or correctely set up in the context value key
// by the ServerGenerator implementation.
type TypeErrorGetURIParamFunctionNotDefined struct {
	errorf
}

// TypeErrorFailResponseNotDefined typed error will be trigger when the fail response in the Operation
// was not defined, but the Execute method returns a false value in the success return value.
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
