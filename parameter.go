package resource

import (
	"net/http"
	"reflect"
)

type ParameterType string

const (
	BodyParameter     ParameterType = "body"
	FormDataParameter ParameterType = "formData"
	HeaderParameter   ParameterType = "header"
	QueryParameter    ParameterType = "query"
	URIParameter      ParameterType = "uri"
)

type Parameter struct {
	Description string
	Name        string
	GetFunc     func(r *http.Request) string
	HttpType    ParameterType
	Type        reflect.Kind
	Body        interface{}
}
