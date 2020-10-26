package resource

type Security struct {
	Type        string
	Name        string
	Description string
	Parameters
	validator                    Validator
	FailedAuthenticationResponse Response
	OAuthFlows                   []OAuthFlow
}

const (
	BasicSecurityType  = "basic"
	ApiKeySecurityType = "apiKey"
	Oauth2SecurityType = "oauth2"
)

func NewSecurity(name string, typ string, securityOperation Validator, failedAuthenticationResponse Response) *Security {
	s := &Security{Name: name, Type: typ, validator: securityOperation, FailedAuthenticationResponse: failedAuthenticationResponse}
	s.parameters = make(map[ParameterType]map[string]Parameter)
	return s
}

func (s *Security) WithOAuth2Flow(flow OAuthFlow) *Security {
	s.OAuthFlows = append(s.OAuthFlows, flow)
	return s
}

type OAuthFlow struct {
	Name             string
	AuthorizationURL string
	TokenURL         string
	RefreshURL       string
	Scopes           map[string]string
}

const (
	FlowAuthCodeType         = "authorization_code"
	FlowPasswordType         = "password"
	FlowClientCredentialType = "client_credentials"
	FlowImplicitType         = "implicit"
)
