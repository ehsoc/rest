package httputil

import (
	"errors"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

const (
	defaultMaxMemory = 32 << 20 // 32 MB
)

// From the std library /net/http/request.go
// multipartByReader is a sentinel value.
// Its presence in Request.MultipartForm indicates that parsing of the request
// body has been handed off to a MultipartReader instead of ParseMultipartForm.
var multipartByReader = &multipart.Form{
	Value: make(map[string][]string),
	File:  make(map[string][]*multipart.FileHeader),
}

var ErrMissingFile = errors.New("httputil: no such file")

// GetFormFile returns the first file content and header for the provided form key.
// GetFormfile uses FormFile function underneth so calls ParseMultipartForm and ParseForm if necessary
func GetFormFile(r *http.Request, key string) (fileContent []byte, fileHeader *multipart.FileHeader, err error) {
	f, fh, err := r.FormFile(key)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	if b, err := ioutil.ReadAll(f); err == nil {
		return b, fh, nil
	}

	return nil, nil, err
}

// GetFiles returns all the files for the provided form key.
func GetFiles(r *http.Request, key string) ([]*multipart.FileHeader, error) {
	if r.MultipartForm == multipartByReader {
		return nil, errors.New("httputil: multipart handled by MultipartReader")
	}

	if r.MultipartForm == nil {
		err := r.ParseMultipartForm(defaultMaxMemory)
		if err != nil {
			return nil, err
		}
	}

	if r.MultipartForm != nil && r.MultipartForm.File != nil {
		if files := r.MultipartForm.File[key]; len(files) > 0 {
			return files, nil
		}
	}

	return nil, ErrMissingFile
}

// GetFormValues returns all the form values for the provided form key.
func GetFormValues(r *http.Request, key string) ([]string, error) {
	if r.Form == nil {
		r.ParseMultipartForm(defaultMaxMemory)
	}

	if vs := r.Form[key]; len(vs) > 0 {
		return vs, nil
	}

	return nil, nil
}
