package resource

import "net/http"

type ParameterType string

const (
	BodyParameter     ParameterType = "body"
	FormDataParameter ParameterType = "formData"
	HeaderParameter   ParameterType = "header"
	QueryParameter    ParameterType = "query"
	URIParameter      ParameterType = "uri"
)

type Parameter struct {
	Name    string
	GetFunc func(r *http.Request) string
	Type    ParameterType
}
