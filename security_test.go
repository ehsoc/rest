package resource_test

import (
	"reflect"

	"testing"

	"github.com/ehsoc/resource"
)

func TestNewOAuth2Security(t *testing.T) {
	so := resource.SecurityOperation{}
	want := &resource.Security{
		Type:                resource.OAuth2SecurityType,
		Name:                "myName",
		SecurityOperation:   so,
		ParameterCollection: resource.NewParameterCollection(),
	}
	got := resource.NewOAuth2Security("myName", so)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot:\n %#v \nwant:\n %#v", got, want)
	}
}

func TestNewSecurity(t *testing.T) {
	so := resource.SecurityOperation{}
	want := &resource.Security{
		Type:                resource.OAuth2SecurityType,
		Name:                "myName",
		SecurityOperation:   so,
		ParameterCollection: resource.NewParameterCollection(),
	}
	got := resource.NewSecurity("myName", resource.OAuth2SecurityType, so)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot:\n %#v \nwant:\n %#v", got, want)
	}
}

func TestWithImplicitOAuth2Flow(t *testing.T) {
	so := resource.SecurityOperation{}
	authURL := "http://localhost:7070"
	scopes := map[string]string{"a": "aa"}
	want := &resource.Security{
		Type:                resource.OAuth2SecurityType,
		Name:                "myName",
		SecurityOperation:   so,
		ParameterCollection: resource.NewParameterCollection(),
		OAuth2Flows: []resource.OAuth2Flow{
			resource.OAuth2Flow{
				Name:             resource.FlowImplicitType,
				AuthorizationURL: authURL,
				Scopes:           scopes,
			},
		},
	}
	got := resource.NewSecurity("myName", resource.OAuth2SecurityType, so).WithImplicitOAuth2Flow(authURL, scopes)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot:\n %#v \nwant:\n %#v", got, want)
	}
}

func TestWithPasswordOAuth2Flow(t *testing.T) {
	so := resource.SecurityOperation{}
	tokenURL := "http://localhost:7070"
	scopes := map[string]string{"a": "aa"}
	want := &resource.Security{
		Type:                resource.OAuth2SecurityType,
		Name:                "myName",
		SecurityOperation:   so,
		ParameterCollection: resource.NewParameterCollection(),
		OAuth2Flows: []resource.OAuth2Flow{
			resource.OAuth2Flow{
				Name:     resource.FlowPasswordType,
				TokenURL: tokenURL,
				Scopes:   scopes,
			},
		},
	}
	got := resource.NewSecurity("myName", resource.OAuth2SecurityType, so).WithPasswordOAuth2Flow(tokenURL, scopes)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot:\n %#v \nwant:\n %#v", got, want)
	}
}

func TestWithAuthCodeOAuth2Flow(t *testing.T) {
	so := resource.SecurityOperation{}
	authURL := "http://localhost:7070/auth"
	tokenURL := "http://localhost:7070/token"
	scopes := map[string]string{"a": "aa"}
	want := &resource.Security{
		Type:                resource.OAuth2SecurityType,
		Name:                "myName",
		SecurityOperation:   so,
		ParameterCollection: resource.NewParameterCollection(),
		OAuth2Flows: []resource.OAuth2Flow{
			resource.OAuth2Flow{
				Name:             resource.FlowAuthCodeType,
				AuthorizationURL: authURL,
				TokenURL:         tokenURL,
				Scopes:           scopes,
			},
		},
	}
	got := resource.NewSecurity("myName", resource.OAuth2SecurityType, so).WithAuthCodeOAuth2Flow(authURL, tokenURL, scopes)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot:\n %#v \nwant:\n %#v", got, want)
	}
}

func TestWithClientCredentialOAuth2Flow(t *testing.T) {
	so := resource.SecurityOperation{}
	tokenURL := "http://localhost:7070"
	scopes := map[string]string{"a": "aa"}
	want := &resource.Security{
		Type:                resource.OAuth2SecurityType,
		Name:                "myName",
		SecurityOperation:   so,
		ParameterCollection: resource.NewParameterCollection(),
		OAuth2Flows: []resource.OAuth2Flow{
			resource.OAuth2Flow{
				Name:     resource.FlowClientCredentialType,
				TokenURL: tokenURL,
				Scopes:   scopes,
			},
		},
	}
	got := resource.NewSecurity("myName", resource.OAuth2SecurityType, so).WithClientCredentialOAuth2Flow(tokenURL, scopes)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot:\n %#v \nwant:\n %#v", got, want)
	}
}

func TestWithOAuth2Flow(t *testing.T) {
	so := resource.SecurityOperation{}
	tokenURL := "http://localhost:7070"
	scopes := map[string]string{"a": "aa"}
	flow := resource.OAuth2Flow{
		Name:     resource.FlowClientCredentialType,
		TokenURL: tokenURL,
		Scopes:   scopes,
	}
	want := &resource.Security{
		Type:                resource.OAuth2SecurityType,
		Name:                "myName",
		SecurityOperation:   so,
		ParameterCollection: resource.NewParameterCollection(),
		OAuth2Flows: []resource.OAuth2Flow{
			flow,
		},
	}
	got := resource.NewSecurity("myName", resource.OAuth2SecurityType, so).WithOAuth2Flow(flow)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot:\n %#v \nwant:\n %#v", got, want)
	}
}
