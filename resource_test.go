package rest_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"testing"

	"github.com/ehsoc/rest"
)

var moTest = rest.NewMethodOperation(&OperationStub{}, rest.NewResponse(200)).WithFailResponse(rest.NewResponse(404))

func TestAddMethod(t *testing.T) {
	ct := rest.NewContentTypes()
	t.Run("get method", func(t *testing.T) {
		r := rest.NewResource("pet")
		m := rest.NewMethod("GET", moTest, ct)
		r.AddMethod(m)
		if len(r.Methods()) != 1 {
			t.Errorf("expecting on method")
		}
		getMethod := r.Methods()[0]
		if getMethod.HTTPMethod != http.MethodGet {
			t.Errorf("got: %v want: %v", getMethod.HTTPMethod, http.MethodGet)
		}
		if !reflect.DeepEqual(getMethod.MethodOperation, moTest) {
			t.Errorf("got: %v want: %v", getMethod.MethodOperation, moTest)
		}
	})
	t.Run("methods map nil", func(t *testing.T) {
		defer assertNoPanic(t)
		r := rest.Resource{}
		r.Get(rest.MethodOperation{}, rest.ContentTypes{})
	})
	t.Run("methods map nil with constructor", func(t *testing.T) {
		defer assertNoPanic(t)
		r := rest.NewResource("pet")
		r.Get(rest.MethodOperation{}, rest.ContentTypes{})
	})
	t.Run("override existing method", func(t *testing.T) {
		description := "second get"
		r := rest.NewResource("")
		r.Get(rest.MethodOperation{}, rest.ContentTypes{}).WithDescription("first get")
		r.Post(rest.MethodOperation{}, rest.ContentTypes{})
		r.Get(rest.MethodOperation{}, rest.ContentTypes{}).WithDescription(description)
		if len(r.Methods()) != 2 {
			t.Fatalf("expecting 2 methods, got: %v", len(r.Methods()))
		}
		for _, m := range r.Methods() {
			if m.HTTPMethod == "GET" {
				if m.Description != "second get" {
					t.Errorf("got: %v want: %v", m.Description, description)
				}
			}
		}
	})
}

func TestGet(t *testing.T) {
	ct := rest.NewContentTypes()
	t.Run("get method", func(t *testing.T) {
		r := rest.NewResource("pet")
		r.Get(moTest, ct)
		if len(r.Methods()) != 1 {
			t.Errorf("expecting on method")
		}
		getMethod := r.Methods()[0]
		if getMethod.HTTPMethod != http.MethodGet {
			t.Errorf("got: %v want: %v", getMethod.HTTPMethod, http.MethodGet)
		}
		if !reflect.DeepEqual(getMethod.MethodOperation, moTest) {
			t.Errorf("got: %v want: %v", getMethod.MethodOperation, moTest)
		}
	})
	t.Run("post method", func(t *testing.T) {
		r := rest.NewResource("pet")
		r.Post(moTest, ct)
		if len(r.Methods()) != 1 {
			t.Errorf("expecting on method")
		}
		getMethod := r.Methods()[0]
		if getMethod.HTTPMethod != http.MethodPost {
			t.Errorf("got: %v want: %v", getMethod.HTTPMethod, http.MethodPost)
		}
		if !reflect.DeepEqual(getMethod.MethodOperation, moTest) {
			t.Errorf("got: %v want: %v", getMethod.MethodOperation, moTest)
		}
	})
	t.Run("options put", func(t *testing.T) {
		r := rest.NewResource("pet")
		r.Put(moTest, ct)
		if len(r.Methods()) != 1 {
			t.Errorf("expecting on method")
		}
		getMethod := r.Methods()[0]
		if getMethod.HTTPMethod != http.MethodPut {
			t.Errorf("got: %v want: %v", getMethod.HTTPMethod, http.MethodPut)
		}
		if !reflect.DeepEqual(getMethod.MethodOperation, moTest) {
			t.Errorf("got: %v want: %v", getMethod.MethodOperation, moTest)
		}
	})
	t.Run("options patch", func(t *testing.T) {
		r := rest.NewResource("pet")
		r.Patch(moTest, ct)
		if len(r.Methods()) != 1 {
			t.Errorf("expecting on method")
		}
		getMethod := r.Methods()[0]
		if getMethod.HTTPMethod != http.MethodPatch {
			t.Errorf("got: %v want: %v", getMethod.HTTPMethod, http.MethodPatch)
		}
		if !reflect.DeepEqual(getMethod.MethodOperation, moTest) {
			t.Errorf("got: %v want: %v", getMethod.MethodOperation, moTest)
		}
	})
	t.Run("delete method", func(t *testing.T) {
		r := rest.NewResource("pet")
		r.Delete(moTest, ct)
		if len(r.Methods()) != 1 {
			t.Errorf("expecting on method")
		}
		getMethod := r.Methods()[0]
		if getMethod.HTTPMethod != http.MethodDelete {
			t.Errorf("got: %v want: %v", getMethod.HTTPMethod, http.MethodDelete)
		}
		if !reflect.DeepEqual(getMethod.MethodOperation, moTest) {
			t.Errorf("got: %v want: %v", getMethod.MethodOperation, moTest)
		}
	})
	t.Run("options method", func(t *testing.T) {
		r := rest.NewResource("pet")
		r.Options(moTest, ct)
		if len(r.Methods()) != 1 {
			t.Errorf("expecting on method")
		}
		getMethod := r.Methods()[0]
		if getMethod.HTTPMethod != http.MethodOptions {
			t.Errorf("got: %v want: %v", getMethod.HTTPMethod, http.MethodOptions)
		}
		if !reflect.DeepEqual(getMethod.MethodOperation, moTest) {
			t.Errorf("got: %v want: %v", getMethod.MethodOperation, moTest)
		}
	})
	t.Run("connect method", func(t *testing.T) {
		r := rest.NewResource("pet")
		r.Connect(moTest, ct)
		if len(r.Methods()) != 1 {
			t.Errorf("expecting on method")
		}
		getMethod := r.Methods()[0]
		if getMethod.HTTPMethod != http.MethodConnect {
			t.Errorf("got: %v want: %v", getMethod.HTTPMethod, http.MethodConnect)
		}
		if !reflect.DeepEqual(getMethod.MethodOperation, moTest) {
			t.Errorf("got: %v want: %v", getMethod.MethodOperation, moTest)
		}
	})
	t.Run("head method", func(t *testing.T) {
		r := rest.NewResource("pet")
		r.Head(moTest, ct)
		if len(r.Methods()) != 1 {
			t.Errorf("expecting on method")
		}
		getMethod := r.Methods()[0]
		if getMethod.HTTPMethod != http.MethodHead {
			t.Errorf("got: %v want: %v", getMethod.HTTPMethod, http.MethodHead)
		}
		if !reflect.DeepEqual(getMethod.MethodOperation, moTest) {
			t.Errorf("got: %v want: %v", getMethod.MethodOperation, moTest)
		}
	})
	t.Run("trace method", func(t *testing.T) {
		r := rest.NewResource("pet")
		r.Trace(moTest, ct)
		if len(r.Methods()) != 1 {
			t.Errorf("expecting on method")
		}
		getMethod := r.Methods()[0]
		if getMethod.HTTPMethod != http.MethodTrace {
			t.Errorf("got: %v want: %v", getMethod.HTTPMethod, http.MethodTrace)
		}
		if !reflect.DeepEqual(getMethod.MethodOperation, moTest) {
			t.Errorf("got: %v want: %v", getMethod.MethodOperation, moTest)
		}
	})
}

func TestNewResource(t *testing.T) {
	t.Run("new resource", func(t *testing.T) {
		name := "pet"
		rest.NewResource(name)
		defer assertNoPanic(t)
	})

	testCase := []struct {
		name           string
		expectingPanic bool
	}{
		{"abcd122344_", false},
		{"abc_2344323wfwef123123weddvsdGGGGGGG", false},
		{"eerewrWERWERWERWrw", false},
	}

	// appending invalid character test cases for every char in rest.ResourceReservedChar
	for _, char := range rest.ResourceReservedChar {
		testCase = append(testCase, struct {
			name           string
			expectingPanic bool
		}{"aasd123" + string(char) + "abcd1234_", true})
	}

	t.Run("invalid char set panic", func(t *testing.T) {
		for _, tt := range testCase {
			t.Run(tt.name, func(t *testing.T) {
				if tt.expectingPanic {
					defer func() {
						err := recover()
						if err != nil {
							if _, ok := err.(*rest.ErrorResourceCharNotAllowed); !ok {
								t.Errorf("got: %T want: %T", err, &rest.ErrorResourceCharNotAllowed{})
							}
						}
					}()
				} else {
					defer assertNoPanic(t)
				}

				rest.NewResource(tt.name)

				if tt.expectingPanic {
					t.Fatalf("The code did not panic")
				}
			})
		}
	})
}

func TestResourceIntegration(t *testing.T) {
	r := rest.NewResource("car")
	r.Resource("find", func(r *rest.Resource) {
		r.Resource("left", func(r *rest.Resource) {
		})
		r.Resource("right", func(r *rest.Resource) {
		})
	})
	if r.Path() != "car" {
		t.Errorf("got : %v want: %v", r.Path(), "car")
	}
	findNode := r.Resources()[0]
	if findNode.Path() != "find" {
		t.Errorf("got : %v want: %v", findNode.Path(), "find")
	}
	if len(findNode.Resources()) != 2 {
		t.Fatalf("expecting 2 sub nodes got: %v", len(findNode.Resources()))
	}
	directionResources := findNode.Resources()
	sort.Slice(directionResources, func(i, j int) bool {
		return directionResources[i].Path() < directionResources[j].Path()
	})
	if directionResources[0].Path() != "left" {
		t.Errorf("got : %v want: %v", findNode.Resources()[0].Path(), "left")
	}
	if directionResources[1].Path() != "right" {
		t.Errorf("got : %v want: %v", findNode.Resources()[1].Path(), "right")
	}
}

type MiddlewareSpy struct {
	name   string
	writer io.Writer
	called bool
}

func (m *MiddlewareSpy) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		m.called = true
		if m.writer != nil {
			m.writer.Write([]byte(m.name + "->"))
		}
		next.ServeHTTP(rw, r)
	})
}

func TestUseMiddlewareOnMethods(t *testing.T) {
	t.Run("method get middleware stack from resource stack", func(t *testing.T) {
		middleware := &MiddlewareSpy{}
		r := rest.NewResource("my name")
		r.Use(middleware.Middleware)
		op := &OperationStub{}
		m := r.Get(rest.NewMethodOperation(op, rest.NewResponse(200)), mustGetJSONContentType())
		req, _ := http.NewRequest("GET", "", nil)
		resp := httptest.NewRecorder()
		m.ServeHTTP(resp, req)
		if !op.wasCall {
			t.Errorf("operation was not called")
		}
		if !middleware.called {
			t.Errorf("middleware was not called")
		}
	})
	t.Run("only methods declared after calling Use ", func(t *testing.T) {
		middleware := &MiddlewareSpy{}
		r := rest.NewResource("my name")
		opnm := &OperationStub{}

		mNoMiddleware := r.Get(rest.NewMethodOperation(opnm, rest.NewResponse(200)), mustGetJSONContentType())

		r.Use(middleware.Middleware)

		opm := &OperationStub{}

		mMiddleware := r.Get(rest.NewMethodOperation(opm, rest.NewResponse(200)), mustGetJSONContentType())
		req, _ := http.NewRequest("GET", "", nil)
		resp := httptest.NewRecorder()

		mMiddleware.ServeHTTP(resp, req)

		if !opm.wasCall {
			t.Errorf("operation was not called")
		}

		if !middleware.called {
			t.Errorf("middleware was not called")
		}

		middleware.called = false
		respnm := httptest.NewRecorder()

		mNoMiddleware.ServeHTTP(respnm, req)

		if !opnm.wasCall {
			t.Errorf("operation was not called")
		}

		if middleware.called {
			t.Errorf("not expecting calling the middleware")
		}
	})
	t.Run("middleware order", func(t *testing.T) {
		orderWriter := bytes.NewBufferString("")
		m1 := &MiddlewareSpy{name: "m1", writer: orderWriter}
		m2 := &MiddlewareSpy{name: "m2", writer: orderWriter}
		m3 := &MiddlewareSpy{name: "m3", writer: orderWriter}
		r := rest.NewResource("test")

		r.Use(m1.Middleware, m2.Middleware, m3.Middleware)

		op := &OperationStub{}
		method := r.Get(rest.NewMethodOperation(op, rest.NewResponse(200)), mustGetJSONContentType())
		req, _ := http.NewRequest("GET", "", nil)
		res := httptest.NewRecorder()

		method.ServeHTTP(res, req)

		assertResponseCode(t, res, 200)

		want := "m1->m2->m3->"
		if orderWriter.String() != want {
			t.Errorf("got: %v want: %v", orderWriter.String(), want)
		}
	})
}

func TestUseMiddlewareOnResource(t *testing.T) {
	middleware := &MiddlewareSpy{}
	r := rest.NewResource("my name")

	r.Use(middleware.Middleware)

	op := &OperationStub{}
	var m *rest.Method
	r.Resource("sub1", func(r *rest.Resource) {
		m = r.Get(rest.NewMethodOperation(op, rest.NewResponse(200)), mustGetJSONContentType())
	})

	req, _ := http.NewRequest("GET", "", nil)
	resp := httptest.NewRecorder()

	m.ServeHTTP(resp, req)

	if !op.wasCall {
		t.Errorf("operation was not called")
	}

	if !middleware.called {
		t.Errorf("middleware was not called")
	}
}

type AuthenticatorSpy struct {
	called bool
}

func (s *AuthenticatorSpy) Authenticate(i rest.Input) rest.AuthError {
	s.called = true
	return nil
}

func TestOverwriteCoreSecurityMiddleware(t *testing.T) {
	middleware := &MiddlewareSpy{}
	r := rest.NewResource("my name")

	r.OverwriteCoreSecurityMiddleware(middleware.Middleware)

	op := &OperationStub{}
	auth := &AuthenticatorSpy{}
	so := rest.SecurityOperation{auth, rest.NewResponse(401), rest.NewResponse(403)}
	apiKeyScheme := rest.NewAPIKeySecurityScheme("apikey", rest.NewHeaderParameter("X-Key", reflect.String), so)
	var m *rest.Method
	r.Resource("sub1", func(r *rest.Resource) {
		m = r.Get(rest.NewMethodOperation(op, rest.NewResponse(200)), mustGetJSONContentType()).
			WithSecurity(apiKeyScheme)
	})

	req, _ := http.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()

	m.ServeHTTP(resp, req)

	if !op.wasCall {
		t.Errorf("operation was not called")
	}

	if auth.called {
		t.Errorf("not expecting calling auth operation")
	}

	if !middleware.called {
		t.Errorf("middleware was not called")
	}

	t.Run("order", func(t *testing.T) {
		orderWriter := bytes.NewBufferString("")
		m1 := &MiddlewareSpy{name: "m1", writer: orderWriter}
		m2 := &MiddlewareSpy{name: "m2", writer: orderWriter}
		m3 := &MiddlewareSpy{name: "m3", writer: orderWriter}
		m4 := &MiddlewareSpy{name: "m4", writer: orderWriter}
		m5 := &MiddlewareSpy{name: "m5", writer: orderWriter}

		r := rest.NewResource("test")

		op := &OperationStub{}
		method := r.Get(rest.NewMethodOperation(op, rest.NewResponse(200)), mustGetJSONContentType())

		var sub1Method *rest.Method
		var sub1MethodSecMidd *rest.Method
		var sub2Method *rest.Method

		r.Resource("sub1", func(r *rest.Resource) {
			r.Use(m1.Middleware)
			sub1Method = r.Get(rest.NewMethodOperation(op, rest.NewResponse(200)), mustGetJSONContentType())
			r.OverwriteCoreSecurityMiddleware(m4.Middleware)
			r.Use(m2.Middleware, m3.Middleware)
			sub1MethodSecMidd = r.Post(rest.NewMethodOperation(op, rest.NewResponse(200)), mustGetJSONContentType())
			r.Resource("sub2", func(r *rest.Resource) {
				r.Use(m5.Middleware)
				sub2Method = r.Get(rest.NewMethodOperation(op, rest.NewResponse(200)), mustGetJSONContentType())
			})
		})

		req, _ := http.NewRequest("GET", "", nil)
		res := httptest.NewRecorder()
		orderWriter.Reset()
		method.ServeHTTP(res, req)

		assertResponseCode(t, res, 200)

		want := ""
		if orderWriter.String() != want {
			t.Errorf("got: %v want: %v", orderWriter.String(), want)
		}

		orderWriter.Reset()
		sub1Method.ServeHTTP(res, req)

		assertResponseCode(t, res, 200)

		want = "m1->"
		if orderWriter.String() != want {
			t.Errorf("got: %v want: %v", orderWriter.String(), want)
		}

		orderWriter.Reset()
		sub1MethodSecMidd.ServeHTTP(res, req)

		assertResponseCode(t, res, 200)

		want = "m1->m2->m3->m4->"
		if orderWriter.String() != want {
			t.Errorf("got: %v want: %v", orderWriter.String(), want)
		}

		orderWriter.Reset()
		sub2Method.ServeHTTP(res, req)

		assertResponseCode(t, res, 200)

		want = "m1->m2->m3->m5->m4->"
		if orderWriter.String() != want {
			t.Errorf("got: %v want: %v", orderWriter.String(), want)
		}
	})
}
