package dir_layout

import (
	"io"
	"path"
	"strings"
	"testing"

	"code.linenisgreat.com/zit/go/zit/src/delta/age"
)

func makeAge(t *testing.T) *age.Age {
	t.Helper()

	d := t.TempDir()
	age, err := age.Generate(path.Join(d, "AgeIdentity"))
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

	o := WriteOptions{
		Age:    age,
		Writer: out,
	}

	if w, err = NewWriter(o); err != nil {
		t.Fatalf("%s", err)
	}

	if _, err = io.Copy(w, in); err != nil {
		t.Fatalf("%s", err)
	}

	w.Close()

	in = strings.NewReader(out.String())
	out = &strings.Builder{}

	var r *reader

	ro := ReadOptions{
		Age:    age,
		Reader: in,
	}

	if r, err = NewReader(ro); err != nil {
		t.Fatalf("%s", err)
	}

	if _, err = io.Copy(out, r); err != nil {
		t.Fatalf("%s", err)
	}

	if text != out.String() {
		t.Fatalf("expected '%s', but got '%s'", text, out.String())
	}
}
