package resource_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/ehsoc/resource"
	"github.com/ehsoc/resource/encdec"
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

func TestAdd(t *testing.T) {
	t.Run("nil maps", func(t *testing.T) {
		defer assertNoPanic(t)
		ct := resource.HTTPContentTypeSelector{}
		ct.Add("", encdec.JSONEncoderDecoder{}, true)
	})
}

func TestEncoderDecoderSelector(t *testing.T) {
	t.Run("getting an existent key on encoder", func(t *testing.T) {
		e := &EncodeDecoderSpy{}
		wantContentType := "test/message"
		contentTypes := resource.NewHTTPContentTypeSelector()
		contentTypes.Add(wantContentType, e, false)
		encoder, err := contentTypes.GetEncoder(wantContentType)
		if err != nil {
			t.Fatalf("Not expecting error: %v", err)
		}
		buf := bytes.NewBufferString("")
		encoder.Encode(buf, "")
		assertTrue(t, e.encodeCalled)
	})
	t.Run("getting a non existent key on encoder", func(t *testing.T) {
		e := &EncodeDecoderSpy{}
		wantContentType := "test/message"
		contentTypes := resource.NewHTTPContentTypeSelector()
		contentTypes.Add(wantContentType, e, false)
		_, err := contentTypes.GetEncoder("randomkey")
		if err == nil {
			t.Errorf("Was expecting error.")
		}
	})
	t.Run("getting an existent key on decoder", func(t *testing.T) {
		e := &EncodeDecoderSpy{}
		wantContentType := "test/message"
		contentTypes := resource.NewHTTPContentTypeSelector()
		contentTypes.Add(wantContentType, e, false)
		encoderDecoder, err := contentTypes.GetDecoder(wantContentType)
		if err != nil {
			t.Fatalf("Not expecting error: %v", err)
		}
		buf := bytes.NewBufferString("")
		e.Encode(buf, "")
		assertTrue(t, e.encodeCalled)
		encoderDecoder.Decode(buf, "")
		assertTrue(t, e.decodeCalled)
	})
	t.Run("getting a non existent key on decoder", func(t *testing.T) {
		e := &EncodeDecoderSpy{}
		wantContentType := "test/message"
		contentTypes := resource.NewHTTPContentTypeSelector()
		contentTypes.Add(wantContentType, e, false)
		_, err := contentTypes.GetDecoder("randomkey")
		if err == nil {
			t.Errorf("Was expecting error.")
		}
	})

}

func TestGetDefaultEncoderDecoder(t *testing.T) {
	t.Run("no default encdec", func(t *testing.T) {
		e := &EncodeDecoderSpy{}
		wantContentType := "test/message"
		contentTypes := resource.NewHTTPContentTypeSelector()
		contentTypes.Add(wantContentType, e, false)
		_, _, err := contentTypes.GetDefaultEncoder()
		if err == nil {
			t.Errorf("Was expecting error.")
		}
	})
	t.Run("get default encdec", func(t *testing.T) {
		e := &EncodeDecoderSpy{}
		wantContentType := "test/message"
		contentTypes := resource.NewHTTPContentTypeSelector()
		contentTypes.Add("random/json", e, false)
		contentTypes.Add("r/xml", e, true)
		//The last overwrites all
		contentTypes.Add(wantContentType, e, true)
		contentTypes.Add("r/tson", e, false)
		contentTypes.Add("r/ttext", e, false)
		got, _, err := contentTypes.GetDefaultDecoder()
		if err != nil {
			t.Fatalf("Not expecting error: %v", err)
		}
		if got != wantContentType {
			t.Errorf("got:%s want:%s", got, wantContentType)
		}
	})

}
