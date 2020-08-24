package resource_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/ehsoc/resource"
	"github.com/ehsoc/resource/encdec"
)

type ResponseBody struct {
	Code    int
	Message string
}

type OperationStub struct {
	wasCall bool
}

func (o *OperationStub) Execute(id string, query url.Values, entity interface{}) (interface{}, error) {
	o.wasCall = true
	return nil, nil
}

func TestOperations(t *testing.T) {
	successResponse := resource.Response{http.StatusCreated, nil}
	contentTypes := resource.NewHTTPContentTypeSelector()
	contentTypes.Add("application/json", encdec.JSONEncoderDecoder{}, true)
	failResponse := resource.Response{http.StatusInternalServerError, ResponseBody{http.StatusInternalServerError, ""}}
	operation := &OperationStub{}
	mo := resource.NewMethodOperation(operation, successResponse, failResponse, nil)
	method := resource.NewMethod(http.MethodPost, mo, contentTypes)
	request, _ := http.NewRequest(http.MethodPost, "/", nil)
	response := httptest.NewRecorder()
	method.ServeHTTP(response, request)
	if !operation.wasCall {
		t.Errorf("Expecting operation execution.")
	}
	if response.Code != successResponse.Code {
		t.Errorf("got:%d want:%d", response.Code, successResponse.Code)
	}
}
