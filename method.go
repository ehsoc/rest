package resource

import (
	"context"
	"net/http"

	"github.com/ehsoc/resource/encdec"
)

type Method struct {
	HTTPMethod                string
	request                   Request
	responses                 []Response
	bodyRequiredErrorResponse Response
	methodOperation           MethodOperation
	contentTypeSelector       HTTPContentTypeSelector
	http.Handler
}

func NewMethod(httpMethod string, methodOperation MethodOperation, contentTypeSelector HTTPContentTypeSelector) Method {
	m := Method{}
	m.HTTPMethod = httpMethod
	m.methodOperation = methodOperation
	m.contentTypeSelector = contentTypeSelector
	m.newResponse(m.methodOperation.successResponse)
	m.newResponse(m.methodOperation.failResponse)
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
			http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
			return
		}
		_, decoder, err := m.contentTypeSelector.NegotiateDecoder(r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
			return
		}
		ctx := context.WithValue(r.Context(), encoderDecoderContextKey("decoder"), decoder)
		ctx = context.WithValue(ctx, encoderDecoderContextKey("encoder"), encoder)
		w.Header().Add("Content-Type", responseContentType)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Method) mainHandler(w http.ResponseWriter, r *http.Request) {
	getIdFunc := m.methodOperation.GetIdURLParam
	id := ""
	if getIdFunc != nil {
		id = m.methodOperation.GetIdURLParam(r)
	}
	_, err := m.methodOperation.Execute(id, r.URL.Query(), nil)
	if err != nil {
		writeResponse(w, r.Context(), m.methodOperation.failResponse)
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
	w.WriteHeader(resp.Code)
	if resp.Body != nil {
		encoder.Encode(w, resp.Body)
	}
}
