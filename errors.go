package resource

import "errors"

var ErrorNoDefaultContentTypeIsSet = errors.New("no default content-type is set")
var ErrorResourceMethodCollition = errors.New("method is already define for this resource")
