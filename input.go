package resource

import (
	"io"
	"mime/multipart"
	"net/http"

	"github.com/ehsoc/resource/encdec"
	"github.com/ehsoc/resource/httputil"
)

// Input type is used as Operation Execute method as input.
// Input Get methods can help to check property definition against implementation, returning an error if parameter was not defined.
// Request property provide access to the http.Request pointer.
// Paramters is the collection of defined method's parameters.
// RequestBodyParameter parameter is the Request Body defined for the method.
type Input struct {
	Request              *http.Request
	Parameters           Parameters
	RequestBodyParameter RequestBody
	BodyDecoder          encdec.Decoder
}

// GetURIParam gets the URI Param using the InputContextKey("uriparamfunc") context value.
// If the InputContextKey("uriparamfunc") context value is not set will return an error.
// If the URI parameter is not set, will return an error.
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

// GetHeader gets the request query slice associated to the given key.
// If the parameter is not defined it will return an error.
func (i Input) GetHeader(key string) (string, error) {
	//Check param is defined
	_, err := i.Parameters.GetParameter(HeaderParameter, key)
	if err != nil {
		return "", err
	}
	return i.Request.Header.Get(key), nil
}

// GetQuery gets the request query slice associated to the given key.
// If the parameter is not defined, will return an error.
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

// GetFormValues gets the values associated with the provided key.
// If the parameter is not defined will return error.
// Will also return an error if any error is found getting the values.
func (i Input) GetFormValues(key string) ([]string, error) {
	//Check param is defined
	_, err := i.Parameters.GetParameter(FormDataParameter, key)
	if err != nil {
		return nil, err
	}
	return httputil.GetFormValues(i.Request, key)
}

// GetFormFiles gets all the files of a multipart form with the provided key, in a slice of *multipart.FileHeader
// If the parameter is not defined will return error, as well any other error will returned if a problem
// is find getting the files.
func (i Input) GetFormFiles(key string) ([]*multipart.FileHeader, error) {
	//Check param is defined
	_, err := i.Parameters.GetParameter(FileParameter, key)
	if err != nil {
		return nil, err
	}
	return httputil.GetFiles(i.Request, key)
}

// GetFormFile gets the first file content and header for the provided form key.
// If the parameter is not defined, will return error.
func (i Input) GetFormFile(key string) ([]byte, *multipart.FileHeader, error) {
	//Check param is defined
	_, err := i.Parameters.GetParameter(FileParameter, key)
	if err != nil {
		return nil, nil, err
	}
	return httputil.GetFormFile(i.Request, key)
}

// GetBody returns the request body, error if is not defined.
func (i Input) GetBody() (io.ReadCloser, error) {
	//Check param is defined
	if i.RequestBodyParameter.Body == nil {
		return nil, ErrorRequestBodyNotDefined
	}
	return i.Request.Body, nil
}
