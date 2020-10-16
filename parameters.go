package resource

type Parameters struct {
	parameters map[ParameterType]map[string]Parameter
}

// AddParameter adds a new parameter to the collection with the unique composite key by HTTPType and Name properties.
// It will silently override a parameter if the same key is already set.
func (p *Parameters) AddParameter(parameter Parameter) {
	p.checkNilMap()
	if _, ok := p.parameters[parameter.HTTPType]; !ok {
		p.parameters[parameter.HTTPType] = make(map[string]Parameter)
	}
	p.parameters[parameter.HTTPType][parameter.Name] = parameter
}

func (p *Parameters) checkNilMap() {
	if p.parameters == nil {
		p.parameters = make(map[ParameterType]map[string]Parameter)
	}
}

// GetParameters gets the parameter collection.
// The order of the slice elements will not be consistent.
func (p *Parameters) GetParameters() []Parameter {
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
func (p *Parameters) GetParameter(paramType ParameterType, name string) (Parameter, error) {
	p.checkNilMap()
	params, ok := p.parameters[paramType]
	if !ok {
		return Parameter{}, &TypeErrorParameterNotDefined{Errorf{MessageErrParameterNotDefined, name}}
	}
	if params == nil {
		return Parameter{}, &TypeErrorParameterNotDefined{Errorf{MessageErrParameterNotDefined, name}}
	}
	if parameter, ok := params[name]; ok {
		return parameter, nil
	}
	return Parameter{}, &TypeErrorParameterNotDefined{Errorf{MessageErrParameterNotDefined, name}}
}
