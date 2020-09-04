package resource

import "errors"

var ErrorNoDefaultContentTypeIsSet = errors.New("no default content-type is set")
var ErrorResourceMethodCollition = errors.New("method is already define for this resource")
var ErrorResourceURIParamNoParamFound = errors.New("path must include a parameter name in brackets, like {myParamId}")
var ErrorResourceURIParamMoreThanOne = errors.New("path just must include one parameter name in brackets")
var ErrorResourceBracketsNotAllowed = errors.New("brackets are not allowed, if you are trying to define a uri paramameter use NewResourceWithURIParam instead")
