package resource

import (
	"context"
	"net/http"

	"github.com/ehsoc/resource/encdec"
)

//Method represents a http operation that is performed on a resource.
type Method struct {
	HTTPMethod                string
	request                   Request
	responses                 []Response
	bodyRequiredErrorResponse Response
	methodOperation           MethodOperation
	contentTypeSelector       HTTPContentTypeSelector
	http.Handler
}

//NewMethod returns a Method instance
func NewMethod(httpMethod string, methodOperation MethodOperation, contentTypeSelector HTTPContentTypeSelector) Method {
	m := Method{}
	m.HTTPMethod = httpMethod
	m.methodOperation = methodOperation
	m.contentTypeSelector = contentTypeSelector
	m.newResponse(m.methodOperation.successResponse)
	m.newResponse(m.methodOperation.failResponse)
	m.newResponse(m.contentTypeSelector.unsupportedMediaTypeResponse)
	m.Handler = m.contentTypeMiddleware(http.HandlerFunc(m.mainHandler))
	return m
}

func (m *Method) newResponse(response Response) Response {
	m.responses = append(m.responses, response)
	return response
}

func (m *Method) contentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseContentType, encoder, err := m.contentTypeSelector.NegotiateEncoder(r)
		if err != nil {
			m.writeResponseFallBack(w, m.contentTypeSelector.unsupportedMediaTypeResponse)
			return
		}
		_, decoder, err := m.contentTypeSelector.NegotiateDecoder(r)
		if err != nil {
			writeResponse(w, r.Context(), m.contentTypeSelector.unsupportedMediaTypeResponse)
			return
		}
		ctx := context.WithValue(r.Context(), encoderDecoderContextKey("decoder"), decoder)
		ctx = context.WithValue(ctx, encoderDecoderContextKey("encoder"), encoder)
		w.Header().Add("Content-Type", responseContentType)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Method) writeResponseFallBack(w http.ResponseWriter, response Response) {
	_, encoder, err := m.contentTypeSelector.GetDefaultEncoderDecoder()
	//if no default encdec is set will only return the header code
	if err != nil {
		w.WriteHeader(response.Code)
		return
	}
	write(w, encoder, response)
}

func (m *Method) mainHandler(w http.ResponseWriter, r *http.Request) {
	getIdFunc := m.methodOperation.GetIdURLParam
	id := ""
	if getIdFunc != nil {
		id = m.methodOperation.GetIdURLParam(r)
	}
	entity, err := m.methodOperation.Execute(id, r.URL.Query(), nil)
	if err != nil {
		writeResponse(w, r.Context(), m.methodOperation.failResponse)
		return
	}
	if m.methodOperation.returnEntityOnSuccess {
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
