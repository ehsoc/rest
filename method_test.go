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

type NegotiatorErrorStub struct {
}

func (n NegotiatorErrorStub) NegotiateEncoder(*http.Request, *resource.HTTPContentTypeSelector) (mimeType string, encoder encdec.Encoder, err error) {
	return "", nil, errors.New("content type not available")
}

func (n NegotiatorErrorStub) NegotiateDecoder(*http.Request, *resource.HTTPContentTypeSelector) (mimeType string, encoder encdec.Decoder, err error) {
	return "", nil, errors.New("content type not available")
}

func TestOperations(t *testing.T) {
	t.Run("POST no response body", func(t *testing.T) {
		successResponse := resource.Response{http.StatusCreated, nil}
		contentTypes := resource.NewHTTPContentTypeSelector(resource.Response{})
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
		contentTypes := resource.NewHTTPContentTypeSelector(resource.Response{})
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
		contentTypes := resource.NewHTTPContentTypeSelector(resource.Response{})
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
	t.Run("unsupported media response in negotiation in POST with body", func(t *testing.T) {
		responseBody := ResponseBody{http.StatusUnsupportedMediaType, "we do not support that"}
		unsupportedMediaResponse := resource.Response{http.StatusUnsupportedMediaType, responseBody}
		contentTypes := resource.NewHTTPContentTypeSelector(unsupportedMediaResponse)
		contentTypes.Add("application/json", encdec.JSONEncoderDecoder{}, true)
		contentTypes.Negotiator = NegotiatorErrorStub{}
		operation := &OperationStub{}
		mo := resource.NewMethodOperation(operation, resource.Response{}, resource.Response{}, nil)
		method := resource.NewMethod(http.MethodPost, mo, contentTypes)
		request, _ := http.NewRequest(http.MethodPost, "/?error=error", nil)
		response := httptest.NewRecorder()
		method.ServeHTTP(response, request)
		AssertResponseCode(t, response, unsupportedMediaResponse.Code)
		enc := encdec.JSONEncoderDecoder{}
		gotResponse := ResponseBody{}
		enc.Decode(response.Body, &gotResponse)
		if !reflect.DeepEqual(gotResponse, unsupportedMediaResponse.Body) {
			t.Errorf("got:%v want:%v", gotResponse, unsupportedMediaResponse.Body)
		}
	})
	t.Run("unsupported media response in negotiation in POST with body no default type", func(t *testing.T) {
		responseBody := ResponseBody{http.StatusUnsupportedMediaType, "we do not support that"}
		unsupportedMediaResponse := resource.Response{http.StatusUnsupportedMediaType, responseBody}
		contentTypes := resource.NewHTTPContentTypeSelector(unsupportedMediaResponse)
		contentTypes.Add("application/json", encdec.JSONEncoderDecoder{}, false)
		contentTypes.Negotiator = NegotiatorErrorStub{}
		operation := &OperationStub{}
		mo := resource.NewMethodOperation(operation, resource.Response{}, resource.Response{}, nil)
		method := resource.NewMethod(http.MethodPost, mo, contentTypes)
		request, _ := http.NewRequest(http.MethodPost, "/?error=error", nil)
		response := httptest.NewRecorder()
		method.ServeHTTP(response, request)
		AssertResponseCode(t, response, unsupportedMediaResponse.Code)
		if response.Body.String() != "" {
			t.Errorf("Was not expecting body, got:%v", response.Body.String())
		}
	})
}
