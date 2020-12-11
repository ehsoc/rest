package rest_test

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"reflect"
	"testing"

	"github.com/ehsoc/rest"
	"github.com/ehsoc/rest/encdec"
)

func TestGetQuery(t *testing.T) {
	t.Run("get query", func(t *testing.T) {
		r, _ := http.NewRequest("POST", "/?foo=1&foo=2", nil)
		p := rest.NewQueryArrayParameter("foo", nil)
		parameters := rest.ParameterCollection{}
		parameters.AddParameter(p)
		input := rest.Input{r, parameters, rest.RequestBody{}, nil}
		querySlice, err := input.GetQuery("foo")
		assertNoErrorFatal(t, err)
		if len(querySlice) != 2 {
			t.Fatalf("got :%v want: %v", len(querySlice), 2)
		}
	})
	t.Run("parameter not found", func(t *testing.T) {
		r, _ := http.NewRequest("POST", "/?foo=1&foo=2", nil)
		p := rest.NewQueryArrayParameter("bar", nil)
		parameters := rest.ParameterCollection{}
		parameters.AddParameter(p)
		input := rest.Input{r, parameters, rest.RequestBody{}, nil}
		_, err := input.GetQuery("foo")
		if _, ok := err.(*rest.ErrorParameterNotDefined); !ok {
			t.Errorf("got: %T want: %T", err, &rest.ErrorParameterNotDefined{})
		}
	})
}

func TestGetHeader(t *testing.T) {
	t.Run("get header", func(t *testing.T) {
		r, _ := http.NewRequest("POST", "/", nil)
		p := rest.NewHeaderParameter("foo", reflect.String)
		parameters := rest.ParameterCollection{}
		parameters.AddParameter(p)
		want := "myHeaderValue"
		r.Header.Set(p.Name, want)
		input := rest.Input{r, parameters, rest.RequestBody{}, nil}
		got, err := input.GetHeader(p.Name)
		assertNoErrorFatal(t, err)
		assertStringEqual(t, got, want)
	})
	t.Run("parameter not found", func(t *testing.T) {
		r, _ := http.NewRequest("POST", "/", nil)
		parameters := rest.ParameterCollection{}
		input := rest.Input{r, parameters, rest.RequestBody{}, nil}
		_, err := input.GetHeader("foo")
		if _, ok := err.(*rest.ErrorParameterNotDefined); !ok {
			t.Errorf("got: %T want: %T", err, &rest.ErrorParameterNotDefined{})
		}
	})
}

func GetURIParamStub(r *http.Request, key string) string {
	return "myIdValue"
}

func TestGetURIParam(t *testing.T) {
	t.Run("get uri", func(t *testing.T) {
		// I set up the get function on InputContextKey("uriparamfunc") key in the context.
		ctx := context.WithValue(context.Background(), rest.InputContextKey("uriparamfunc"), GetURIParamStub)
		r, _ := http.NewRequest("POST", "/", nil)
		r = r.WithContext(ctx)
		p := rest.NewURIParameter("myId", reflect.String)
		parameters := rest.ParameterCollection{}
		parameters.AddParameter(p)
		input := rest.Input{r, parameters, rest.RequestBody{}, nil}
		uriValue, err := input.GetURIParam("myId")
		assertNoErrorFatal(t, err)
		wantURIValue := "myIdValue"
		if uriValue != wantURIValue {
			t.Errorf("got: %v want: %v", uriValue, wantURIValue)
		}
	})
	t.Run("get uri parameter function not defined", func(t *testing.T) {
		r, _ := http.NewRequest("POST", "/", nil)
		p := rest.NewURIParameter("myId", reflect.String)
		parameters := rest.ParameterCollection{}
		parameters.AddParameter(p)
		input := rest.Input{r, parameters, rest.RequestBody{}, nil}
		_, err := input.GetURIParam("myId")
		if _, ok := err.(*rest.ErrorGetURIParamFunctionNotDefined); !ok {
			t.Errorf("got: %T want: %T", err, &rest.ErrorGetURIParamFunctionNotDefined{})
		}
	})
	t.Run("get uri parameter not defined", func(t *testing.T) {
		// I set up the get function on InputContextKey("uriparamfunc") key in the context.
		ctx := context.WithValue(context.Background(), rest.InputContextKey("uriparamfunc"), GetURIParamStub)
		r, _ := http.NewRequest("POST", "/", nil)
		r = r.WithContext(ctx)
		p := rest.NewURIParameter("myId", reflect.String)
		parameters := rest.ParameterCollection{}
		parameters.AddParameter(p)
		input := rest.Input{r, parameters, rest.RequestBody{}, nil}
		_, err := input.GetURIParam("foo")
		if _, ok := err.(*rest.ErrorParameterNotDefined); !ok {
			t.Errorf("got: %T want: %T", err, &rest.ErrorParameterNotDefined{})
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
		input := rest.Input{r, rest.ParameterCollection{}, rest.RequestBody{"", Car{}, true}, nil}
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
		input := rest.Input{r, rest.ParameterCollection{}, rest.RequestBody{}, nil}
		_, err := input.GetBody()
		assertEqualError(t, err, rest.ErrorRequestBodyNotDefined)
	})
}

func TestGetFormFiles(t *testing.T) {
	t.Run("get form file", func(t *testing.T) {
		buf := new(bytes.Buffer)
		w := multipart.NewWriter(buf)
		fileW, _ := w.CreateFormFile("file", "MyFileName.jpg")
		fileData := "filerandomstrings!"
		fileW.Write([]byte(fileData))
		file2W, _ := w.CreateFormFile("file", "MyFileName2.jpg")
		fileData2 := "filerandomstrings2!"
		file2W.Write([]byte(fileData2))
		w.Close()
		r, _ := http.NewRequest("POST", "/", buf)
		r.Header.Set("Content-Type", w.FormDataContentType())
		p := rest.NewFileParameter("file")
		parameters := rest.ParameterCollection{}
		parameters.AddParameter(p)
		input := rest.Input{r, parameters, rest.RequestBody{}, nil}
		files, err := input.GetFormFiles("file")
		assertNoErrorFatal(t, err)
		if len(files) != 2 {
			t.Fatalf("expecting 2 files")
		}
		if files[0].Filename != "MyFileName.jpg" {
			t.Errorf("got: %v want: %v", files[0].Filename, "MyFileName.jpg")
		}
		if files[1].Filename != "MyFileName2.jpg" {
			t.Errorf("got: %v want: %v", files[1].Filename, "MyFileName2.jpg")
		}
	})
	t.Run("parameter not found", func(t *testing.T) {
		r, _ := http.NewRequest("POST", "/", nil)
		p := rest.NewFileParameter("file")
		parameters := rest.ParameterCollection{}
		parameters.AddParameter(p)
		input := rest.Input{r, parameters, rest.RequestBody{}, nil}
		_, _, err := input.GetFormFile("foo")
		if _, ok := err.(*rest.ErrorParameterNotDefined); !ok {
			t.Errorf("got: %T want: %T", err, &rest.ErrorParameterNotDefined{})
		}
	})
}

func TestGetFormFile(t *testing.T) {
	t.Run("get form file", func(t *testing.T) {
		buf := new(bytes.Buffer)
		w := multipart.NewWriter(buf)
		fileW, _ := w.CreateFormFile("file", "MyFileName.jpg")
		fileData := "filerandomstrings!"
		fileW.Write([]byte(fileData))
		w.Close()
		r, _ := http.NewRequest("POST", "/", buf)
		r.Header.Set("Content-Type", w.FormDataContentType())
		p := rest.NewFileParameter("file")
		parameters := rest.ParameterCollection{}
		parameters.AddParameter(p)
		input := rest.Input{r, parameters, rest.RequestBody{}, nil}
		file, _, err := input.GetFormFile("file")
		assertNoErrorFatal(t, err)
		if string(file) != fileData {
			t.Errorf("got: %v want: %v", string(file), fileData)
		}
	})
	t.Run("parameter not found", func(t *testing.T) {
		r, _ := http.NewRequest("POST", "/", nil)
		p := rest.NewFileParameter("file")
		parameters := rest.ParameterCollection{}
		parameters.AddParameter(p)
		input := rest.Input{r, parameters, rest.RequestBody{}, nil}
		_, _, err := input.GetFormFile("foo")
		if _, ok := err.(*rest.ErrorParameterNotDefined); !ok {
			t.Errorf("got: %T want: %T", err, &rest.ErrorParameterNotDefined{})
		}
	})
}

func TestGetFormValues(t *testing.T) {
	t.Run("get form file", func(t *testing.T) {
		b := new(bytes.Buffer)
		w := multipart.NewWriter(b)
		key := "additionalMetadata"
		additionalMetaData := "My Additional Metadata"
		fieldW, _ := w.CreateFormField(key)
		fieldW.Write([]byte(additionalMetaData))
		additionalMetaData2 := "My Additional Metadata2"
		field2W, _ := w.CreateFormField(key)
		field2W.Write([]byte(additionalMetaData2))
		w.Close()
		r, _ := http.NewRequest(http.MethodPost, "/", b)
		r.Header.Set("Content-Type", w.FormDataContentType())

		p := rest.NewFormDataParameter(key, reflect.String, encdec.TextDecoder{})
		parameters := rest.ParameterCollection{}
		parameters.AddParameter(p)
		input := rest.Input{r, parameters, rest.RequestBody{}, nil}
		values, err := input.GetFormValues(key)
		if err != nil {
			t.Fatalf("was not expecting an error")
		}
		if len(values) != 2 {
			t.Errorf("expecting 2 elements")
		}
		if values[0] != additionalMetaData {
			t.Errorf("got: %v want: %v", values[0], additionalMetaData)
		}
		if values[1] != additionalMetaData2 {
			t.Errorf("got: %v want: %v", values[1], additionalMetaData2)
		}
	})
	t.Run("parameter not found", func(t *testing.T) {
		r, _ := http.NewRequest("POST", "/", nil)
		p := rest.NewFileParameter("file")
		parameters := rest.ParameterCollection{}
		parameters.AddParameter(p)
		input := rest.Input{r, parameters, rest.RequestBody{}, nil}
		_, _, err := input.GetFormFile("foo")
		if _, ok := err.(*rest.ErrorParameterNotDefined); !ok {
			t.Errorf("got: %T want: %T", err, &rest.ErrorParameterNotDefined{})
		}
	})
}
