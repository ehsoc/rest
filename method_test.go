package resource_test

import (
	"bytes"
	"encoding/json"
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

type TestResponseBody struct {
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

func (o *OperationStub) Execute(i resource.Input) (interface{}, bool, error) {
	o.wasCall = true
	fbytes, _, _ := i.GetFormFile("file")
	o.FileData = string(fbytes)
	metadata, _ := i.GetFormValue("additionalMetadata")
	o.Metadata = metadata
	car := Car{}
	body, _ := i.GetBody()
	if body != nil && body != http.NoBody {
		i.BodyDecoder.Decode(body, &car)
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
		return nil, false, errors.New("Error")
	}
	fail, _ := i.GetQueryString("fail")
	if fail != "" {
		return nil, false, nil
	}
	return o.Car, true, nil
}

type NegotiatorErrorStub struct {
}

func (n NegotiatorErrorStub) NegotiateEncoder(*http.Request, *resource.Renderers) (mimeType string, encoder encdec.Encoder, err error) {
	return "", nil, errors.New("renderer not available")
}

func (n NegotiatorErrorStub) NegotiateDecoder(*http.Request, *resource.Renderers) (mimeType string, encoder encdec.Decoder, err error) {
	return "", nil, errors.New("renderer not available")
}

type Color struct {
	Name string `json:"name"`
}

type Car struct {
	ID     int     `json:"id"`
	Brand  string  `json:"brand"`
	Colors []Color `json:"colors"`
}

var ISErrorResponse = resource.NewResponse(http.StatusNotFound).WithBody(TestResponseBody{http.StatusNotFound, ""})

func TestOperations(t *testing.T) {
	t.Run("POST no response body", func(t *testing.T) {
		successResponse := resource.NewResponse(http.StatusCreated)
		renderers := resource.NewRenderers()
		renderers.Add("application/json", encdec.JSONEncoderDecoder{}, true)
		failResponse := ISErrorResponse
		operation := &OperationStub{}
		mo := resource.NewMethodOperation(operation, successResponse).WithFailResponse(failResponse)
		method := resource.NewMethod(http.MethodPost, mo, renderers)
		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		response := httptest.NewRecorder()
		method.ServeHTTP(response, request)
		if !operation.wasCall {
			t.Errorf("Expecting operation execution.")
		}
		assertResponseCode(t, response, successResponse.Code())
		enc := encdec.JSONEncoderDecoder{}
		gotResponse := TestResponseBody{}
		enc.Decode(response.Body, &gotResponse)
		body := response.Body.String()
		if body != "" {
			t.Errorf("Not expecting response Body.")
		}
	})
	t.Run("POST response body", func(t *testing.T) {
		responseBody := TestResponseBody{http.StatusCreated, ""}
		successResponse := resource.NewResponse(http.StatusCreated).WithBody(responseBody)
		renderers := resource.NewRenderers()
		renderers.Add("application/json", encdec.JSONEncoderDecoder{}, true)
		failResponse := ISErrorResponse
		operation := &OperationStub{}
		mo := resource.NewMethodOperation(operation, successResponse).WithFailResponse(failResponse)
		method := resource.NewMethod(http.MethodPost, mo, renderers)
		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		response := httptest.NewRecorder()
		method.ServeHTTP(response, request)
		if !operation.wasCall {
			t.Errorf("Expecting operation execution.")
		}
		assertResponseCode(t, response, successResponse.Code())
		enc := encdec.JSONEncoderDecoder{}
		gotResponse := TestResponseBody{}
		enc.Decode(response.Body, &gotResponse)
		if !reflect.DeepEqual(gotResponse, responseBody) {
			t.Errorf("got:%v want:%v", gotResponse, responseBody)
		}
	})
	t.Run("POST Operation Error with query parameter trigger", func(t *testing.T) {
		responseBody := TestResponseBody{http.StatusCreated, ""}
		successResponse := resource.NewResponse(http.StatusCreated).WithBody(responseBody)
		renderers := resource.NewRenderers()
		renderers.Add("application/json", encdec.JSONEncoderDecoder{}, true)
		failResponse := resource.NewResponse(http.StatusFailedDependency).WithBody(TestResponseBody{http.StatusFailedDependency, ""})
		operation := &OperationStub{}
		mo := resource.NewMethodOperation(operation, successResponse).WithFailResponse(failResponse)
		method := resource.NewMethod(http.MethodPost, mo, renderers)
		method.AddParameter(resource.NewQueryParameter("error", reflect.String))
		request, _ := http.NewRequest(http.MethodPost, "/?error=error", nil)
		response := httptest.NewRecorder()
		method.ServeHTTP(response, request)
		if !operation.wasCall {
			t.Errorf("Expecting operation execution.")
		}
		assertResponseCode(t, response, 500)
	})
	t.Run("POST Operation Failed with query parameter trigger, no failed response defined", func(t *testing.T) {
		defer func() {
			err := recover()
			if err != nil {
				if _, ok := err.(*resource.TypeErrorFailResponseNotDefined); !ok {
					t.Errorf("got: %T want: %T", err, &resource.TypeErrorFailResponseNotDefined{})
				}
			}
		}()
		responseBody := TestResponseBody{http.StatusCreated, ""}
		successResponse := resource.NewResponse(http.StatusCreated).WithBody(responseBody)
		renderers := resource.NewRenderers()
		renderers.Add("application/json", encdec.JSONEncoderDecoder{}, true)
		operation := &OperationStub{}
		mo := resource.NewMethodOperation(operation, successResponse)
		method := resource.NewMethod(http.MethodPost, mo, renderers)
		method.AddParameter(resource.NewQueryParameter("fail", reflect.String))
		request, _ := http.NewRequest(http.MethodPost, "/?fail=fail", nil)
		response := httptest.NewRecorder()
		method.ServeHTTP(response, request)
		if !operation.wasCall {
			t.Errorf("Expecting operation execution.")
		}
		//assertResponseCode(t, response, failResponse.Code())
	})
	t.Run("GET Operation Failed with query parameter trigger", func(t *testing.T) {
		successResponse := resource.NewResponse(http.StatusCreated)
		renderers := resource.NewRenderers()
		renderers.Add("application/json", encdec.JSONEncoderDecoder{}, true)
		failResponse := resource.NewResponse(http.StatusNotFound)
		operation := &OperationStub{}
		mo := resource.NewMethodOperation(operation, successResponse).WithFailResponse(failResponse)
		method := resource.NewMethod(http.MethodGet, mo, renderers)
		method.AddParameter(resource.NewQueryParameter("fail", reflect.String))
		request, _ := http.NewRequest(http.MethodPost, "/?fail=fail", nil)
		response := httptest.NewRecorder()
		method.ServeHTTP(response, request)
		if !operation.wasCall {
			t.Errorf("Expecting operation execution.")
		}
		assertResponseCode(t, response, failResponse.Code())
	})
	t.Run("unsupported media response in negotiation in POST with no body", func(t *testing.T) {
		responseBody := TestResponseBody{http.StatusUnsupportedMediaType, "we do not support that"}
		unsupportedMediaResponse := resource.NewResponse(http.StatusUnsupportedMediaType).WithBody(responseBody)
		renderers := resource.NewRenderers()
		renderers.UnsupportedMediaTypeResponse = unsupportedMediaResponse
		renderers.Add("application/json", encdec.JSONEncoderDecoder{}, true)
		operation := &OperationStub{}
		mo := resource.NewMethodOperation(operation, resource.NewResponse(200))
		method := resource.NewMethod(http.MethodPost, mo, renderers)
		method.Negotiator = NegotiatorErrorStub{}
		request, _ := http.NewRequest(http.MethodPost, "/?error=error", nil)
		response := httptest.NewRecorder()
		method.ServeHTTP(response, request)
		assertResponseCode(t, response, unsupportedMediaResponse.Code())
		enc := encdec.JSONEncoderDecoder{}
		gotResponse := TestResponseBody{}
		enc.Decode(response.Body, &gotResponse)
		if !reflect.DeepEqual(gotResponse, unsupportedMediaResponse.Body()) {
			t.Errorf("got:%v want:%v", gotResponse, unsupportedMediaResponse.Body())
		}
	})
	t.Run("unsupported media response in negotiation in POST with body no default type", func(t *testing.T) {
		responseBody := TestResponseBody{http.StatusUnsupportedMediaType, "we do not support that"}
		unsupportedMediaResponse := resource.NewResponse(http.StatusUnsupportedMediaType).WithBody(responseBody)
		renderers := resource.NewRenderers()
		renderers.UnsupportedMediaTypeResponse = unsupportedMediaResponse
		renderers.Add("application/json", encdec.JSONEncoderDecoder{}, false)
		operation := &OperationStub{}
		mo := resource.NewMethodOperation(operation, resource.NewResponse(200))
		method := resource.NewMethod(http.MethodPost, mo, renderers)
		method.Negotiator = NegotiatorErrorStub{}
		request, _ := http.NewRequest(http.MethodPost, "/?error=error", nil)
		response := httptest.NewRecorder()
		method.ServeHTTP(response, request)
		assertResponseCode(t, response, unsupportedMediaResponse.Code())
		if response.Body.String() != "" {
			t.Errorf("Was not expecting body, got:%v", response.Body.String())
		}
	})
	t.Run("unsupported media response decoder negotiation", func(t *testing.T) {
		responseBody := TestResponseBody{http.StatusUnsupportedMediaType, "we do not support that"}
		unsupportedMediaResponse := resource.NewResponse(http.StatusUnsupportedMediaType).WithBody(responseBody)
		renderers := resource.NewRenderers()
		renderers.UnsupportedMediaTypeResponse = unsupportedMediaResponse
		renderers.Add("application/json", encdec.JSONEncoderDecoder{}, true)
		operation := &OperationStub{}
		mo := resource.NewMethodOperation(operation, resource.NewResponse(200))
		method := resource.NewMethod(http.MethodPost, mo, renderers)
		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBufferString("{}"))
		request.Header.Set("Content-Type", "unknown")
		request.Header.Set("Accept", "application/json")
		response := httptest.NewRecorder()
		method.ServeHTTP(response, request)
		assertResponseCode(t, response, unsupportedMediaResponse.Code())
		enc := encdec.JSONEncoderDecoder{}
		gotResponse := TestResponseBody{}
		enc.Decode(response.Body, &gotResponse)
		if !reflect.DeepEqual(gotResponse, unsupportedMediaResponse.Body()) {
			t.Errorf("got:%v want:%v", gotResponse, unsupportedMediaResponse.Body())
		}
	})
	t.Run("GET id return entity on Body response", func(t *testing.T) {
		successResponse := resource.NewResponse(200).WithOperationResultBody(Car{})
		renderers := resource.NewRenderers()
		renderers.Add("application/json", encdec.JSONEncoderDecoder{}, true)
		failResponse := ISErrorResponse
		wantedCar := Car{2, "Fiat", []Color{{"blue"}, {"red"}}}
		operation := &OperationStub{Car: wantedCar}
		mo := resource.NewMethodOperation(operation, successResponse).WithFailResponse(failResponse)
		method := resource.NewMethod(http.MethodPost, mo, renderers)
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
		renderers := resource.NewRenderers()
		renderers.Add("application/json", encdec.JSONEncoderDecoder{}, true)
		failResponse := ISErrorResponse
		wantedCar := Car{200, "Fiat", []Color{{"blue"}, {"red"}}}
		operation := &OperationStub{Car: Car{}}
		mo := resource.NewMethodOperation(operation, successResponse).WithFailResponse(failResponse)
		method := resource.NewMethod(http.MethodPost, mo, renderers).WithRequestBody("", Car{})
		buf := bytes.NewBufferString("")
		_, encoder, _ := renderers.GetDefaultEncoder()
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
		mo := resource.NewMethodOperation(operation, resource.NewResponse(200).WithDescription("successful operation"))
		renderers := resource.NewRenderers()
		renderers.AddEncoder("application/json", encdec.JSONEncoderDecoder{}, true)
		renderers.AddDecoder("multipart/form-data", encdec.XMLEncoderDecoder{}, true)
		method := resource.NewMethod(http.MethodPost, mo, renderers)
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
		successResponse := resource.NewResponse(200).WithOperationResultBody(Car{})
		renderers := resource.NewRenderers()
		renderers.AddEncoder("application/json", encdec.JSONEncoder{}, true)
		failResponse := ISErrorResponse
		wantedCar := Car{2, "Fiat", []Color{{"blue"}, {"red"}}}
		operation := &OperationStub{Car: wantedCar}
		mo := resource.NewMethodOperation(operation, successResponse).WithFailResponse(failResponse)
		method := resource.NewMethod(http.MethodPost, mo, renderers)
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
		m := resource.NewMethod("POST", resource.MethodOperation{}, resource.Renderers{})
		m.AddParameter(resource.NewQueryParameter("myparam", reflect.String))
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
	ct := resource.NewRenderers()
	ct.AddEncoder("application/json", encdec.JSONEncoder{}, true)
	m := resource.NewMethod("POST", resource.NewMethodOperation(nil, resource.NewResponse(200)), ct)
	request, _ := http.NewRequest("POST", "/", nil)
	response := httptest.NewRecorder()
	defer func() { recover() }()
	m.ServeHTTP(response, request)
	t.Errorf("The code did not panic")
}

func TestGetParameters(t *testing.T) {
	t.Run("nil parameters", func(t *testing.T) {
		defer assertNoPanic(t)
		m := resource.NewMethod("POST", resource.MethodOperation{}, resource.Renderers{})
		params := m.GetParameters()
		if len(params) != 0 {
			t.Errorf("got: %v want: %v", len(params), 0)
		}
	})
}

func TestGetEncoderMediaTypes(t *testing.T) {
	ct := resource.NewRenderers()
	ct.AddEncoder("b", encdec.JSONEncoder{}, true)
	ct.AddEncoder("c", encdec.JSONEncoder{}, false)
	ct.AddEncoder("a", encdec.JSONEncoder{}, false)
	m := resource.NewMethod("GET", resource.MethodOperation{}, ct)
	mt := m.GetEncoderMediaTypes()
	number := 3
	if len(mt) != number {
		t.Fatalf("expecting %v elements", number)
	}
	assertStringEqual(t, mt[0], "a")
	assertStringEqual(t, mt[1], "b")
	assertStringEqual(t, mt[2], "c")
}

func TestGetDecoderMediaTypes(t *testing.T) {
	ct := resource.NewRenderers()
	ct.AddDecoder("b", encdec.JSONDecoder{}, true)
	ct.AddDecoder("c", encdec.JSONDecoder{}, false)
	ct.AddDecoder("a", encdec.JSONDecoder{}, false)
	m := resource.NewMethod("GET", resource.MethodOperation{}, ct)
	mt := m.GetDecoderMediaTypes()
	number := 3
	if len(mt) != number {
		t.Fatalf("expecting %v elements", number)
	}
	assertStringEqual(t, mt[0], "a")
	assertStringEqual(t, mt[1], "b")
	assertStringEqual(t, mt[2], "c")
}

func TestWithRequestBody(t *testing.T) {
	car := Car{}
	description := "my request body"
	m := resource.NewMethod("POST", resource.MethodOperation{}, resource.NewRenderers()).
		WithRequestBody("my request body", car)
	if !reflect.DeepEqual(m.RequestBody.Body, car) {
		t.Errorf("got: %v want:%v", m.RequestBody.Body, car)
	}
	assertStringEqual(t, m.RequestBody.Description, description)
}

type MethodValidatorSpy struct {
	called bool
	passed bool
}

func (vs *MethodValidatorSpy) Validate(i resource.Input) error {
	vs.called = true
	r, _ := i.GetQueryString("requiredparam")
	if r == "" {
		return errors.New("requiredparam is required")
	}
	vs.passed = true
	return nil
}

func TestMethodWithValidation(t *testing.T) {
	t.Run("pass validation", func(t *testing.T) {
		v := &MethodValidatorSpy{}
		m := resource.NewMethod("GET", resource.NewMethodOperation(&OperationStub{}, resource.NewResponse(200)).WithFailResponse(resource.NewResponse(404)), mustGetCTS()).
			WithValidation(v, resource.NewResponse(400)).
			WithParameter(resource.NewQueryParameter("requiredparam", reflect.String))
		req, _ := http.NewRequest("GET", "/?requiredparam=something", nil)
		resp := httptest.NewRecorder()
		m.ServeHTTP(resp, req)
		if !v.called {
			t.Errorf("Validator was not called %v", resp)
		}
		if !v.passed {
			t.Errorf("expecting passed to be true")
		}
		assertResponseCode(t, resp, 200)
	})
	t.Run("don't pass validation", func(t *testing.T) {
		v := &MethodValidatorSpy{}
		m := resource.NewMethod("GET", resource.NewMethodOperation(&OperationStub{}, resource.NewResponse(200)).WithFailResponse(resource.NewResponse(404)), mustGetCTS()).
			WithValidation(v, resource.NewResponse(400))
		req, _ := http.NewRequest("GET", "/", nil)
		resp := httptest.NewRecorder()
		m.ServeHTTP(resp, req)
		if !v.called {
			t.Errorf("Validator was not called %v", resp)
		}
		if v.passed {
			t.Errorf("not expecting passed to be true")
		}
		assertResponseCode(t, resp, 400)
	})
}

func TestParameterValidation(t *testing.T) {
	t.Run("pass validation", func(t *testing.T) {
		v := &MethodValidatorSpy{}
		param := resource.NewQueryParameter("requiredparam", reflect.String).WithValidation(v, resource.NewResponse(415))
		m := resource.NewMethod("GET", resource.NewMethodOperation(&OperationStub{}, resource.NewResponse(200)).WithFailResponse(resource.NewResponse(500)), mustGetCTS()).
			WithParameter(param)
		req, _ := http.NewRequest("GET", "/?requiredparam=something", nil)
		resp := httptest.NewRecorder()
		m.ServeHTTP(resp, req)
		if !v.called {
			t.Errorf("Validator was not called %v", resp)
		}
		if !v.passed {
			t.Errorf("expecting passed to be true")
		}
		assertResponseCode(t, resp, 200)
	})
	t.Run("don't pass validation", func(t *testing.T) {
		v := &MethodValidatorSpy{}
		param := resource.NewQueryParameter("requiredparam", reflect.String).WithValidation(v, resource.NewResponse(415))
		m := resource.NewMethod("GET", resource.NewMethodOperation(&OperationStub{}, resource.NewResponse(200)).WithFailResponse(resource.NewResponse(500)), mustGetCTS()).
			WithParameter(param)
		req, _ := http.NewRequest("GET", "/", nil)
		resp := httptest.NewRecorder()
		m.ServeHTTP(resp, req)
		if !v.called {
			t.Errorf("Validator was not called %v", resp)
		}
		if v.passed {
			t.Errorf("not expecting passed to be true")
		}
		assertResponseCode(t, resp, 415)
	})
}

func TestGetResponses(t *testing.T) {
	v := &MethodValidatorSpy{}
	successResponse := resource.NewResponse(200).WithBody(Car{})
	renderers := resource.NewRenderers()
	renderers.Add("application/json", encdec.JSONEncoderDecoder{}, true)
	failResponse := ISErrorResponse
	mo := resource.NewMethodOperation(&OperationStub{}, successResponse).WithFailResponse(failResponse)
	validationResponse := resource.NewResponse(405)
	methodValidationResponse := resource.NewResponse(400)
	method := resource.NewMethod(http.MethodPost, mo, renderers).
		WithParameter(resource.NewQueryParameter("p", reflect.String).WithValidation(v, validationResponse)).
		WithValidation(v, methodValidationResponse)
	responses := method.GetResponses()
	wantResponses := []resource.Response{
		successResponse,
		failResponse,
		methodValidationResponse,
		validationResponse,
	}
	if !reflect.DeepEqual(responses, wantResponses) {
		t.Errorf("got: %v want: %v", responses, wantResponses)
	}
}

func TestMutableResponse(t *testing.T) {
	want := &MutableBodyStub{500, myError, errorMessage}
	mutableResponseBody := &MutableBodyStub{}
	successResponse := resource.NewResponse(415).WithMutableBody(mutableResponseBody)
	m := resource.NewMethod("GET", resource.NewMethodOperation(&OperationStub{}, successResponse).WithFailResponse(resource.NewResponse(404)), mustGetCTS())
	req, _ := http.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()
	m.ServeHTTP(resp, req)
	got := &MutableBodyStub{}
	json.NewDecoder(resp.Body).Decode(got)
	assertResponseCode(t, resp, 415)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got: %#v want: %#v", got, want)
	}
}

type SecurityOperationStub struct {
	authenticateCalled bool
	authorizeCalled    bool
}

func (s *SecurityOperationStub) Authorize(i resource.Input) error {
	s.authorizeCalled = true
	return nil
}

func (s *SecurityOperationStub) Authenticate(i resource.Input) error {
	s.authenticateCalled = true
	return nil
}

type SecurityValidationStub struct {
	called bool
}

func (s *SecurityValidationStub) Validate(i resource.Input) error {
	s.called = true
	return nil
}

func TestSecurity(t *testing.T) {
	t.Run("GET ", func(t *testing.T) {
		successResponse := resource.NewResponse(200)
		renderers := resource.NewRenderers()
		renderers.Add("application/json", encdec.JSONEncoderDecoder{}, true)
		failResponse := resource.NewResponse(404)
		operation := &OperationStub{}
		sv := &SecurityValidationStub{}
		security := resource.NewSecurity("apiKey", resource.ApiKeySecurityType, sv, resource.NewResponse(401))
		mo := resource.NewMethodOperation(operation, successResponse).WithFailResponse(failResponse)
		method := resource.NewMethod(http.MethodGet, mo, renderers).WithSecurity(security)
		method.AddParameter(resource.NewQueryParameter("fail", reflect.String))
		request, _ := http.NewRequest(http.MethodPost, "/?fail=fail", nil)
		response := httptest.NewRecorder()
		method.ServeHTTP(response, request)
		if !operation.wasCall {
			t.Errorf("Expecting operation execution.")
		}
		assertResponseCode(t, response, failResponse.Code())
		assertTrue(t, sv.called)
	})
}
