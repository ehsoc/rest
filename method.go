package rest

import (
	"context"
	"fmt"
	"net/http"
	"sort"

	"github.com/ehsoc/rest/encdec"
)

// Method represents a http operation that is performed on a rest.
type Method struct {
	HTTPMethod      string
	Summary         string
	Description     string
	RequestBody     RequestBody
	MethodOperation MethodOperation
	contentTypes    ContentTypes
	Negotiator
	SecuritySchemes []*Security
	http.Handler
	ParameterCollection
	validation Validation
}

// NewMethod returns a Method instance
func NewMethod(httpMethod string, methodOperation MethodOperation, contentTypes ContentTypes) *Method {
	m := Method{}
	m.HTTPMethod = httpMethod
	m.MethodOperation = methodOperation
	m.contentTypes = contentTypes
	m.Negotiator = DefaultNegotiator{}
	m.parameters = make(map[ParameterType]map[string]Parameter)
	m.Handler = m.negotiationMiddleware(http.HandlerFunc(m.mainHandler))
	return &m
}

func (m *Method) negotiationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseContentType, encoder, err := m.Negotiator.NegotiateEncoder(r, &m.contentTypes)
		if err != nil {
			mutateResponseBody(&m.contentTypes.UnsupportedMediaTypeResponse, nil, false, err)
			m.writeResponseFallBack(w, m.contentTypes.UnsupportedMediaTypeResponse)
			return
		}
		ctx := context.WithValue(r.Context(), EncoderDecoderContextKey("encoder"), encoder)
		ctx = context.WithValue(ctx, ContentTypeContextKey("encoder"), responseContentType)
		decoderContentType, decoder, err := m.Negotiator.NegotiateDecoder(r, &m.contentTypes)
		ctx = context.WithValue(ctx, ContentTypeContextKey("decoder"), decoderContentType)
		if err != nil && r.Body != http.NoBody && r.Body != nil {
			mutateResponseBody(&m.contentTypes.UnsupportedMediaTypeResponse, nil, false, err)
			writeResponse(ctx, w, m.contentTypes.UnsupportedMediaTypeResponse)
			return
		}
		ctx = context.WithValue(ctx, EncoderDecoderContextKey("decoder"), decoder)
		w.Header().Add("Content-Type", responseContentType)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Method) writeResponseFallBack(w http.ResponseWriter, response Response) {
	_, encoder, err := m.contentTypes.GetDefaultEncoder()
	// if no default encdec is set will only return the header code
	if err != nil {
		w.WriteHeader(response.Code())
		return
	}

	write(w, encoder, response)
}

func (m *Method) mainHandler(w http.ResponseWriter, r *http.Request) {
	decoder, ok := r.Context().Value(EncoderDecoderContextKey("decoder")).(encdec.Decoder)
	if !ok {
		// Fallback decoder is a simple string decoder, so we will avoid to pass a nil decoder
		decoder = encdec.TextDecoder{}
	}
	if m.MethodOperation.Operation == nil {
		panic(fmt.Sprintf("resource: resource %s method %s doesn't have an operation.", r.URL.Path, m.HTTPMethod))
	}
	input := Input{r, m.ParameterCollection, m.RequestBody, decoder}

	// Security only if SecuritySchemeEnforcement is true
	if len(m.SecuritySchemes) > 0 {
		passSecurity := false

		var securityFailedResponse Response

		for _, ss := range m.SecuritySchemes {
			if ss.Enforce {
				resp, err := processSecurity(ss, input)
				if err != nil {
					securityFailedResponse = resp
					continue
				}
			}
			passSecurity = true

			break
		}
		if !passSecurity {
			writeResponse(r.Context(), w, securityFailedResponse)
			return
		}
	}

	// Validation
	// Method validation
	if m.validation.Validator != nil {
		err := m.validation.Validate(input)
		if err != nil {
			mutateResponseBody(&m.validation.Response, nil, false, err)
			writeResponse(r.Context(), w, m.validation.Response)
			return
		}
	}
	// Parameter validation
	for _, p := range m.Parameters() {
		if p.validation.Validator != nil && p.validation.Response.code != 0 {
			err := p.validation.Validate(input)
			if err != nil {
				mutateResponseBody(&p.validation.Response, nil, false, err)
				writeResponse(r.Context(), w, p.validation.Response)
				return
			}
		}
	}

	// Operation
	entity, success, err := m.MethodOperation.Execute(input)
	if err != nil {
		errResponse := NewResponse(500)
		mutateResponseBody(&errResponse, entity, success, err)
		writeResponse(r.Context(), w, errResponse)
		return
	}
	if !success {
		if m.MethodOperation.failResponse.disabled {
			panic(&TypeErrorFailResponseNotDefined{errorf{messageErrFailResponseNotDefined, r.URL.Path + " " + m.HTTPMethod}})
		}

		mutateResponseBody(&m.MethodOperation.failResponse, entity, success, err)
		writeResponse(r.Context(), w, m.MethodOperation.failResponse)
		return
	}

	mutateResponseBody(&m.MethodOperation.successResponse, entity, success, err)
	writeResponse(r.Context(), w, m.MethodOperation.successResponse)
}

func processSecurity(ss *Security, input Input) (Response, error) {
	err := ss.Authenticate(input)
	if err != nil {
		if authErr, ok := err.(AuthError); ok {
			if authErr.isAuthorization() {
				return ss.FailedAuthorizationResponse, authErr
			}
			return ss.FailedAuthenticationResponse, authErr
		}
	}
	return Response{}, nil
}

func writeResponse(ctx context.Context, w http.ResponseWriter, resp Response) {
	encoder, ok := ctx.Value(EncoderDecoderContextKey("encoder")).(encdec.Encoder)
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	write(w, encoder, resp)
}

func mutateResponseBody(response *Response, entity interface{}, success bool, err error) {
	if response.MutableResponseBody != nil {
		response.MutableResponseBody.Mutate(entity, success, err)
	}
}

func write(w http.ResponseWriter, encoder encdec.Encoder, resp Response) {
	w.WriteHeader(resp.Code())
	if resp.Body() != nil {
		encoder.Encode(w, resp.Body())
	}
}

// GetEncoderMediaTypes gets a string slice of the method's encoder media types
func (m *Method) GetEncoderMediaTypes() []string {
	mediaTypes := []string{}
	for m := range m.contentTypes.encoderContentTypes {
		mediaTypes = append(mediaTypes, m)
	}
	// Sorting map keys for order consistency
	sort.Strings(mediaTypes)
	return mediaTypes
}

// GetDecoderMediaTypes gets a string slice of the method's decoder media types
func (m *Method) GetDecoderMediaTypes() []string {
	mediaTypes := []string{}
	for m := range m.contentTypes.decoderContentTypes {
		mediaTypes = append(mediaTypes, m)
	}
	// Sorting map keys for consistency
	sort.Strings(mediaTypes)
	return mediaTypes
}

// WithParameter will add a new parameter to the collection with the unique key compose by HTTPType and Name properties.
// It will silently override a parameter if the same key was already set.
func (m *Method) WithParameter(parameter Parameter) *Method {
	m.AddParameter(parameter)
	return m
}

// WithDescription sets the description property
func (m *Method) WithDescription(description string) *Method {
	m.Description = description
	return m
}

// WithSummary sets the summary property
func (m *Method) WithSummary(summary string) *Method {
	m.Summary = summary
	return m
}

// WithRequestBody sets the RequestBody property
func (m *Method) WithRequestBody(description string, body interface{}) *Method {
	m.RequestBody = RequestBody{description, body, true}
	return m
}

// WithValidation sets the validation operation and the response in case of validation error.
func (m *Method) WithValidation(v Validation) *Method {
	m.validation = v
	return m
}

// WithSecurity adds a security to SecuritySchemes slice
func (m *Method) WithSecurity(security *Security) *Method {
	m.SecuritySchemes = append(m.SecuritySchemes, security)
	return m
}

// Responses gets the response collection of the method.
func (m *Method) Responses() []Response {
	responses := make([]Response, 0)
	if !m.MethodOperation.successResponse.disabled {
		responses = append(responses, m.MethodOperation.successResponse)
	}
	if !m.MethodOperation.failResponse.disabled {
		responses = append(responses, m.MethodOperation.failResponse)
	}
	if m.validation.Validator != nil && !m.validation.Response.disabled {
		responses = append(responses, m.validation.Response)
	}

	for _, p := range m.Parameters() {
		if p.validation.Validator != nil && !p.validation.Response.disabled {
			responses = append(responses, p.validation.Response)
		}
	}
	return responses
}
