package rest

import (
	"errors"
	"fmt"
)

// ErrorNoDefaultContentTypeIsSet error when no default content-type is set
var ErrorNoDefaultContentTypeIsSet = errors.New("rest: no default content-type is set")

// ErrorRequestBodyNotDefined error describes a specification/parameter check error trying to get the request body,
// but it was not declared as parameter.
var ErrorRequestBodyNotDefined = errors.New("rest: a request body was not defined")

var msgErrResourceCharNotAllowed = "rest: char not allowed on resource name '%s'"
var msgErrParameterCharNotAllowed = "rest: char not allowed on parameter name '%s'"
var msgErrParameterNotDefined = "rest: parameter '%s' not defined"
var msgErrGetURIParamFunctionNotDefined = "rest: no get uri parameter function is defined in context value InputContextKey(\"uriparamfunc\") for '%v' parameter"
var msgErrFailResponseNotDefined = "rest: resource '%s' failedResponse was not defined, but the operation was expecting one"

// ErrorResourceCharNotAllowed error when a forbidden character is included in the `name` parameter of a `Resource`.
type ErrorResourceCharNotAllowed struct {
	Name string
}

func (e *ErrorResourceCharNotAllowed) Error() string {
	return fmt.Sprintf(msgErrResourceCharNotAllowed, e.Name)
}

// ErrorParameterCharNotAllowed error when a forbidden character is included in the `name` parameter of a `Parameter`.
type ErrorParameterCharNotAllowed struct {
	Name string
}

func (e *ErrorParameterCharNotAllowed) Error() string {
	return fmt.Sprintf(msgErrParameterCharNotAllowed, e.Name)
}

// ErrorParameterNotDefined describes a specification/parameter check error trying to get a parameter,
// but it was not declared.
type ErrorParameterNotDefined struct {
	Name string
}

func (e *ErrorParameterNotDefined) Error() string {
	return fmt.Sprintf(msgErrParameterNotDefined, e.Name)
}

// ErrorGetURIParamFunctionNotDefined describes a
// problem when the method GetURIParam has not been declared or correctely set up in the context value key
// by the ServerGenerator implementation.
type ErrorGetURIParamFunctionNotDefined struct {
	Name string
}

func (e *ErrorGetURIParamFunctionNotDefined) Error() string {
	return fmt.Sprintf(msgErrGetURIParamFunctionNotDefined, e.Name)
}

// ErrorFailResponseNotDefined will be trigger when the fail response in the Operation
// was not defined, but the Execute method returns a false value in the success return value.
type ErrorFailResponseNotDefined struct {
	Name string
}

func (e *ErrorFailResponseNotDefined) Error() string {
	return fmt.Sprintf(msgErrFailResponseNotDefined, e.Name)
}

// AuthError describes an authentication/authorization error.
// Use the following implementations:
// For an authentication failure use the TypeErrorAuthentication error.
// For an authorization failure use TypeErrorAuthorization error
type AuthError interface {
	isAuthorization() bool
	Error() string
}

// ErrorAuthorization is an AuthError implementation respresenting an authorization failure
type ErrorAuthorization struct {
	Message string
}

func (ia ErrorAuthorization) Error() string {
	return ia.Message
}

func (ia ErrorAuthorization) isAuthorization() bool {
	return true
}

// ErrorAuthentication is an AuthError implementation respresenting an authentication failure
type ErrorAuthentication struct {
	Message string
}

func (ia ErrorAuthentication) isAuthorization() bool {
	return false
}

func (ia ErrorAuthentication) Error() string {
	return ia.Message
}
