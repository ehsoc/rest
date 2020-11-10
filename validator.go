package rest

// Validator type implementation will provide a method that validates the given Input.
type Validator interface {
	Validate(Input) error
}

// The ValidatorFunc type is an adapter to allow the use of
// ordinary functions as Validator. If f is a function
// with the appropriate signature, ValidatorFunc(f) is a
// Validator that calls f.
type ValidatorFunc func(i Input) error

// Validate calls f(i)
func (f ValidatorFunc) Validate(i Input) error {
	return f(i)
}
