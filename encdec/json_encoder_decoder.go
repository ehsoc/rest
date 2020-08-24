package encdec

//JSONEncoderDecoder is a json EncoderDecoder implementation
//composed by embedding JSONEncoder and JSONDecoder
type JSONEncoderDecoder struct {
	JSONEncoder
	JSONDecoder
}
