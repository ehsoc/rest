package encdec

import "errors"

var ErrorTextDecoderNoString = errors.New("v is not a string")
var ErrorTextDecoderNoValidPointer = errors.New("v is not a valid pointer")
