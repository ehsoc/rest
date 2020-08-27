package resource

type Resource struct {
	Name    string   `json:"name"`
	Type    Type     `json:"type"`
	Path    string   `json:"path"`
	Methods []Method `json:"methods"`
}
