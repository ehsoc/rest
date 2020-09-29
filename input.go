package resource

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/ehsoc/resource/httputil"
)

//Input will enforce properties declaration for access the request inputs
type Input struct {
	r *http.Request
	Parameters
	requestBodyParameter interface{}
}

func (i Input) GetURIParam(key string) string {
	//Check param is defined
	i.checkParameterIsDefined(URIParameter, key)
	//Execute param Getter
	getURIValue := i.r.Context().Value(InputContextKey("uriparamfunc"))
	getURIParamFunc, ok := getURIValue.(GetURIParamFunc)
	if !ok {
		fmt.Printf("resource: no get uri parameter function is defined in context value InputContextKey(\"uriparamfunc\") for '%v' parameter\n", key)
	}
	return getURIParamFunc(i.r, key)
}

// func (i Input) GetQuery() ([]string, error) {
// 	//Check param is defined
// 	//Execute param Getter
// }

func (i Input) GetQueryString(key string) string {
	//Check param is defined
	i.checkParameterIsDefined(QueryParameter, key)
	//Execute param Getter
	return i.r.URL.Query().Get(key)
}

func (i Input) GetFormValue(key string) string {
	//Check param is defined
	i.checkParameterIsDefined(FormDataParameter, key)
	//Execute param Getter
	return i.r.FormValue(key)
}

func (i Input) GetFormFile(key string) ([]byte, *multipart.FileHeader, error) {
	//Check param is defined
	i.checkParameterIsDefined(FormDataParameter, key)
	//Execute param Getter
	return httputil.GetFormFile(i.r, key)
}

func (i Input) GetBody() io.ReadCloser {
	//Check param is defined
	if i.requestBodyParameter == nil {
		fmt.Println("resource: no request body parameter is defined")
	}
	//Execute param Getter
	return i.r.Body
}

// func (i Input) GetCookie() ([]string, error) {
// 	//Check param is defined
// 	//Execute param Getter
// }

// func (i Input) GetHeader() (string, error) {
// 	//Check param is defined
// 	//Execute param Getter
// }

func (i *Input) checkParameterIsDefined(pType ParameterType, key string) {
	_, err := i.GetParameter(pType, key)
	if err != nil {
		fmt.Printf("resource: error getting parameter value: %v\n", err)
	}
}
