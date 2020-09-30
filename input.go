package resource

import (
	"io"
	"mime/multipart"
	"net/http"

	"github.com/ehsoc/resource/httputil"
)

// Input is a set of components that are intendend to help you to check specification-implementation consistency control.
// Input Get methods can help to do this. They will return an error if the parameter or body is not defined.
// Request is the http.Request of the handler, you can use it, but you need to check if the parameter is defined,
// if you want specification-implementation consistency control.
// Parameters is the Parameter collection defined for the method.
// RequestBodyParameter parameter is the Request Body defined for the method.
type Input struct {
	Request              *http.Request
	Parameters           Parameters
	RequestBodyParameter interface{}
}

// GetURIParam gets the URI Param using the InputContextKey("uriparamfunc") context value function.
// If the InputContextKey("uriparamfunc") cintext value is not set will return error.
// If the URI parameter is not set, will return error.
func (i Input) GetURIParam(key string) (string, error) {
	//Check param is defined
	_, err := i.Parameters.GetParameter(URIParameter, key)
	if err != nil {
		return "", err
	}
	getURIValue := i.Request.Context().Value(InputContextKey("uriparamfunc"))
	getURIParamFunc, ok := getURIValue.(func(r *http.Request, key string) string)
	if !ok {
		return "", &TypeErrorGetURIParamFunctionNotDefined{Errorf{MessageErrGetURIParamFunctionNotDefined, key}}
	}
	return getURIParamFunc(i.Request, key), nil
}

// GetQuery gets the request query slice associated to the given key.
// If the parameter is not defined, will return error.
func (i Input) GetQuery(key string) ([]string, error) {
	//Check param is defined
	_, err := i.Parameters.GetParameter(QueryParameter, key)
	if err != nil {
		return nil, err
	}

	return i.Request.URL.Query()[key], nil
}

// GetQueryString gets the first value associated with the given key.
// If there are no values associated with the key, GetQueryString returns the empty string.
// If the parameter is not defined, will return error.
// To access multiple values, use GetQuery
func (i Input) GetQueryString(key string) (string, error) {
	//Check if parameter is defined
	_, err := i.Parameters.GetParameter(QueryParameter, key)
	if err != nil {
		return "", err
	}

	return i.Request.URL.Query().Get(key), nil
}

// GetFormValue gets the first value for the named component of the query.
// FormValue calls FormValue from the standard library.
// If the parameter is not defined, will return error.
func (i Input) GetFormValue(key string) (string, error) {
	//Check param is defined
	_, err := i.Parameters.GetParameter(FormDataParameter, key)
	if err != nil {
		return "", err
	}

	return i.Request.FormValue(key), nil
}

//GetFormFile gets the first file content and header for the provided form key.
// If the parameter is not defined, will return error.
func (i Input) GetFormFile(key string) ([]byte, *multipart.FileHeader, error) {
	//Check param is defined
	_, err := i.Parameters.GetParameter(FileParameter, key)
	if err != nil {
		return nil, nil, err
	}

	return httputil.GetFormFile(i.Request, key)
}

//GetBody returns the request body, error if is not defined.
func (i Input) GetBody() (io.ReadCloser, error) {
	//Check param is defined
	if i.RequestBodyParameter == nil {
		return nil, ErrorRequestBodyNotDefined
	}
	return i.Request.Body, nil
}

// func (i Input) GetCookie() ([]string, error) {
// 	//Check param is defined
//
// }

// func (i Input) GetHeader() (string, error) {
// 	//Check param is defined
//
// }
