package resource_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/ehsoc/resource"
)

type EncodeDecoderSpy struct {
	encodeCalled bool
	decodeCalled bool
}

func (e *EncodeDecoderSpy) Encode(w io.Writer, v interface{}) error {
	e.encodeCalled = true
	return nil
}

func (e *EncodeDecoderSpy) Decode(r io.Reader, v interface{}) error {
	e.decodeCalled = true
	return nil
}

func TestEncoderDecoderSelector(t *testing.T) {
	t.Run("getting an existent key", func(t *testing.T) {
		e := &EncodeDecoderSpy{}
		wantContentType := "test/message"
		contentTypes := resource.NewHTTPContentTypeSelector(resource.Response{})
		contentTypes.Add(wantContentType, e, false)
		encoderDecoder, err := contentTypes.GetEncoderDecoder(wantContentType)
		if err != nil {
			t.Fatalf("Not expecting error: %v", err)
		}
		buf := bytes.NewBufferString("")
		e.Encode(buf, "")
		AssertTrue(t, e.encodeCalled)
		encoderDecoder.Decode(buf, "")
		AssertTrue(t, e.decodeCalled)
	})
	t.Run("getting a non existent key", func(t *testing.T) {
		e := &EncodeDecoderSpy{}
		wantContentType := "test/message"
		contentTypes := resource.NewHTTPContentTypeSelector(resource.Response{})
		contentTypes.Add(wantContentType, e, false)
		_, err := contentTypes.GetEncoderDecoder("randomkey")
		if err == nil {
			t.Errorf("Was expecting error.")
		}
	})

}

func TestGetDefaultEncoderDecoder(t *testing.T) {
	t.Run("no default encdec", func(t *testing.T) {
		e := &EncodeDecoderSpy{}
		wantContentType := "test/message"
		contentTypes := resource.NewHTTPContentTypeSelector(resource.Response{})
		contentTypes.Add(wantContentType, e, false)
		_, _, err := contentTypes.GetDefaultEncoderDecoder()
		if err == nil {
			t.Errorf("Was expecting error.")
		}
	})
	t.Run("get default encdec", func(t *testing.T) {
		e := &EncodeDecoderSpy{}
		wantContentType := "test/message"
		contentTypes := resource.NewHTTPContentTypeSelector(resource.Response{})
		contentTypes.Add("random/json", e, false)
		contentTypes.Add("r/xml", e, true)
		//The last overwrites all
		contentTypes.Add(wantContentType, e, true)
		contentTypes.Add("r/tson", e, false)
		contentTypes.Add("r/ttext", e, false)
		got, _, err := contentTypes.GetDefaultEncoderDecoder()
		if err != nil {
			t.Fatalf("Not expecting error: %v", err)
		}
		if got != wantContentType {
			t.Errorf("got:%s want:%s", got, wantContentType)
		}
	})

}
