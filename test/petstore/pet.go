package petstore

type Pet struct {
	// id
	ID int64 `json:"id,omitempty"`
	// name
	// Required: true
	Name string `json:"name"`
	// photo urls
	// Required: true
	PhotoUrls []string `json:"photoUrls" xml:"photoUrl"`
	// pet status in the store
	// Enum: [available pending sold]
	Status string `json:"status,omitempty"`
	// tags
	Tags []Tag `json:"tags" xml:"tag"`
}

// Tag tag
//
// swagger:model Tag
type Tag struct {
	// id
	ID int64 `json:"id,omitempty"`
	// name
	Name string `json:"name,omitempty"`
}
