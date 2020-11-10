package rest_test

import (
	"reflect"

	"testing"

	"github.com/ehsoc/rest"
)

func TestNewOAuth2Security(t *testing.T) {
	so := rest.SecurityOperation{}
	want := &rest.Security{
		Type:                rest.OAuth2SecurityType,
		Name:                "myName",
		SecurityOperation:   so,
		ParameterCollection: rest.NewParameterCollection(),
	}
	got := rest.NewOAuth2Security("myName", so)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot:\n %#v \nwant:\n %#v", got, want)
	}
}

func TestNewSecurity(t *testing.T) {
	so := rest.SecurityOperation{}
	want := &rest.Security{
		Type:                rest.OAuth2SecurityType,
		Name:                "myName",
		SecurityOperation:   so,
		ParameterCollection: rest.NewParameterCollection(),
	}
	got := rest.NewSecurity("myName", rest.OAuth2SecurityType, so)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot:\n %#v \nwant:\n %#v", got, want)
	}
}

func TestWithImplicitOAuth2Flow(t *testing.T) {
	so := rest.SecurityOperation{}
	authURL := "http://localhost:7070"
	scopes := map[string]string{"a": "aa"}
	want := &rest.Security{
		Type:                rest.OAuth2SecurityType,
		Name:                "myName",
		SecurityOperation:   so,
		ParameterCollection: rest.NewParameterCollection(),
		OAuth2Flows: []rest.OAuth2Flow{
			rest.OAuth2Flow{
				Name:             rest.FlowImplicitType,
				AuthorizationURL: authURL,
				Scopes:           scopes,
			},
		},
	}
	got := rest.NewSecurity("myName", rest.OAuth2SecurityType, so).WithImplicitOAuth2Flow(authURL, scopes)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot:\n %#v \nwant:\n %#v", got, want)
	}
}

func TestWithPasswordOAuth2Flow(t *testing.T) {
	so := rest.SecurityOperation{}
	tokenURL := "http://localhost:7070"
	scopes := map[string]string{"a": "aa"}
	want := &rest.Security{
		Type:                rest.OAuth2SecurityType,
		Name:                "myName",
		SecurityOperation:   so,
		ParameterCollection: rest.NewParameterCollection(),
		OAuth2Flows: []rest.OAuth2Flow{
			rest.OAuth2Flow{
				Name:     rest.FlowPasswordType,
				TokenURL: tokenURL,
				Scopes:   scopes,
			},
		},
	}
	got := rest.NewSecurity("myName", rest.OAuth2SecurityType, so).WithPasswordOAuth2Flow(tokenURL, scopes)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot:\n %#v \nwant:\n %#v", got, want)
	}
}

func TestWithAuthCodeOAuth2Flow(t *testing.T) {
	so := rest.SecurityOperation{}
	authURL := "http://localhost:7070/auth"
	tokenURL := "http://localhost:7070/token"
	scopes := map[string]string{"a": "aa"}
	want := &rest.Security{
		Type:                rest.OAuth2SecurityType,
		Name:                "myName",
		SecurityOperation:   so,
		ParameterCollection: rest.NewParameterCollection(),
		OAuth2Flows: []rest.OAuth2Flow{
			rest.OAuth2Flow{
				Name:             rest.FlowAuthCodeType,
				AuthorizationURL: authURL,
				TokenURL:         tokenURL,
				Scopes:           scopes,
			},
		},
	}
	got := rest.NewSecurity("myName", rest.OAuth2SecurityType, so).WithAuthCodeOAuth2Flow(authURL, tokenURL, scopes)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot:\n %#v \nwant:\n %#v", got, want)
	}
}

func TestWithClientCredentialOAuth2Flow(t *testing.T) {
	so := rest.SecurityOperation{}
	tokenURL := "http://localhost:7070"
	scopes := map[string]string{"a": "aa"}
	want := &rest.Security{
		Type:                rest.OAuth2SecurityType,
		Name:                "myName",
		SecurityOperation:   so,
		ParameterCollection: rest.NewParameterCollection(),
		OAuth2Flows: []rest.OAuth2Flow{
			rest.OAuth2Flow{
				Name:     rest.FlowClientCredentialType,
				TokenURL: tokenURL,
				Scopes:   scopes,
			},
		},
	}
	got := rest.NewSecurity("myName", rest.OAuth2SecurityType, so).WithClientCredentialOAuth2Flow(tokenURL, scopes)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot:\n %#v \nwant:\n %#v", got, want)
	}
}

func TestWithOAuth2Flow(t *testing.T) {
	so := rest.SecurityOperation{}
	tokenURL := "http://localhost:7070"
	scopes := map[string]string{"a": "aa"}
	flow := rest.OAuth2Flow{
		Name:     rest.FlowClientCredentialType,
		TokenURL: tokenURL,
		Scopes:   scopes,
	}
	want := &rest.Security{
		Type:                rest.OAuth2SecurityType,
		Name:                "myName",
		SecurityOperation:   so,
		ParameterCollection: rest.NewParameterCollection(),
		OAuth2Flows: []rest.OAuth2Flow{
			flow,
		},
	}
	got := rest.NewSecurity("myName", rest.OAuth2SecurityType, so).WithOAuth2Flow(flow)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot:\n %#v \nwant:\n %#v", got, want)
	}
}
