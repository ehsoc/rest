package httputil_test

import (
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"testing"

	"github.com/ehsoc/resource/encdec"
	"github.com/ehsoc/resource/httputil"
)

type Body struct {
	Name string
	Year int
}

func TestGetFile(t *testing.T) {
	t.Run("getting file", func(t *testing.T) {
		request := newMultiformRequest()
		fc, fh, err := httputil.GetFormFile(request, "file")
		if err != nil {
			t.Fatalf("Not expecting error: %v ", err)
		}
		assertMultipartFileEqual(t, file1Name, file1Content, fc, fh)
	})
	t.Run("file not found", func(t *testing.T) {
		request := newMultiformRequest()
		_, _, err := httputil.GetFormFile(request, "filenotfound")
		if err == nil {
			t.Fatalf("expecting error")
		}
	})

}

var testGetFiles = []struct {
	fileName    string
	fileContent string
}{
	{file1Name, file1Content},
	{file2Name, file2Content},
}

func TestGetFiles(t *testing.T) {
	t.Run("get files", func(t *testing.T) {
		request := newMultiformRequest()
		files, err := httputil.GetFiles(request, "file")
		if err != nil {
			t.Fatalf("Not expecting error: %v ", err)
		}

		for i := 0; i < len(files); i++ {
			tfile := testGetFiles[i]
			fh := files[i]
			t.Run(tfile.fileName, func(t *testing.T) {
				f, _ := fh.Open()
				defer f.Close()
				fc, _ := ioutil.ReadAll(f)
				assertMultipartFileEqual(t, tfile.fileName, tfile.fileContent, fc, fh)
			})
		}
	})
	t.Run("get files not found", func(t *testing.T) {
		request, _ := http.NewRequest("POST", "/", nil)
		_, err := httputil.GetFiles(request, "filenotfound")
		if err == nil {
			t.Fatalf("expecting error")
		}
	})

}

func assertMultipartFileEqual(t *testing.T, fileName, fileContent string, fc []byte, fh *multipart.FileHeader) {
	t.Helper()
	gotFileContent := string(fc)
	if gotFileContent != fileContent {
		t.Errorf("got: %v want: %v", gotFileContent, fileContent)
	}
	if fh.Filename != fileName {
		t.Errorf("got: %v want: %v", fh.Filename, fileName)
	}
}

var file1Name = "MyFileName1.jpg"
var file2Name = "MyFileName2.jpg"
var file1Content = "filerandomstrings!1"
var file2Content = "filerandomstrings!2"

func newMultiformRequest() *http.Request {
	b := new(bytes.Buffer)
	w := multipart.NewWriter(b)
	fileW, _ := w.CreateFormFile("file", file1Name)
	fileW.Write([]byte(file1Content))
	fileW, _ = w.CreateFormFile("file", file2Name)
	fileW.Write([]byte(file2Content))
	additionalMetaData := "My Additional Metadata"
	fieldW, _ := w.CreateFormField("additionalMetadata")
	fieldW.Write([]byte(additionalMetaData))
	mediaHeader := textproto.MIMEHeader{}
	mediaHeader.Set("Content-Type", "application/json; charset=UTF-8")
	mediaHeader.Set("Content-Disposition", "form-data; name=\"jsonPetData\"")
	jsonPetDataW, _ := w.CreatePart(mediaHeader)
	encoder := encdec.JSONEncoder{}
	wantCar := Body{}
	encoder.Encode(jsonPetDataW, wantCar)
	w.Close()
	request, _ := http.NewRequest(http.MethodPost, "/", b)
	request.Header.Set("Content-Type", w.FormDataContentType())
	return request
}
