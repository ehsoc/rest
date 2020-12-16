package encdec_test

import (
	"bytes"
	"testing"

	"github.com/ehsoc/rest/encdec"
)

func TestTextEncoder(t *testing.T) {
	t.Run("text", func(t *testing.T) {
		buf := bytes.NewBufferString("")
		mytext := "mytext"
		encoder := encdec.TextEncoder{}
		err := encoder.Encode(buf, mytext)
		if err != nil {
			t.Fatalf("not expecting error: %v", err)
		}
		got := buf.String()
		if mytext != got {
			t.Errorf("got: %q want: %q", got, mytext)
		}
	})
	t.Run("text pointer", func(t *testing.T) {
		buf := bytes.NewBufferString("")
		mytext := "mytext"
		encoder := encdec.TextEncoder{}
		err := encoder.Encode(buf, &mytext)
		if err != nil {
			t.Fatalf("not expecting error: %v", err)
		}
		got := buf.String()
		if mytext != got {
			t.Errorf("got: %q want: %q", got, mytext)
		}
	})
	t.Run("nil", func(t *testing.T) {
		buf := bytes.NewBufferString("")
		encoder := encdec.TextEncoder{}
		err := encoder.Encode(buf, nil)
		if err != encdec.ErrorTextDecoderNoString {
			t.Fatalf("got: %v want: %v", err, encdec.ErrorTextDecoderNoString)
		}
	})
}
