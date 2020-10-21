package encdec

// XMLEncoderDecoder is a xml EncoderDecoder implementation
// composed by embedding XMLEncoder and XMLDecoder
type XMLEncoderDecoder struct {
	XMLEncoder
	XMLDecoder
}
