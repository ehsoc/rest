package resource

import (
	"reflect"
)

type ParameterType string

const (
	BodyParameter     ParameterType = "body"
	FormDataParameter ParameterType = "formData"
	FileParameter     ParameterType = "file"
	HeaderParameter   ParameterType = "header"
	QueryParameter    ParameterType = "query"
	URIParameter      ParameterType = "uri"
)

type Parameter struct {
	Description string
	Name        string
	GetFunc     GetParamFunc
	HTTPType    ParameterType
	Type        reflect.Kind
	Body        interface{}
	Required    bool
}

//NewURIParameter creates a URIParameter Parameter. Required is true by default
func NewURIParameter(name string, tpe reflect.Kind, getFunc GetParamFunc) *Parameter {
	return &Parameter{"", name, getFunc, URIParameter, tpe, nil, true}
}

//NewHeaderParameter creates a HeaderParameter Parameter. Required is true by default
func NewHeaderParameter(name string, tpe reflect.Kind, getFunc GetParamFunc) *Parameter {
	return &Parameter{"", name, getFunc, HeaderParameter, tpe, nil, true}
}

//NewQueryParameter creates a QueryParameter Parameter. Required is false by default
func NewQueryParameter(name string, tpe reflect.Kind, getFunc GetParamFunc) *Parameter {
	return &Parameter{"", name, getFunc, QueryParameter, tpe, nil, false}
}

//NewFormDataParameter creates a FormDataParameter Parameter. Required is false by default
func NewFormDataParameter(name string, tpe reflect.Kind, getFunc GetParamFunc) *Parameter {
	return &Parameter{"", name, getFunc, FormDataParameter, tpe, nil, false}
}

//NewFileParameter creates a FileParameter Parameter. Required is false by default
func NewFileParameter(name string, getFunc GetParamFunc) *Parameter {
	return &Parameter{"", name, getFunc, FileParameter, reflect.Slice, nil, false}
}

//WithDescription set description property
func (p *Parameter) WithDescription(d string) *Parameter {
	p.Description = d
	return p
}

//AsOptional set Required property to false
func (p *Parameter) AsOptional() *Parameter {
	p.Required = false
	return p
}

//AsRequired set Required property to true
func (p *Parameter) AsRequired() *Parameter {
	p.Required = true
	return p
}