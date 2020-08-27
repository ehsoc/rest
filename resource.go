package resource

type Resource struct {
	Name    string   `json:"name"`
	Path    string   `json:"path"`
	Methods []Method `json:"methods"`
}
