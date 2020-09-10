package resource_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
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
	entity  interface{}
	Car     Car
}

func (o *OperationStub) Execute(body io.ReadCloser, params url.Values, decoder encdec.Decoder) (interface{}, error) {
	o.wasCall = true
	car := Car{}
	if body != nil && body != http.NoBody {
		decoder.Decode(body, &car)
		o.entity = car
	}
	if params.Get("error") != "" {
		return nil, errors.New("Failed")
	}
	return o.Car, nil
}

type NegotiatorErrorStub struct {
}

func (n NegotiatorErrorStub) NegotiateEncoder(*http.Request, *resource.HTTPContentTypeSelector) (mimeType string, encoder encdec.Encoder, err error) {
	return "", nil, errors.New("content type not available")
}

func (n NegotiatorErrorStub) NegotiateDecoder(*http.Request, *resource.HTTPContentTypeSelector) (mimeType string, encoder encdec.Decoder, err error) {
	return "", nil, errors.New("content type not available")
}

type Color struct {
	Name string `json:"name"`
}

type Car struct {
	ID     int     `json:"id"`
	Brand  string  `json:"brand"`
	Colors []Color `json:"colors"`
}

func TestOperations(t *testing.T) {
	t.Run("POST no response body", func(t *testing.T) {
		successResponse := resource.Response{http.StatusCreated, nil, ""}
		contentTypes := resource.NewHTTPContentTypeSelector(resource.Response{})
		contentTypes.Add("application/json", encdec.JSONEncoderDecoder{}, true)
		failResponse := resource.Response{http.StatusInternalServerError, ResponseBody{http.StatusInternalServerError, ""}, ""}
		operation := &OperationStub{}
		mo := resource.NewMethodOperation(operation, successResponse, failResponse, false)
		method := resource.NewMethod(http.MethodPost, mo, contentTypes)
		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		response := httptest.NewRecorder()
		method.ServeHTTP(response, request)
		if !operation.wasCall {
			t.Errorf("Expecting operation execution.")
		}
		assertResponseCode(t, response, successResponse.Code)
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
		successResponse := resource.Response{http.StatusCreated, responseBody, ""}
		contentTypes := resource.NewHTTPContentTypeSelector(resource.Response{})
		contentTypes.Add("application/json", encdec.JSONEncoderDecoder{}, true)
		failResponse := resource.Response{http.StatusInternalServerError, ResponseBody{http.StatusInternalServerError, ""}, ""}
		operation := &OperationStub{}
		mo := resource.NewMethodOperation(operation, successResponse, failResponse, false)
		method := resource.NewMethod(http.MethodPost, mo, contentTypes)
		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		response := httptest.NewRecorder()
		method.ServeHTTP(response, request)
		if !operation.wasCall {
			t.Errorf("Expecting operation execution.")
		}
		assertResponseCode(t, response, successResponse.Code)
		enc := encdec.JSONEncoderDecoder{}
		gotResponse := ResponseBody{}
		enc.Decode(response.Body, &gotResponse)
		if !reflect.DeepEqual(gotResponse, responseBody) {
			t.Errorf("got:%v want:%v", gotResponse, responseBody)
		}
	})
	t.Run("POST Operation Failed with query parameter trigger", func(t *testing.T) {
		responseBody := ResponseBody{http.StatusCreated, ""}
		successResponse := resource.Response{http.StatusCreated, responseBody, ""}
		contentTypes := resource.NewHTTPContentTypeSelector(resource.Response{})
		contentTypes.Add("application/json", encdec.JSONEncoderDecoder{}, true)
		failResponse := resource.Response{http.StatusFailedDependency, ResponseBody{http.StatusFailedDependency, ""}, ""}
		operation := &OperationStub{}
		mo := resource.NewMethodOperation(operation, successResponse, failResponse, false)
		method := resource.NewMethod(http.MethodPost, mo, contentTypes)
		method.AddParameter(*resource.NewQueryParameter("error", reflect.String, resource.GetterFunc(func(r *http.Request) string {
			return r.URL.Query().Get("error")
		})))
		request, _ := http.NewRequest(http.MethodPost, "/?error=error", nil)
		response := httptest.NewRecorder()
		method.ServeHTTP(response, request)
		if !operation.wasCall {
			t.Errorf("Expecting operation execution.")
		}
		assertResponseCode(t, response, failResponse.Code)
		enc := encdec.JSONEncoderDecoder{}
		gotResponse := ResponseBody{}
		enc.Decode(response.Body, &gotResponse)
		if !reflect.DeepEqual(gotResponse, failResponse.Body) {
			t.Errorf("got:%v want:%v", gotResponse, failResponse.Body)
		}
	})
	t.Run("unsupported media response in negotiation in POST with no body", func(t *testing.T) {
		responseBody := ResponseBody{http.StatusUnsupportedMediaType, "we do not support that"}
		unsupportedMediaResponse := resource.Response{http.StatusUnsupportedMediaType, responseBody, ""}
		contentTypes := resource.NewHTTPContentTypeSelector(unsupportedMediaResponse)
		contentTypes.Add("application/json", encdec.JSONEncoderDecoder{}, true)
		contentTypes.Negotiator = NegotiatorErrorStub{}
		operation := &OperationStub{}
		mo := resource.NewMethodOperation(operation, resource.Response{}, resource.Response{}, false)
		method := resource.NewMethod(http.MethodPost, mo, contentTypes)
		request, _ := http.NewRequest(http.MethodPost, "/?error=error", nil)
		response := httptest.NewRecorder()
		method.ServeHTTP(response, request)
		assertResponseCode(t, response, unsupportedMediaResponse.Code)
		enc := encdec.JSONEncoderDecoder{}
		gotResponse := ResponseBody{}
		enc.Decode(response.Body, &gotResponse)
		if !reflect.DeepEqual(gotResponse, unsupportedMediaResponse.Body) {
			t.Errorf("got:%v want:%v", gotResponse, unsupportedMediaResponse.Body)
		}
	})
	t.Run("unsupported media response in negotiation in POST with body no default type", func(t *testing.T) {
		responseBody := ResponseBody{http.StatusUnsupportedMediaType, "we do not support that"}
		unsupportedMediaResponse := resource.Response{http.StatusUnsupportedMediaType, responseBody, ""}
		contentTypes := resource.NewHTTPContentTypeSelector(unsupportedMediaResponse)
		contentTypes.Add("application/json", encdec.JSONEncoderDecoder{}, false)
		contentTypes.Negotiator = NegotiatorErrorStub{}
		operation := &OperationStub{}
		mo := resource.NewMethodOperation(operation, resource.Response{}, resource.Response{}, false)
		method := resource.NewMethod(http.MethodPost, mo, contentTypes)
		request, _ := http.NewRequest(http.MethodPost, "/?error=error", nil)
		response := httptest.NewRecorder()
		method.ServeHTTP(response, request)
		assertResponseCode(t, response, unsupportedMediaResponse.Code)
		if response.Body.String() != "" {
			t.Errorf("Was not expecting body, got:%v", response.Body.String())
		}
	})
	t.Run("unsupported media response decoder negotiation", func(t *testing.T) {
		responseBody := ResponseBody{http.StatusUnsupportedMediaType, "we do not support that"}
		unsupportedMediaResponse := resource.Response{http.StatusUnsupportedMediaType, responseBody, ""}
		contentTypes := resource.NewHTTPContentTypeSelector(unsupportedMediaResponse)
		contentTypes.Add("application/json", encdec.JSONEncoderDecoder{}, true)
		operation := &OperationStub{}
		mo := resource.NewMethodOperation(operation, resource.Response{}, resource.Response{}, false)
		method := resource.NewMethod(http.MethodPost, mo, contentTypes)
		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBufferString("{}"))
		request.Header.Set("Content-Type", "unknown")
		request.Header.Set("Accept", "application/json")
		response := httptest.NewRecorder()
		method.ServeHTTP(response, request)
		assertResponseCode(t, response, unsupportedMediaResponse.Code)
		enc := encdec.JSONEncoderDecoder{}
		gotResponse := ResponseBody{}
		enc.Decode(response.Body, &gotResponse)
		if !reflect.DeepEqual(gotResponse, unsupportedMediaResponse.Body) {
			t.Errorf("got:%v want:%v", gotResponse, unsupportedMediaResponse.Body)
		}
	})
	t.Run("GET id return entity on Body response", func(t *testing.T) {
		successResponse := resource.Response{http.StatusOK, Car{}, ""}
		contentTypes := resource.NewHTTPContentTypeSelector(resource.Response{})
		contentTypes.Add("application/json", encdec.JSONEncoderDecoder{}, true)
		failResponse := resource.Response{http.StatusInternalServerError, ResponseBody{http.StatusInternalServerError, ""}, ""}
		wantedCar := Car{2, "Fiat", []Color{{"blue"}, {"red"}}}
		operation := &OperationStub{Car: wantedCar}
		mo := resource.NewMethodOperation(operation, successResponse, failResponse, true)
		method := resource.NewMethod(http.MethodPost, mo, contentTypes)
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		method.ServeHTTP(response, request)
		if !operation.wasCall {
			t.Errorf("Expecting operation execution.")
		}
		assertResponseCode(t, response, successResponse.Code)
		enc := encdec.JSONEncoderDecoder{}
		gotResponse := Car{}
		enc.Decode(response.Body, &gotResponse)
		if !reflect.DeepEqual(gotResponse, wantedCar) {
			t.Errorf("got:%v want:%v", gotResponse, wantedCar)
		}
	})
	t.Run("read body", func(t *testing.T) {
		successResponse := resource.Response{http.StatusCreated, Car{}, ""}
		contentTypes := resource.NewHTTPContentTypeSelector(resource.Response{})
		contentTypes.Add("application/json", encdec.JSONEncoderDecoder{}, true)
		failResponse := resource.Response{http.StatusInternalServerError, ResponseBody{http.StatusInternalServerError, ""}, ""}
		wantedCar := Car{200, "Fiat", []Color{{"blue"}, {"red"}}}
		operation := &OperationStub{Car: Car{}}
		mo := resource.NewMethodOperation(operation, successResponse, failResponse, false)
		method := resource.NewMethod(http.MethodPost, mo, contentTypes)
		buf := bytes.NewBufferString("")
		_, encoder, _ := contentTypes.GetDefaultEncoder()
		encoder.Encode(buf, wantedCar)
		request, _ := http.NewRequest(http.MethodPost, "/", buf)
		request.Header.Set("content-type", "application/json")
		response := httptest.NewRecorder()
		method.ServeHTTP(response, request)
		if !operation.wasCall {
			t.Errorf("Expecting operation execution.")
		}
		assertResponseCode(t, response, successResponse.Code)
		gotOpCar, ok := operation.entity.(Car)
		if !ok {
			t.Fatalf("Expecting type Car, got: %T", operation.entity)
		}
		if !reflect.DeepEqual(gotOpCar, wantedCar) {
			t.Errorf("got:%v want:%v", gotOpCar, wantedCar)
		}
	})

	t.Run("read multipart/form-data", func(t *testing.T) {
		operation := &OperationStub{}
		mo := resource.NewMethodOperation(operation, resource.Response{200, nil, "successful operation"}, resource.Response{}, false)
		contentTypes := resource.NewHTTPContentTypeSelector(resource.Response{http.StatusUnsupportedMediaType, nil, ""})
		contentTypes.AddEncoder("application/json", encdec.JSONEncoderDecoder{}, true)
		contentTypes.AddDecoder("multipart/form-data", encdec.XMLEncoderDecoder{}, true)
		method := resource.NewMethod(http.MethodPost, mo, contentTypes)
		method.AddParameter(*resource.NewURIParameter("petId", reflect.String, resource.GetterFunc(func(r *http.Request) string { return "" })))
		method.AddParameter(*resource.NewFormDataParameter("additionalMetadata", reflect.String, encdec.JSONDecoder{}).WithDescription("Additional data to pass to server"))
		method.AddParameter(*resource.NewFileParameter("file").WithDescription("file to upload"))
		buf := new(bytes.Buffer)
		w := multipart.NewWriter(buf)
		fileW, _ := w.CreateFormFile("file", "MyFileName.jpg")
		fieldW, _ := w.CreateFormField("additionalMetadata")
		fileW.Write([]byte("randomstrings!"))
		fieldW.Write([]byte("My Additional Metadata"))
		request, _ := http.NewRequest(http.MethodPost, "/", buf)
		request.Header.Set("Content-Type", w.FormDataContentType())
		response := httptest.NewRecorder()
		method.ServeHTTP(response, request)
		if !operation.wasCall {
			t.Errorf("Expecting operation execution.")
		}
		assertResponseCode(t, response, http.StatusOK)
		fmt.Println(response)
	})
}

func TestAddParameter(t *testing.T) {
	m := resource.Method{}
	p := resource.Parameter{}
	p.Name = "myParam"
	m.AddParameter(p)
	if !reflect.DeepEqual(m.Parameters[p.Name], p) {
		t.Errorf("got: %v want: %v", m.Parameters[p.Name], p)
	}
}
