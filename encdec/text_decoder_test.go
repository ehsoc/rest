package encdec_test

import (
	"bytes"
	"testing"

	"github.com/ehsoc/rest/encdec"
)

func TestTextDecoder(t *testing.T) {
	t.Run("decode text", func(t *testing.T) {
		want := "my text"
		got := ""
		buf := bytes.NewBufferString(want)
		decoder := encdec.TextDecoder{}
		err := decoder.Decode(buf, &got)
		if err != nil {
			t.Errorf("not expecting error: %v", err)
		}
		if got != want {
			t.Errorf("got:%v want:%v", got, want)
		}
	})
	t.Run("v nil", func(t *testing.T) {
		want := "my text"
		buf := bytes.NewBufferString(want)
		decoder := encdec.TextDecoder{}
		err := decoder.Decode(buf, nil)
		if err != encdec.ErrorTextDecoderNoValidPointer {
			t.Errorf("got error : %v want: %v", err, encdec.ErrorTextDecoderNoValidPointer)
		}
	})
	t.Run("v not a string pointer", func(t *testing.T) {
		want := "my text"
		number := 101
		buf := bytes.NewBufferString(want)
		decoder := encdec.TextDecoder{}
		err := decoder.Decode(buf, &number)
		if err != encdec.ErrorTextDecoderNoString {
			t.Errorf("got error : %v want: %v", err, encdec.ErrorTextDecoderNoString)
		}
	})
}
