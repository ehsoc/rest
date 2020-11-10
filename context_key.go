package rest

// EncoderDecoderContextKey is the type used to pass the Encoder/Decoder values the Context of the request.
// The values are set by the Negotiator implementation.
type EncoderDecoderContextKey string

// ContentTypeContextKey is the type used to pass the mime-types values through the Context of the request.
// The values are set by the Negotiator implementation.
type ContentTypeContextKey string

// InputContextKey is the type used to pass the URI Parameter function through the Context of the request.
// The values are set by the GenerateServer method of the API type.
type InputContextKey string
