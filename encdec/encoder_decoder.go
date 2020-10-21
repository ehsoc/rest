package encdec

// EncoderDecoder is the interface that groups the Encode and Decode methods.
type EncoderDecoder interface {
	Encoder
	Decoder
}
