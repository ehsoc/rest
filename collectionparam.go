package rest

// CollectionParam is a subset of properties for array type query parameters
type CollectionParam struct {
	CollectionFormat string
	EnumValues       []interface{}
}
