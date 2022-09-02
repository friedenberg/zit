package zettel

import (
	"testing"
)

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
