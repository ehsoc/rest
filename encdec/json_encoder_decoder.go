package encdec

// JSONEncoderDecoder is a JSON EncoderDecoder implementation
// composed by embedding JSONEncoder and JSONDecoder types
type JSONEncoderDecoder struct {
	JSONEncoder
	JSONDecoder
}
