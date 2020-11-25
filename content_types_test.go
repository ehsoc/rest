package rest_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/ehsoc/rest"
	"github.com/ehsoc/rest/encdec"
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
		ct := rest.ContentTypes{}
		ct.Add("", encdec.JSONEncoderDecoder{}, true)
	})
}

func TestEncoderDecoderSelector(t *testing.T) {
	t.Run("getting an existent key on encoder", func(t *testing.T) {
		e := &EncodeDecoderSpy{}
		wantMIMEType := "test/message"
		ct := rest.NewContentTypes()
		ct.Add(wantMIMEType, e, false)
		encoder, err := ct.GetEncoder(wantMIMEType)
		if err != nil {
			t.Fatalf("Not expecting error: %v", err)
		}
		buf := bytes.NewBufferString("")
		encoder.Encode(buf, "")
		assertTrue(t, e.encodeCalled)
	})
	t.Run("getting a non existent key on encoder", func(t *testing.T) {
		e := &EncodeDecoderSpy{}
		wantMIMEType := "test/message"
		ct := rest.NewContentTypes()
		ct.Add(wantMIMEType, e, false)
		_, err := ct.GetEncoder("randomkey")
		if err == nil {
			t.Errorf("Was expecting error.")
		}
	})
	t.Run("getting an existent key on decoder", func(t *testing.T) {
		e := &EncodeDecoderSpy{}
		wantMIMEType := "test/message"
		ct := rest.NewContentTypes()
		ct.Add(wantMIMEType, e, false)
		encoderDecoder, err := ct.GetDecoder(wantMIMEType)
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
		wantMIMEType := "test/message"
		ct := rest.NewContentTypes()
		ct.Add(wantMIMEType, e, false)
		_, err := ct.GetDecoder("randomkey")
		if err == nil {
			t.Errorf("Was expecting error.")
		}
	})
}

func TestGetDefaultEncoderDecoder(t *testing.T) {
	t.Run("no default encdec", func(t *testing.T) {
		e := &EncodeDecoderSpy{}
		wantMIMEType := "test/message"
		ct := rest.NewContentTypes()
		ct.Add(wantMIMEType, e, false)
		_, _, err := ct.GetDefaultEncoder()
		if err == nil {
			t.Errorf("Was expecting error.")
		}
	})
	t.Run("get default encdec", func(t *testing.T) {
		e := &EncodeDecoderSpy{}
		wantMIMEType := "test/message"
		ct := rest.NewContentTypes()
		ct.Add("random/json", e, false)
		ct.Add("r/xml", e, true)
		// The last overwrites all
		ct.Add(wantMIMEType, e, true)
		ct.Add("r/tson", e, false)
		ct.Add("r/ttext", e, false)
		got, _, err := ct.GetDefaultDecoder()
		if err != nil {
			t.Fatalf("Not expecting error: %v", err)
		}
		if got != wantMIMEType {
			t.Errorf("got:%s want:%s", got, wantMIMEType)
		}
	})
}
