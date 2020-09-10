package resource

import (
	"io/ioutil"
	"net/http"
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

type Parameter struct {
	Description string
	Name        string
	Getter      Getter
	HTTPType    ParameterType
	Type        reflect.Kind
	Body        interface{}
	Decoder     encdec.Decoder
	Required    bool
}

//NewURIParameter creates a URIParameter Parameter. Required is true by default
func NewURIParameter(name string, tpe reflect.Kind, getter Getter) *Parameter {
	return &Parameter{"", name, getter, URIParameter, tpe, nil, nil, true}
}

//NewHeaderParameter creates a HeaderParameter Parameter. Required is true by default
func NewHeaderParameter(name string, tpe reflect.Kind, getter Getter) *Parameter {
	return &Parameter{"", name, getter, HeaderParameter, tpe, nil, nil, true}
}

//NewQueryParameter creates a QueryParameter Parameter. Required is false by default
func NewQueryParameter(name string, tpe reflect.Kind, getter Getter) *Parameter {
	return &Parameter{"", name, getter, QueryParameter, tpe, nil, nil, false}
}

//NewFormDataParameter creates a FormDataParameter Parameter. Required is false by default
func NewFormDataParameter(name string, tpe reflect.Kind, decoder encdec.Decoder) *Parameter {
	return &Parameter{"", name, GetterFunc(func(r *http.Request) string {
		return r.FormValue(name)
	}), FormDataParameter, tpe, decoder, nil, false}
}

//NewFileParameter creates a FileParameter Parameter. Required is false by default
func NewFileParameter(name string) *Parameter {
	return &Parameter{"", name, GetterFunc(func(r *http.Request) string {
		f, _, err := r.FormFile(name)
		if err != nil {
			return ""
		}
		fileString, err := ioutil.ReadAll(f)
		if err != nil {
			return ""
		}
		return string(fileString)
	}), FileParameter, reflect.String, nil, nil, false}
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
