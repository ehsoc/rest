package resource

type Type struct {
	Name       string
	Properties []Property
}

type Property struct {
	Name string
	Type Type
}
