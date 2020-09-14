package resource

import (
	"context"
	"net/http"
	"net/url"
	"sort"

	"github.com/ehsoc/resource/encdec"
)

//Method represents a http operation that is performed on a resource.
type Method struct {
	HTTPMethod                string
	Summary                   string
	Description               string
	RequestBody               RequestBody
	Responses                 []Response
	bodyRequiredErrorResponse Response
	methodOperation           MethodOperation
	contentTypeSelector       HTTPContentTypeSelector
	http.Handler
	Parameters map[string]Parameter
}

//NewMethod returns a Method instance
func NewMethod(HTTPMethod string, methodOperation MethodOperation, contentTypeSelector HTTPContentTypeSelector) Method {
	m := Method{}
	m.HTTPMethod = HTTPMethod
	m.methodOperation = methodOperation
	m.contentTypeSelector = contentTypeSelector
	m.newResponse(m.methodOperation.successResponse)
	m.newResponse(m.methodOperation.failResponse)
	m.newResponse(m.contentTypeSelector.unsupportedMediaTypeResponse)
	m.Parameters = make(map[string]Parameter)
	m.Handler = m.contentTypeMiddleware(http.HandlerFunc(m.mainHandler))
	return m
}

func (m *Method) newResponse(response Response) Response {
	m.Responses = append(m.Responses, response)
	return response
}

func (m *Method) contentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseContentType, encoder, err := m.contentTypeSelector.NegotiateEncoder(r)
		if err != nil {
			m.writeResponseFallBack(w, m.contentTypeSelector.unsupportedMediaTypeResponse)
			return
		}
		ctx := context.WithValue(r.Context(), encoderDecoderContextKey("encoder"), encoder)
		_, decoder, err := m.contentTypeSelector.NegotiateDecoder(r)
		if err != nil {
			writeResponse(w, ctx, m.contentTypeSelector.unsupportedMediaTypeResponse)
			return
		}
		ctx = context.WithValue(ctx, encoderDecoderContextKey("decoder"), decoder)
		w.Header().Add("Content-Type", responseContentType)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Method) writeResponseFallBack(w http.ResponseWriter, response Response) {
	_, encoder, err := m.contentTypeSelector.GetDefaultEncoder()
	//if no default encdec is set will only return the header code
	if err != nil {
		w.WriteHeader(response.Code)
		return
	}
	write(w, encoder, response)
}

func (m *Method) mainHandler(w http.ResponseWriter, r *http.Request) {
	decoder, ok := r.Context().Value(encoderDecoderContextKey("decoder")).(encdec.Decoder)
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	pValues := url.Values{}
	for name, parameter := range m.Parameters {
		value := parameter.Getter.Get(r)
		pValues.Add(name, value)
	}

	entity, err := m.methodOperation.Execute(r.Body, pValues, decoder)
	if err != nil {
		writeResponse(w, r.Context(), m.methodOperation.failResponse)
		return
	}
	if m.methodOperation.operationResultAsBody {
		writeResponse(w, r.Context(), Response{Code: m.methodOperation.successResponse.Code, Body: entity})
		return
	}
	writeResponse(w, r.Context(), m.methodOperation.successResponse)
}

func writeResponse(w http.ResponseWriter, ctx context.Context, resp Response) {
	encoder, ok := ctx.Value(encoderDecoderContextKey("encoder")).(encdec.Encoder)
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	write(w, encoder, resp)
}

func write(w http.ResponseWriter, encoder encdec.Encoder, resp Response) {
	w.WriteHeader(resp.Code)
	if resp.Body != nil {
		encoder.Encode(w, resp.Body)
	}
}

func (m *Method) getEncoderMediaTypes() []string {
	mediaTypes := []string{}
	for m := range m.contentTypeSelector.encoderContentTypes {
		mediaTypes = append(mediaTypes, m)
	}
	//Sorting map keys for order consistency
	sort.Strings(mediaTypes)
	return mediaTypes
}

func (m *Method) getDecoderMediaTypes() []string {
	mediaTypes := []string{}
	for m := range m.contentTypeSelector.decoderContentTypes {
		mediaTypes = append(mediaTypes, m)
	}
	//Sorting map keys for order consistency
	sort.Strings(mediaTypes)
	return mediaTypes
}

func (m *Method) AddParameter(parameter Parameter) error {
	if m.Parameters == nil {
		m.Parameters = make(map[string]Parameter)
	}
	m.Parameters[parameter.Name] = parameter
	return nil
}
