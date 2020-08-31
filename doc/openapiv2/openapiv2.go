package openapiv2

type Resource struct {
	Path    string
	Methods map[string]Method
}

type Method struct {
	Responses   map[string]Response `json:"responses"`
	Summary     string              `json:"summary"`
	Description string              `json:"description"`
	OperationId string              `json:"operationId"`
	Consumes    []string            `json:"consumes"`
	Produces    []string            `json:"produces"`
	Parameters  []Parameter         `json:"parameter"`
}

type Parameter struct {
	Name        string `json:"name,omitempty"`
	In          string `json:"in"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required"`
	Type        string `json:"type,omitempty"`
	Maximum     int    `json:"maximum,omitempty"`
	Minimum     int    `json:"minimum,omitempty"`
	Format      string `json:"format,omitempty"`
}

type Response struct {
}

func (r *Resource) MarshalJSON() ([]byte, error) {
	// return json.Marshal(map[string]interface{}{
	// 	r.Path: r.methods,
	// })
}
