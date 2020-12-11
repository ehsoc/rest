package rest

import "strings"

// ParameterCollection is a collection of parameters
type ParameterCollection struct {
	parameters map[ParameterType]map[string]Parameter
}

// NewParameterCollection returns a new ParameterCollection
func NewParameterCollection() ParameterCollection {
	p := ParameterCollection{}
	p.parameters = make(map[ParameterType]map[string]Parameter)
	return p
}

// AddParameter adds a new parameter to the collection with the unique composite key by HTTPType and Name properties.
// It will silently override a parameter if the same key is already set.
func (p *ParameterCollection) AddParameter(parameter Parameter) {
	p.checkNilMap()
	// The uri charset is checked here because is the parameter's only point of enter
	if strings.ContainsAny(parameter.Name, URIReservedChar) {
		panic(&ErrorParameterCharNotAllowed{parameter.Name})
	}
	if _, ok := p.parameters[parameter.HTTPType]; !ok {
		p.parameters[parameter.HTTPType] = make(map[string]Parameter)
	}
	p.parameters[parameter.HTTPType][parameter.Name] = parameter
}

func (p *ParameterCollection) checkNilMap() {
	if p.parameters == nil {
		p.parameters = make(map[ParameterType]map[string]Parameter)
	}
}

// Parameters gets the parameter collection.
// The order of the slice elements will not be consistent.
func (p *ParameterCollection) Parameters() []Parameter {
	p.checkNilMap()
	ps := make([]Parameter, 0)

	for _, paramType := range p.parameters {
		for _, param := range paramType {
			ps = append(ps, param)
		}
	}
	return ps
}

// GetParameter gets the parameter of the given ParameterType and name, error if is not found.
func (p *ParameterCollection) GetParameter(paramType ParameterType, name string) (Parameter, error) {
	p.checkNilMap()
	params, ok := p.parameters[paramType]

	if !ok {
		return Parameter{}, &ErrorParameterNotDefined{name}
	}
	if params == nil {
		return Parameter{}, &ErrorParameterNotDefined{name}
	}
	if parameter, ok := params[name]; ok {
		return parameter, nil
	}
	return Parameter{}, &ErrorParameterNotDefined{name}
}
