package resource

type Security struct {
	Type                         string
	Name                         string
	Description                  string
	Parameter                    Parameter
	validator                    Validator
	failedAuthenticationResponse Response
}

const (
	BasicSecurityType  = "basic"
	ApiKeySecurityType = "apiKey"
	Oauth2SecurityType = "oauth2"
)

func NewSecurity(name string, typ string, securityOperation Validator, failedAuthenticationResponse Response) Security {
	return Security{Name: name, Type: typ, validator: securityOperation, failedAuthenticationResponse: failedAuthenticationResponse}
}
