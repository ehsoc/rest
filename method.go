package resource

import (
	"context"
	"fmt"
	"net/http"
	"sort"

	"github.com/ehsoc/resource/encdec"
)

//Method represents a http operation that is performed on a resource.
type Method struct {
	HTTPMethod          string
	Summary             string
	Description         string
	RequestBody         RequestBody
	MethodOperation     MethodOperation
	contentTypeSelector HTTPContentTypeSelector
	http.Handler
	Parameters
	validation
}

//NewMethod returns a Method instance
func NewMethod(HTTPMethod string, methodOperation MethodOperation, contentTypeSelector HTTPContentTypeSelector) *Method {
	m := Method{}
	m.HTTPMethod = HTTPMethod
	m.MethodOperation = methodOperation
	m.contentTypeSelector = contentTypeSelector
	m.parameters = make(map[ParameterType]map[string]Parameter)
	m.Handler = m.contentTypeMiddleware(http.HandlerFunc(m.mainHandler))
	return &m
}

func (m *Method) contentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseContentType, encoder, err := m.contentTypeSelector.NegotiateEncoder(r)
		if err != nil {
			m.writeResponseFallBack(w, m.contentTypeSelector.UnsupportedMediaTypeResponse)
			return
		}
		ctx := context.WithValue(r.Context(), EncoderDecoderContextKey("encoder"), encoder)
		ctx = context.WithValue(ctx, ContentTypeContextKey("encoder"), responseContentType)
		decoderContentType, decoder, err := m.contentTypeSelector.NegotiateDecoder(r)
		ctx = context.WithValue(ctx, ContentTypeContextKey("decoder"), decoderContentType)
		if err != nil && r.Body != http.NoBody && r.Body != nil {
			writeResponse(ctx, w, m.contentTypeSelector.UnsupportedMediaTypeResponse)
			return
		}
		ctx = context.WithValue(ctx, EncoderDecoderContextKey("decoder"), decoder)
		w.Header().Add("Content-Type", responseContentType)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Method) writeResponseFallBack(w http.ResponseWriter, response Response) {
	_, encoder, err := m.contentTypeSelector.GetDefaultEncoder()
	//if no default encdec is set will only return the header code
	if err != nil {
		w.WriteHeader(response.Code())
		return
	}
	write(w, encoder, response)
}

func (m *Method) mainHandler(w http.ResponseWriter, r *http.Request) {
	decoder, ok := r.Context().Value(EncoderDecoderContextKey("decoder")).(encdec.Decoder)
	if !ok {
		//Fallback decoder is a simple string decoder, so we will avoid to pass a nil decoder
		decoder = encdec.TextDecoder{}
	}
	if m.MethodOperation.Operation == nil {
		panic(fmt.Sprintf("resource: resource %s method %s doesn't have an operation.", r.URL.Path, m.HTTPMethod))
	}
	input := Input{r, m.Parameters, m.RequestBody, decoder}
	// Validation
	// Method validation
	if m.Validator != nil {
		err := m.Validate(input)
		if err != nil {
			writeResponse(r.Context(), w, m.validationResponse)
			return
		}
	}
	// Parameter validation
	for _, p := range m.GetParameters() {
		if p.Validator != nil && p.validationResponse.code != 0 {
			err := p.Validate(input)
			if err != nil {
				writeResponse(r.Context(), w, p.validationResponse)
				return
			}
		}
	}

	// Operation
	entity, success, err := m.MethodOperation.Execute(input)
	if err != nil {
		writeResponse(r.Context(), w, Response{500, err, ""})
		return
	}
	if !success {
		writeResponse(r.Context(), w, m.MethodOperation.failResponse)
		return
	}
	if m.MethodOperation.operationResultAsBody {
		writeResponse(r.Context(), w, NewResponse(m.MethodOperation.successResponse.Code()).WithBody(entity))
		return
	}
	writeResponse(r.Context(), w, m.MethodOperation.successResponse)
}

func writeResponse(ctx context.Context, w http.ResponseWriter, resp Response) {
	encoder, ok := ctx.Value(EncoderDecoderContextKey("encoder")).(encdec.Encoder)
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	write(w, encoder, resp)
}

func write(w http.ResponseWriter, encoder encdec.Encoder, resp Response) {
	w.WriteHeader(resp.Code())
	if resp.Body() != nil {
		encoder.Encode(w, resp.Body())
	}
}

func (m *Method) GetEncoderMediaTypes() []string {
	mediaTypes := []string{}
	for m := range m.contentTypeSelector.encoderContentTypes {
		mediaTypes = append(mediaTypes, m)
	}
	//Sorting map keys for order consistency
	sort.Strings(mediaTypes)
	return mediaTypes
}

func (m *Method) GetDecoderMediaTypes() []string {
	mediaTypes := []string{}
	for m := range m.contentTypeSelector.decoderContentTypes {
		mediaTypes = append(mediaTypes, m)
	}
	//Sorting map keys for consistency
	sort.Strings(mediaTypes)
	return mediaTypes
}

//WithParameter will add a new parameter to the collection with the unique key of parameter's HTTPType and Name properties.
//It will silently override a parameter if the same key was already set.
func (m *Method) WithParameter(parameter Parameter) *Method {
	m.AddParameter(parameter)
	return m
}

//WithDescription sets description property
func (m *Method) WithDescription(description string) *Method {
	m.Description = description
	return m
}

//WithSummary sets summary property
func (m *Method) WithSummary(summary string) *Method {
	m.Summary = summary
	return m
}

//WithRequestBody sets RequestBody property
func (m *Method) WithRequestBody(description string, body interface{}) *Method {
	m.RequestBody = RequestBody{description, body, true}
	return m
}

//WithValidation sets the validation operation and the response in case of error, this method will override any parameter validation.
func (m *Method) WithValidation(validator Validator, failedValidationResponse Response) *Method {
	m.validation = validation{validator, failedValidationResponse}
	return m
}

// GetResponses gets the valid response collection of the method.
// Responses with 0 code are considered a non valid response.
func (m *Method) GetResponses() []Response {
	responses := make([]Response, 0)
	if m.MethodOperation.successResponse.code != 0 {
		responses = append(responses, m.MethodOperation.successResponse)
	}
	if m.MethodOperation.failResponse.code != 0 {
		responses = append(responses, m.MethodOperation.failResponse)
	}
	if m.Validator != nil && m.validationResponse.code != 0 {
		responses = append(responses, m.validationResponse)
	}
	for _, p := range m.GetParameters() {
		if p.Validator != nil && p.validationResponse.code != 0 {
			responses = append(responses, p.validationResponse)
		}
	}
	return responses
}
