package resource_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
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
	if query.Get("error") != "" {
		return nil, errors.New("Failed")
	}
	return nil, nil
}

func TestOperations(t *testing.T) {
	t.Run("POST no response body", func(t *testing.T) {
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
		AssertResponseCode(t, response, successResponse.Code)
		enc := encdec.JSONEncoderDecoder{}
		gotResponse := ResponseBody{}
		enc.Decode(response.Body, &gotResponse)
		body := response.Body.String()
		if body != "" {
			t.Errorf("Not expecting response Body.")
		}
	})
	t.Run("POST response body", func(t *testing.T) {
		responseBody := ResponseBody{http.StatusCreated, ""}
		successResponse := resource.Response{http.StatusCreated, responseBody}
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
		AssertResponseCode(t, response, successResponse.Code)
		enc := encdec.JSONEncoderDecoder{}
		gotResponse := ResponseBody{}
		enc.Decode(response.Body, &gotResponse)
		if !reflect.DeepEqual(gotResponse, responseBody) {
			t.Errorf("got:%v want:%v", gotResponse, responseBody)
		}
	})
	t.Run("POST Operation Failed", func(t *testing.T) {
		responseBody := ResponseBody{http.StatusCreated, ""}
		successResponse := resource.Response{http.StatusCreated, responseBody}
		contentTypes := resource.NewHTTPContentTypeSelector()
		contentTypes.Add("application/json", encdec.JSONEncoderDecoder{}, true)
		failResponse := resource.Response{http.StatusFailedDependency, ResponseBody{http.StatusFailedDependency, ""}}
		operation := &OperationStub{}
		mo := resource.NewMethodOperation(operation, successResponse, failResponse, nil)
		method := resource.NewMethod(http.MethodPost, mo, contentTypes)
		request, _ := http.NewRequest(http.MethodPost, "/?error=error", nil)
		response := httptest.NewRecorder()
		method.ServeHTTP(response, request)
		if !operation.wasCall {
			t.Errorf("Expecting operation execution.")
		}
		AssertResponseCode(t, response, failResponse.Code)
		enc := encdec.JSONEncoderDecoder{}
		gotResponse := ResponseBody{}
		enc.Decode(response.Body, &gotResponse)
		if !reflect.DeepEqual(gotResponse, failResponse.Body) {
			t.Errorf("got:%v want:%v", gotResponse, failResponse.Body)
		}
	})

}
