package resource_test

import (
	"bytes"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"reflect"
	"sort"
	"testing"

	"github.com/ehsoc/resource"
	"github.com/ehsoc/resource/encdec"
)

type ResponseBody struct {
	Code    int
	Message string
}

type OperationStub struct {
	wasCall     bool
	entity      interface{}
	Car         Car
	JsonCarData Car
	FileData    string
	Metadata    string
}

func (o *OperationStub) Execute(i resource.Input, decoder encdec.Decoder) (interface{}, error) {
	o.wasCall = true
	fbytes, _, _ := i.GetFormFile("file")
	o.FileData = string(fbytes)
	metadata, _ := i.GetFormValue("additionalMetadata")
	o.Metadata = metadata
	car := Car{}
	body, _ := i.GetBody()
	if body != nil && body != http.NoBody {
		decoder.Decode(body, &car)
		o.entity = car
	}
	jsonPetData, _ := i.GetFormValue("jsonPetData")
	if jsonPetData != "" {
		buf := bytes.NewBufferString(jsonPetData)
		car := Car{}
		jsonDec := encdec.JSONDecoder{}
		jsonDec.Decode(buf, &car)
		o.JsonCarData = car
	}
	error, _ := i.GetQueryString("error")
	if error != "" {
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

var ISErrorResponse = resource.NewResponse(http.StatusInternalServerError).WithBody(ResponseBody{http.StatusInternalServerError, ""})

func TestOperations(t *testing.T) {
	t.Run("POST no response body", func(t *testing.T) {
		successResponse := resource.NewResponse(http.StatusCreated)
		contentTypes := resource.NewHTTPContentTypeSelector()
		contentTypes.Add("application/json", encdec.JSONEncoderDecoder{}, true)
		failResponse := ISErrorResponse
		operation := &OperationStub{}
		mo := resource.NewMethodOperation(operation, successResponse, failResponse, false)
		method := resource.NewMethod(http.MethodPost, mo, contentTypes)
		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		response := httptest.NewRecorder()
		method.ServeHTTP(response, request)
		if !operation.wasCall {
			t.Errorf("Expecting operation execution.")
		}
		assertResponseCode(t, response, successResponse.Code())
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
		successResponse := resource.NewResponse(http.StatusCreated).WithBody(responseBody)
		contentTypes := resource.NewHTTPContentTypeSelector()
		contentTypes.Add("application/json", encdec.JSONEncoderDecoder{}, true)
		failResponse := ISErrorResponse
		operation := &OperationStub{}
		mo := resource.NewMethodOperation(operation, successResponse, failResponse, false)
		method := resource.NewMethod(http.MethodPost, mo, contentTypes)
		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		response := httptest.NewRecorder()
		method.ServeHTTP(response, request)
		if !operation.wasCall {
			t.Errorf("Expecting operation execution.")
		}
		assertResponseCode(t, response, successResponse.Code())
		enc := encdec.JSONEncoderDecoder{}
		gotResponse := ResponseBody{}
		enc.Decode(response.Body, &gotResponse)
		if !reflect.DeepEqual(gotResponse, responseBody) {
			t.Errorf("got:%v want:%v", gotResponse, responseBody)
		}
	})
	t.Run("POST Operation Failed with query parameter trigger", func(t *testing.T) {
		responseBody := ResponseBody{http.StatusCreated, ""}
		successResponse := resource.NewResponse(http.StatusCreated).WithBody(responseBody)
		contentTypes := resource.NewHTTPContentTypeSelector()
		contentTypes.Add("application/json", encdec.JSONEncoderDecoder{}, true)
		failResponse := resource.NewResponse(http.StatusFailedDependency).WithBody(ResponseBody{http.StatusFailedDependency, ""})
		operation := &OperationStub{}
		mo := resource.NewMethodOperation(operation, successResponse, failResponse, false)
		method := resource.NewMethod(http.MethodPost, mo, contentTypes)
		method.AddParameter(resource.NewQueryParameter("error"))
		request, _ := http.NewRequest(http.MethodPost, "/?error=error", nil)
		response := httptest.NewRecorder()
		method.ServeHTTP(response, request)
		if !operation.wasCall {
			t.Errorf("Expecting operation execution.")
		}
		assertResponseCode(t, response, failResponse.Code())
		enc := encdec.JSONEncoderDecoder{}
		gotResponse := ResponseBody{}
		enc.Decode(response.Body, &gotResponse)
		if !reflect.DeepEqual(gotResponse, failResponse.Body()) {
			t.Errorf("got:%v want:%v", gotResponse, failResponse.Body())
		}
	})
	t.Run("unsupported media response in negotiation in POST with no body", func(t *testing.T) {
		responseBody := ResponseBody{http.StatusUnsupportedMediaType, "we do not support that"}
		unsupportedMediaResponse := resource.NewResponse(http.StatusUnsupportedMediaType).WithBody(responseBody)
		contentTypes := resource.NewHTTPContentTypeSelector()
		contentTypes.UnsupportedMediaTypeResponse = unsupportedMediaResponse
		contentTypes.Add("application/json", encdec.JSONEncoderDecoder{}, true)
		contentTypes.Negotiator = NegotiatorErrorStub{}
		operation := &OperationStub{}
		mo := resource.NewMethodOperation(operation, resource.NewResponse(0), resource.NewResponse(0), false)
		method := resource.NewMethod(http.MethodPost, mo, contentTypes)
		request, _ := http.NewRequest(http.MethodPost, "/?error=error", nil)
		response := httptest.NewRecorder()
		method.ServeHTTP(response, request)
		assertResponseCode(t, response, unsupportedMediaResponse.Code())
		enc := encdec.JSONEncoderDecoder{}
		gotResponse := ResponseBody{}
		enc.Decode(response.Body, &gotResponse)
		if !reflect.DeepEqual(gotResponse, unsupportedMediaResponse.Body()) {
			t.Errorf("got:%v want:%v", gotResponse, unsupportedMediaResponse.Body())
		}
	})
	t.Run("unsupported media response in negotiation in POST with body no default type", func(t *testing.T) {
		responseBody := ResponseBody{http.StatusUnsupportedMediaType, "we do not support that"}
		unsupportedMediaResponse := resource.NewResponse(http.StatusUnsupportedMediaType).WithBody(responseBody)
		contentTypes := resource.NewHTTPContentTypeSelector()
		contentTypes.UnsupportedMediaTypeResponse = unsupportedMediaResponse
		contentTypes.Add("application/json", encdec.JSONEncoderDecoder{}, false)
		contentTypes.Negotiator = NegotiatorErrorStub{}
		operation := &OperationStub{}
		mo := resource.NewMethodOperation(operation, resource.NewResponse(0), resource.NewResponse(0), false)
		method := resource.NewMethod(http.MethodPost, mo, contentTypes)
		request, _ := http.NewRequest(http.MethodPost, "/?error=error", nil)
		response := httptest.NewRecorder()
		method.ServeHTTP(response, request)
		assertResponseCode(t, response, unsupportedMediaResponse.Code())
		if response.Body.String() != "" {
			t.Errorf("Was not expecting body, got:%v", response.Body.String())
		}
	})
	t.Run("unsupported media response decoder negotiation", func(t *testing.T) {
		responseBody := ResponseBody{http.StatusUnsupportedMediaType, "we do not support that"}
		unsupportedMediaResponse := resource.NewResponse(http.StatusUnsupportedMediaType).WithBody(responseBody)
		contentTypes := resource.NewHTTPContentTypeSelector()
		contentTypes.UnsupportedMediaTypeResponse = unsupportedMediaResponse
		contentTypes.Add("application/json", encdec.JSONEncoderDecoder{}, true)
		operation := &OperationStub{}
		mo := resource.NewMethodOperation(operation, resource.NewResponse(0), resource.NewResponse(0), false)
		method := resource.NewMethod(http.MethodPost, mo, contentTypes)
		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBufferString("{}"))
		request.Header.Set("Content-Type", "unknown")
		request.Header.Set("Accept", "application/json")
		response := httptest.NewRecorder()
		method.ServeHTTP(response, request)
		assertResponseCode(t, response, unsupportedMediaResponse.Code())
		enc := encdec.JSONEncoderDecoder{}
		gotResponse := ResponseBody{}
		enc.Decode(response.Body, &gotResponse)
		if !reflect.DeepEqual(gotResponse, unsupportedMediaResponse.Body()) {
			t.Errorf("got:%v want:%v", gotResponse, unsupportedMediaResponse.Body())
		}
	})
	t.Run("GET id return entity on Body response", func(t *testing.T) {
		successResponse := resource.NewResponse(200).WithBody(Car{})
		contentTypes := resource.NewHTTPContentTypeSelector()
		contentTypes.Add("application/json", encdec.JSONEncoderDecoder{}, true)
		failResponse := ISErrorResponse
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
		assertResponseCode(t, response, successResponse.Code())
		enc := encdec.JSONEncoderDecoder{}
		gotResponse := Car{}
		enc.Decode(response.Body, &gotResponse)
		if !reflect.DeepEqual(gotResponse, wantedCar) {
			t.Errorf("got:%v want:%v", gotResponse, wantedCar)
		}
	})
	t.Run("read body", func(t *testing.T) {
		successResponse := resource.NewResponse(201).WithBody(Car{})
		contentTypes := resource.NewHTTPContentTypeSelector()
		contentTypes.Add("application/json", encdec.JSONEncoderDecoder{}, true)
		failResponse := ISErrorResponse
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
		assertResponseCode(t, response, successResponse.Code())
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
		mo := resource.NewMethodOperation(operation, resource.NewResponse(200).WithDescription("successful operation"), resource.NewResponse(0), false)
		contentTypes := resource.NewHTTPContentTypeSelector()
		contentTypes.AddEncoder("application/json", encdec.JSONEncoderDecoder{}, true)
		contentTypes.AddDecoder("multipart/form-data", encdec.XMLEncoderDecoder{}, true)
		method := resource.NewMethod(http.MethodPost, mo, contentTypes)
		method.AddParameter(resource.NewURIParameter("petId", reflect.String))
		method.AddParameter(resource.NewFormDataParameter("additionalMetadata", reflect.String, nil).WithDescription("Additional data to pass to server"))
		method.AddParameter(resource.NewFormDataParameter("jsonPetData", reflect.Struct, encdec.JSONDecoder{}).WithDescription("json format data"))
		method.AddParameter(resource.NewFileParameter("file").WithDescription("file to upload"))
		buf := new(bytes.Buffer)
		w := multipart.NewWriter(buf)
		fileW, _ := w.CreateFormFile("file", "MyFileName.jpg")
		fileData := "filerandomstrings!"
		additionalMetaData := "My Additional Metadata"
		fileW.Write([]byte(fileData))
		fieldW, _ := w.CreateFormField("additionalMetadata")
		fieldW.Write([]byte(additionalMetaData))
		mediaHeader := textproto.MIMEHeader{}
		mediaHeader.Set("Content-Type", "application/json; charset=UTF-8")
		mediaHeader.Set("Content-Disposition", "form-data; name=\"jsonPetData\"")
		jsonPetDataW, _ := w.CreatePart(mediaHeader)
		encoder := encdec.JSONEncoder{}
		wantCar := Car{1, "Subaru", []Color{{"red"}, {"blue"}, {"white"}}}
		encoder.Encode(jsonPetDataW, wantCar)
		w.Close()
		request, _ := http.NewRequest(http.MethodPost, "/", buf)
		request.Header.Set("Content-Type", w.FormDataContentType())
		response := httptest.NewRecorder()
		method.ServeHTTP(response, request)
		if !operation.wasCall {
			t.Errorf("Expecting operation execution.")
		}
		assertResponseCode(t, response, http.StatusOK)
		if operation.FileData != fileData {
			t.Errorf("got :%s want: %s", operation.FileData, fileData)
		}
		if operation.Metadata != additionalMetaData {
			t.Errorf("got :%s want: %s", operation.Metadata, additionalMetaData)
		}
		if !reflect.DeepEqual(operation.JsonCarData, wantCar) {
			t.Errorf("got :%v want: %v", operation.JsonCarData, wantCar)
		}

	})

	t.Run("GET id return entity on Body response only encoder", func(t *testing.T) {
		successResponse := resource.NewResponse(200).WithBody(Car{})
		contentTypes := resource.NewHTTPContentTypeSelector()
		contentTypes.AddEncoder("application/json", encdec.JSONEncoder{}, true)
		failResponse := ISErrorResponse
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
		assertResponseCode(t, response, successResponse.Code())
		enc := encdec.JSONEncoderDecoder{}
		gotResponse := Car{}
		enc.Decode(response.Body, &gotResponse)
		if !reflect.DeepEqual(gotResponse, wantedCar) {
			t.Errorf("got:%v want:%v", gotResponse, wantedCar)
		}
	})
}

func TestAddParameter(t *testing.T) {
	t.Run("nil parameters", func(t *testing.T) {
		defer assertNoPanic(t)
		m := resource.NewMethod("POST", resource.MethodOperation{}, resource.HTTPContentTypeSelector{})
		m.AddParameter(resource.NewQueryParameter("myparam"))
	})
	m := resource.Method{}
	p := resource.Parameter{HTTPType: resource.URIParameter, Name: "id"}
	m.AddParameter(p)
	if !reflect.DeepEqual(m.GetParameters()[0], p) {
		t.Errorf("got: %v want: %v", m.GetParameters()[0], p)
	}
}

func TestWithParameter(t *testing.T) {
	m := &resource.Method{}
	p := resource.Parameter{}
	p2 := resource.Parameter{}
	p.Name = "myParam"
	p2.Name = "myParam2"
	p.HTTPType = resource.FileParameter
	p2.HTTPType = resource.URIParameter
	m.WithParameter(p).WithParameter(p2)
	parameters := m.GetParameters()
	sort.Slice(parameters, func(i, j int) bool {
		return parameters[i].Name < parameters[j].Name
	})
	if parameters[0].Name != p.Name {
		t.Errorf("got: %#v want: %#v", m.GetParameters()[0].Name, p.Name)
	}
	if parameters[1].Name != p2.Name {
		t.Errorf("got: %#v want: %#v", m.GetParameters()[1].Name, p2.Name)
	}
}

func TestChainMethods(t *testing.T) {
	m := &resource.Method{}
	p := resource.Parameter{}
	p2 := resource.Parameter{}
	p.Name = "myParam"
	p2.Name = "myParam2"
	p.HTTPType = resource.FileParameter
	p2.HTTPType = resource.URIParameter
	m.WithParameter(p).WithParameter(p2)
	parameters := m.GetParameters()
	sort.Slice(parameters, func(i, j int) bool {
		return parameters[i].Name < parameters[j].Name
	})
	wantDescription := "my description"
	wantSummary := "my summary"
	if parameters[0].Name != p.Name {
		t.Errorf("got: %v want: %v", parameters[0].Name, p.Name)
	}
	if parameters[1].Name != p2.Name {
		t.Errorf("got: %v want: %v", parameters[1].Name, p2.Name)
	}
	m.WithDescription(wantDescription).WithSummary(wantSummary)
	if m.Description != wantDescription {
		t.Errorf("got: %v want: %v", m.Description, wantDescription)
	}
	if m.Summary != wantSummary {
		t.Errorf("got: %v want: %v", m.Summary, wantSummary)
	}
}

func TestNilOperation(t *testing.T) {
	ct := resource.NewHTTPContentTypeSelector()
	ct.AddEncoder("application/json", encdec.JSONEncoder{}, true)
	m := resource.NewMethod("POST", resource.NewMethodOperation(nil, resource.NewResponse(0), resource.NewResponse(0), false), ct)
	request, _ := http.NewRequest("POST", "/", nil)
	response := httptest.NewRecorder()
	defer func() { recover() }()
	m.ServeHTTP(response, request)
	t.Errorf("The code did not panic")
}

func TestGetParameters(t *testing.T) {
	t.Run("nil parameters", func(t *testing.T) {
		defer assertNoPanic(t)
		m := resource.NewMethod("POST", resource.MethodOperation{}, resource.HTTPContentTypeSelector{})
		params := m.GetParameters()
		if len(params) != 0 {
			t.Errorf("got: %v want: %v", len(params), 0)
		}
	})
}
