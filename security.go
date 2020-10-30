package resource

type Security struct {
	Type        string
	Name        string
	Description string
	ParameterCollection
	validator                    Validator
	FailedAuthenticationResponse Response
	OAuth2Flows                  []OAuth2Flow
}

const (
	BasicSecurityType  = "basic"
	ApiKeySecurityType = "apiKey"
	OAuth2SecurityType = "oauth2"
)

// NewSecurity creates a new security scheme
func NewSecurity(name string, typ string, securityOperation Validator, failedAuthenticationResponse Response) *Security {
	s := &Security{
		Name:                         name,
		Type:                         typ,
		validator:                    securityOperation,
		FailedAuthenticationResponse: failedAuthenticationResponse,
	}
	s.parameters = make(map[ParameterType]map[string]Parameter)
	return s
}

// NewOAuth2Security creates a new security scheme of OAuth2SecurityType type
func NewOAuth2Security(name string, securityOperation Validator, failedAuthenticationResponse Response) *Security {
	return NewSecurity(name, OAuth2SecurityType, securityOperation, failedAuthenticationResponse)
}

func (o OAuth2Flow) checkScopesMap() {
	if o.Scopes == nil {
		o.Scopes = make(map[string]string)
	}
}

// WithImplicitOAuth2Flow adds a new oauth2 flow of implicit type with the necessary parameters
func (s *Security) WithImplicitOAuth2Flow(authorizationURL string, scopes map[string]string) *Security {
	flow := OAuth2Flow{Name: FlowImplicitType, AuthorizationURL: authorizationURL}
	flow.Scopes = scopes
	flow.checkScopesMap()
	s.OAuth2Flows = append(s.OAuth2Flows, flow)
	return s
}

// WithPasswordOAuth2Flow adds a new oauth2 flow of password type with the necessary parameters
func (s *Security) WithPasswordOAuth2Flow(authorizationURL string, scopes map[string]string) *Security {
	flow := OAuth2Flow{Name: FlowPasswordType, AuthorizationURL: authorizationURL}
	flow.Scopes = scopes
	flow.checkScopesMap()
	s.OAuth2Flows = append(s.OAuth2Flows, flow)
	return s
}

// WithAuthCodeOAuth2Flow adds a new oauth2 flow of authorization_code type with the necessary parameters
func (s *Security) WithAuthCodeOAuth2Flow(authorizationURL, tokenURL string, scopes map[string]string) *Security {
	flow := OAuth2Flow{Name: FlowAuthCodeType, AuthorizationURL: authorizationURL, TokenURL: tokenURL}
	flow.Scopes = scopes
	flow.checkScopesMap()
	s.OAuth2Flows = append(s.OAuth2Flows, flow)
	return s
}

// WithClientCredentialOAuth2Flow adds a new oauth2 flow of client_credential type with the necessary parameters
func (s *Security) WithClientCredentialOAuth2Flow(tokenURL string, scopes map[string]string) *Security {
	flow := OAuth2Flow{Name: FlowClientCredentialType, TokenURL: tokenURL}
	flow.Scopes = scopes
	flow.checkScopesMap()
	s.OAuth2Flows = append(s.OAuth2Flows, flow)
	return s
}

// WithOAuth2Flow adds a new oauth2 flow
func (s *Security) WithOAuth2Flow(flow OAuth2Flow) *Security {
	flow.checkScopesMap()
	s.OAuth2Flows = append(s.OAuth2Flows, flow)
	return s
}

type OAuth2Flow struct {
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
