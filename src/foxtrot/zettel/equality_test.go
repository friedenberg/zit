package zettel

import (
	"testing"

	"github.com/friedenberg/zit/src/alfa/typ"
	"github.com/friedenberg/zit/src/delta/etikett"
)

func makeEtiketten(t *testing.T, vs ...string) (es etikett.Set) {
	es = etikett.MakeSet()

	for _, v := range vs {
		if err := es.AddString(v); err != nil {
			t.Fatalf("%s", err)
		}
	}

	return
}

func makeAkteExt(t *testing.T, v string) (es typ.Typ) {
	if err := es.Set(v); err != nil {
		t.Fatalf("%s", err)
	}

	return
}

func TestEqualitySelf(t *testing.T) {
	text := Zettel{
		Bezeichnung: "the title",
		Etiketten: makeEtiketten(t,
			"tag1",
			"tag2",
			"tag3",
		),
		Typ: makeAkteExt(t, "text/plain"),
	}

	if !text.Equals(text) {
		t.Fatalf("expected %v to equal itself", text)
	}
}

func TestEqualityNotSelf(t *testing.T) {
	text := Zettel{
		Bezeichnung: "the title",
		Etiketten: makeEtiketten(t,
			"tag1",
			"tag2",
			"tag3",
		),
		Typ: makeAkteExt(t, "text/plain"),
	}

	text1 := Zettel{
		Bezeichnung: "the title",
		Etiketten: makeEtiketten(t,
			"tag1",
			"tag2",
			"tag3",
		),
		Typ: makeAkteExt(t, "text/plain"),
	}

	if !text.Equals(text1) {
		t.Fatalf("expected %v to equal %v", text, text1)
	}
}
