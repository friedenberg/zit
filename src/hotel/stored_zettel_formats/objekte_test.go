package stored_zettel_formats

import (
	"strings"
	"testing"

	"github.com/friedenberg/zit/src/alfa/bezeichnung"
	"github.com/friedenberg/zit/src/bravo/akte_ext"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/etikett"
	"github.com/friedenberg/zit/src/foxtrot/zettel"
	"github.com/friedenberg/zit/src/golf/stored_zettel"
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

	var z stored_zettel.Stored

	var err error

	f := Objekte{}

	if _, err = f.ReadFrom(&z, strings.NewReader(content)); err != nil {
		t.Fatalf("%q", err)
	}

	expected := stored_zettel.Stored{
		Mutter: sha.Sha{Value: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		Kinder: sha.Sha{Value: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		Zettel: zettel.Zettel{
			Akte:        sha.Sha{Value: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
			AkteExt:     akte_ext.AkteExt{Value: "md"},
			Bezeichnung: bezeichnung.Bezeichnung("the title"),
			Etiketten: map[string]etikett.Etikett{
				"a-tag": etikett.Etikett{Value: "a-tag"},
			},
		},
	}

	if !z.Equals(expected) {
		t.Fatalf("\nexpected: %#v\n  actual: %#v", expected, z)
	}
}

func TestWrite(t *testing.T) {
	z := stored_zettel.Stored{
		Mutter: sha.Sha{Value: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		Kinder: sha.Sha{Value: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		Zettel: zettel.Zettel{
			Akte:        sha.Sha{Value: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
			AkteExt:     akte_ext.AkteExt{Value: "md"},
			Bezeichnung: bezeichnung.Bezeichnung("the title"),
			Etiketten: map[string]etikett.Etikett{
				"a-tag": etikett.Etikett{Value: "a-tag"},
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
