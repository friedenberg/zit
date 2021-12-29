package stored_zettel_formats

import (
	"strings"
	"testing"
)

// func toBase64(t *testing.T, in string) (out string) {
// 	t.Helper()

// 	sb := &strings.Builder{}

// 	w := base64.NewEncoder(base64.StdEncoding, sb)
// 	defer w.Close()

// 	var err error

// 	if _, err = io.Copy(w, strings.NewReader(in)); err != nil {
// 		t.Fatalf("%q", err)
// 	}

// 	out = sb.String()

// 	return
// }

// func fromBase64(t *testing.T, in string) (out string) {
// 	t.Helper()

// 	sb := &strings.Builder{}

// 	r := base64.NewDecoder(base64.StdEncoding, strings.NewReader(in))

// 	var err error

// 	if _, err = io.Copy(sb, r); err != nil {
// 		t.Fatalf("%q", err)
// 	}

// 	out = sb.String()

// 	return
// }

func TestRead(t *testing.T) {
	content := `Mutter e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Kinder e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
AkteExt md
Bezeichnung the title
Etikett a-tag
`

	var z _StoredZettel

	var err error

	f := Objekte{}

	if _, err = f.ReadFrom(&z, strings.NewReader(content)); err != nil {
		t.Fatalf("%q", err)
	}

	expected := _StoredZettel{
		Mutter: _Sha{"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		Kinder: _Sha{"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		Zettel: _Zettel{
			Akte:        _Sha{"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
			AkteExt:     _AkteExt{Value: "md"},
			Bezeichnung: _Bezeichnung("the title"),
			Etiketten: map[string]_Etikett{
				"a-tag": _Etikett{Value: "a-tag"},
			},
		},
	}

	if !z.Equals(expected) {
		t.Fatalf("\nexpected: %#v\n  actual: %#v", expected, z)
	}
}

func TestWrite(t *testing.T) {
	z := _StoredZettel{
		Mutter: _Sha{"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		Kinder: _Sha{"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		Zettel: _Zettel{
			Akte:        _Sha{"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
			AkteExt:     _AkteExt{Value: "md"},
			Bezeichnung: _Bezeichnung("the title"),
			Etiketten: map[string]_Etikett{
				"a-tag": _Etikett{Value: "a-tag"},
			},
		},
	}

	sb := &strings.Builder{}
	var err error

	f := Objekte{}

	if _, err = f.WriteTo(z, sb); err != nil {
		t.Fatalf("%q", err)
	}

	expected := `Mutter e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Kinder e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
AkteExt md
Bezeichnung the title
Etikett a-tag
`

	if expected != sb.String() {
		t.Fatalf("\nexpected: %#v\n  actual: %#v", expected, sb.String())
	}
}
