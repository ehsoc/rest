package resource

import (
	"reflect"

	"github.com/ehsoc/resource/encdec"
)

type ParameterType string

const (
	FormDataParameter ParameterType = "formData"
	FileParameter     ParameterType = "file"
	HeaderParameter   ParameterType = "header"
	QueryParameter    ParameterType = "query"
	URIParameter      ParameterType = "uri"
)

//Parameter a unique parameter is defined by a combination of a HTTPType and Name property
type Parameter struct {
	Description string
	Name        string
	HTTPType    ParameterType
	Type        reflect.Kind
	Body        interface{}
	Decoder     encdec.Decoder
	Required    bool
	CollectionParam
}

//NewURIParameter creates a URIParameter Parameter. Required property is true by default
func NewURIParameter(name string, tpe reflect.Kind) Parameter {
	return Parameter{"", name, URIParameter, tpe, nil, nil, true, CollectionParam{}}
}

//NewHeaderParameter creates a HeaderParameter Parameter. Required property is true by default
func NewHeaderParameter(name string, tpe reflect.Kind) Parameter {
	return Parameter{"", name, HeaderParameter, tpe, nil, nil, true, CollectionParam{}}
}

//NewQueryParameter creates a QueryParameter Parameter. Required property is false by default
func NewQueryParameter(name string) Parameter {
	return Parameter{"", name, QueryParameter, reflect.String, nil, nil, false, CollectionParam{}}
}

//NewQueryArrayParameter creates a QueryParameter Parameter. Required property is false by default
func NewQueryArrayParameter(name string, enumValues []interface{}) Parameter {
	return Parameter{"", name, QueryParameter, reflect.Array, nil, nil, false, CollectionParam{"", enumValues}}
}

//NewFormDataParameter creates a FormDataParameter Parameter. Required property is false by default
func NewFormDataParameter(name string, tpe reflect.Kind, decoder encdec.Decoder) Parameter {
	return Parameter{"", name, FormDataParameter, tpe, nil, decoder, false, CollectionParam{}}
}

//NewFileParameter creates a FileParameter Parameter. Required property is false by default
func NewFileParameter(name string) Parameter {
	return Parameter{"", name, FileParameter, reflect.String, nil, nil, false, CollectionParam{}}
}

//WithDescription chain method set description property
func (p Parameter) WithDescription(description string) Parameter {
	p.Description = description
	return p
}

//WithBody chain method set body property
func (p Parameter) WithBody(body interface{}) Parameter {
	p.Body = body
	return p
}

//AsOptional chain method set Required property to false
func (p Parameter) AsOptional() Parameter {
	p.Required = false
	return p
}

//AsRequired chain method set Required property to true
func (p Parameter) AsRequired() Parameter {
	p.Required = true
	return p
}
