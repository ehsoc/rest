package resource_test

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"reflect"
	"testing"

	"github.com/ehsoc/resource"
)

func TestGetQuery(t *testing.T) {
	t.Run("get query", func(t *testing.T) {
		r, _ := http.NewRequest("POST", "/?foo=1&foo=2", nil)
		p := resource.NewQueryArrayParameter("foo", nil)
		parameters := resource.Parameters{}
		parameters.AddParameter(p)
		input := resource.Input{r, parameters, nil}
		querySlice, err := input.GetQuery("foo")
		assertNoErrorFatal(t, err)
		if len(querySlice) != 2 {
			t.Fatalf("got :%v want: %v", len(querySlice), 2)
		}
	})
	t.Run("parameter not found", func(t *testing.T) {
		r, _ := http.NewRequest("POST", "/?foo=1&foo=2", nil)
		p := resource.NewQueryArrayParameter("bar", nil)
		parameters := resource.Parameters{}
		parameters.AddParameter(p)
		input := resource.Input{r, parameters, nil}
		_, err := input.GetQuery("foo")
		if _, ok := err.(*resource.TypeErrorParameterNotDefined); !ok {
			t.Errorf("got: %T want: %T", err, &resource.TypeErrorParameterNotDefined{})
		}
	})
}

func GetURIParamStub(r *http.Request, key string) string {
	return "myIdValue"
}

func TestGetURIParam(t *testing.T) {
	t.Run("get uri", func(t *testing.T) {
		// I set up the get function on InputContextKey("uriparamfunc") key in the context.
		ctx := context.WithValue(context.Background(), resource.InputContextKey("uriparamfunc"), GetURIParamStub)
		r, _ := http.NewRequestWithContext(ctx, "POST", "/", nil)
		p := resource.NewURIParameter("myId", reflect.String)
		parameters := resource.Parameters{}
		parameters.AddParameter(p)
		input := resource.Input{r, parameters, nil}
		uriValue, err := input.GetURIParam("myId")
		assertNoErrorFatal(t, err)
		wantURIValue := "myIdValue"
		if uriValue != wantURIValue {
			t.Errorf("got: %v want: %v", uriValue, wantURIValue)
		}
	})
	t.Run("get uri parameter function not defined", func(t *testing.T) {
		r, _ := http.NewRequest("POST", "/", nil)
		p := resource.NewURIParameter("myId", reflect.String)
		parameters := resource.Parameters{}
		parameters.AddParameter(p)
		input := resource.Input{r, parameters, nil}
		_, err := input.GetURIParam("myId")
		if _, ok := err.(*resource.TypeErrorGetURIParamFunctionNotDefined); !ok {
			t.Errorf("got: %T want: %T", err, &resource.TypeErrorGetURIParamFunctionNotDefined{})
		}
	})
	t.Run("get uri parameter not defined", func(t *testing.T) {
		// I set up the get function on InputContextKey("uriparamfunc") key in the context.
		ctx := context.WithValue(context.Background(), resource.InputContextKey("uriparamfunc"), GetURIParamStub)
		r, _ := http.NewRequestWithContext(ctx, "POST", "/", nil)
		p := resource.NewURIParameter("myId", reflect.String)
		parameters := resource.Parameters{}
		parameters.AddParameter(p)
		input := resource.Input{r, parameters, nil}
		_, err := input.GetURIParam("foo")
		if _, ok := err.(*resource.TypeErrorParameterNotDefined); !ok {
			t.Errorf("got: %T want: %T", err, &resource.TypeErrorParameterNotDefined{})
		}
	})
}

func TestGetBody(t *testing.T) {
	t.Run("get body", func(t *testing.T) {
		car := Car{Brand: "Ford"}
		gotCar := Car{}
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(car)
		r, _ := http.NewRequest("POST", "/", buf)
		input := resource.Input{r, resource.Parameters{}, Car{}}
		body, err := input.GetBody()
		assertNoErrorFatal(t, err)
		json.NewDecoder(body).Decode(&gotCar)
		if !reflect.DeepEqual(gotCar, car) {
			t.Errorf("got: %v want: %v", gotCar, car)
		}
	})
	t.Run("request body not defined", func(t *testing.T) {
		car := Car{Brand: "Ford"}
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(car)
		r, _ := http.NewRequest("POST", "/", buf)
		input := resource.Input{r, resource.Parameters{}, nil}
		_, err := input.GetBody()
		assertEqualError(t, err, resource.ErrorRequestBodyNotDefined)
	})
}

func TestGetFormFile(t *testing.T) {
	t.Run("get form file", func(t *testing.T) {
		buf := new(bytes.Buffer)
		w := multipart.NewWriter(buf)
		fileW, _ := w.CreateFormFile("file", "MyFileName.jpg")
		fileData := "filerandomstrings!"
		fileW.Write([]byte(fileData))
		r, _ := http.NewRequest("POST", "/", buf)

		p := resource.NewFileParameter("file")
		parameters := resource.Parameters{}
		parameters.AddParameter(p)
		input := resource.Input{r, parameters, nil}
		file, _, err := input.GetFormFile("file")
		assertNoErrorFatal(t, err)
		if string(file) != fileData {
			t.Errorf("got: %v want: %v", string(file), fileData)
		}

	})
	t.Run("parameter not found", func(t *testing.T) {
		r, _ := http.NewRequest("POST", "/", nil)
		p := resource.NewFileParameter("file")
		parameters := resource.Parameters{}
		parameters.AddParameter(p)
		input := resource.Input{r, parameters, nil}
		_, _, err := input.GetFormFile("foo")
		if _, ok := err.(*resource.TypeErrorParameterNotDefined); !ok {
			t.Errorf("got: %T want: %T", err, &resource.TypeErrorParameterNotDefined{})
		}
	})
}

//TODO: test other getters
