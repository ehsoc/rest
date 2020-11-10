package rest

// Security contains the authentication and authorization data, and methods.
type Security struct {
	Type        string
	Name        string
	Description string
	ParameterCollection
	SecurityOperation
	Enforce     bool
	OAuth2Flows []OAuth2Flow
}

// SecurityOperation wraps the authentication/authorization method, and the respective fail responses
type SecurityOperation struct {
	Authenticator
	FailedAuthenticationResponse Response
	FailedAuthorizationResponse  Response
}

// Authenticator describes the method for authentication/authorization.
// AuthError represents either a authentication or authorization failure.
// Authentication function should be executed first, then the authorization.
// To indicate an authentication failure return a TypeErrorAuthentication, and
// for an authorization failure TypeErrorAuthorization error type.
// AuthError will be nil when both authentication and authorization are successful
type Authenticator interface {
	Authenticate(Input) AuthError
}

// The AuthenticatorFunc type is an adapter to allow the use of
// ordinary functions as Authenticator. If f is a function
// with the appropriate signature, AuthenticatorFunc(f) is a
// Authenticator that calls f.
type AuthenticatorFunc func(Input) AuthError

// Authenticate calls f(i)
func (f AuthenticatorFunc) Authenticate(i Input) AuthError {
	return f(i)
}

const (
	// BasicSecurityType is the basic authentication security scheme
	BasicSecurityType = "basic"
	// APIKeySecurityType is the API Key authentication security scheme
	APIKeySecurityType = "apiKey"
	// OAuth2SecurityType is the OAuth2 authentication security scheme
	OAuth2SecurityType = "oauth2"
)

// NewSecurity creates a new security scheme
// This security is not enforced by default as is not recommended to turn the enforcement on in production.
// This should serve for specification purposes only, and you should provide the security through middleware implementation.
// If you want to enforce the security you must turn the Enforce property to true.
func NewSecurity(name string, typ string, securityOperation SecurityOperation) *Security {
	s := &Security{
		Name:                name,
		Type:                typ,
		SecurityOperation:   securityOperation,
		ParameterCollection: NewParameterCollection(),
	}

	return s
}

// NewOAuth2Security creates a new security scheme of OAuth2SecurityType type
func NewOAuth2Security(name string, securityOperation SecurityOperation) *Security {
	return NewSecurity(name, OAuth2SecurityType, securityOperation)
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
func (s *Security) WithPasswordOAuth2Flow(tokenURL string, scopes map[string]string) *Security {
	flow := OAuth2Flow{Name: FlowPasswordType, TokenURL: tokenURL}
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

// OAuth2Flow contains the OAuth2 flow or grant information.
type OAuth2Flow struct {
	Name             string
	AuthorizationURL string
	TokenURL         string
	RefreshURL       string
	Scopes           map[string]string
}

const (
	// FlowAuthCodeType is the `authorization` code flow or grant type in a OAuth2 security scheme
	FlowAuthCodeType = "authorization_code"
	// FlowPasswordType is the `password` flow or grant type in a OAuth2 security scheme
	FlowPasswordType = "password"
	// FlowClientCredentialType is the `client credentials` flow or grant type in a OAuth2 security scheme
	FlowClientCredentialType = "client_credentials"
	// FlowImplicitType is the `implicit` flow or grant type in a OAuth2 security scheme
	FlowImplicitType = "implicit"
)
