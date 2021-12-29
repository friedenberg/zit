package objekte

import (
	"io"
	"strings"
	"testing"

	"github.com/friedenberg/zit/charlie/age"
)

func makeAge(t *testing.T) age.Age {
	t.Helper()

	d := t.TempDir()
	age, err := age.Generate(d)

	if err != nil {
		t.Fatalf("%s", err)
	}

	return age
}

func Test1(t *testing.T) {
	var err error
	age := makeAge(t)

	text := `test string`
	in := strings.NewReader(text)
	out := &strings.Builder{}

	var w *writer

	if w, err = NewWriter(age, out); err != nil {
		t.Fatalf("%s", err)
	}

	if _, err = io.Copy(w, in); err != nil {
		t.Fatalf("%s", err)
	}

	w.Close()

	in = strings.NewReader(out.String())
	out = &strings.Builder{}

	var r *reader

	if r, err = NewReader(age, in); err != nil {
		t.Fatalf("%s", err)
	}

	if _, err = io.Copy(out, r); err != nil {
		t.Fatalf("%s", err)
	}

	if text != out.String() {
		t.Fatalf("expected '%s', but got '%s'", text, out.String())
	}
}

func TestReadWrite(t *testing.T) {
	age := makeAge(t)

	text := `the text`

	eText, _, err := EncodeBase64(age, text)

	if err != nil {
		t.Fatalf("%s", err)
	}

	dText, _, err := DecodeBase64(age, eText)

	if err != nil {
		t.Fatalf("%s", err)
	}

	if text != dText {
		t.Fatalf("expected '%s' but got '%s'", text, dText)
	}
}
