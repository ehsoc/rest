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
}

func NewMethod(httpMethod string, methodOperation MethodOperation, contentTypeSelector HTTPContentTypeSelector) Method {
	m := Method{}
	m.HTTPMethod = httpMethod
	m.methodOperation = methodOperation
	m.contentTypeSelector = contentTypeSelector
	m.newResponse(m.methodOperation.successResponse)
	m.newResponse(m.methodOperation.failResponse)
	return m
}

func (m *Method) newResponse(response Response) Response {
	m.responses = append(m.responses, response)
	return response
}

func (m *Method) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	encoder.Encode(w, resp.Body)
}
