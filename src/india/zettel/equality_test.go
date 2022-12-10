package zettel

import (
	"reflect"
	"testing"

	"github.com/friedenberg/zit/src/bravo/test_logz"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/echo/bezeichnung"
)

func TestMakeEtiketten(t1 *testing.T) {
	t := test_logz.T{T: t1}

	vs := []string{
		"tag1",
		"tag2",
		"tag3",
	}

	var sut kennung.EtikettSet
	var err error

	if sut, err = kennung.MakeSetStrings(vs...); err != nil {
		t.Fatalf("%s", err)
	}

	if sut.Len() != 3 {
		t.Fatalf("expected len 3 but got %d", sut.Len())
	}

	{
		ac := len(sut.Elements())

		if ac != 3 {
			t.Fatalf("expected len 3 but got %d", ac)
		}
	}

	sut2 := sut.Copy()

	if sut2.Len() != 3 {
		t.Fatalf("expected len 3 but got %d", sut2.Len())
	}

	{
		ac := sut.SortedString()

		if !reflect.DeepEqual(ac, vs) {
			t.Fatalf("expected %q but got %q", vs, ac)
		}
	}

	{
		ex := "tag1, tag2, tag3"
		ac := sut.String()

		if ac != ex {
			t.Fatalf("expected %q but got %q", ex, ac)
		}
	}

	{
		ex := "tag1, tag2, tag3"
		ac := sut.Copy().String()

		if ac != ex {
			t.Fatalf("expected %q but got %q", ex, ac)
		}
	}
}

func TestEqualitySelf(t1 *testing.T) {
	t := test_logz.T{T: t1}

	text := &Objekte{
		Bezeichnung: bezeichnung.Make("the title"),
		Etiketten: makeEtiketten(t,
			"tag1",
			"tag2",
			"tag3",
		).Copy(),
		Typ: makeAkteExt(t, "text/plain"),
	}

	if !text.Equals(text) {
		t.Fatalf("expected %v to equal itself", text)
	}
}

func TestEqualityNotSelf(t1 *testing.T) {
	t := test_logz.T{T: t1}

	text := &Objekte{
		Bezeichnung: bezeichnung.Make("the title"),
		Etiketten: makeEtiketten(t,
			"tag1",
			"tag2",
			"tag3",
		).Copy(),
		Typ: makeAkteExt(t, "text/plain"),
	}

	text1 := &Objekte{
		Bezeichnung: bezeichnung.Make("the title"),
		Etiketten: makeEtiketten(t,
			"tag1",
			"tag2",
			"tag3",
		).Copy(),
		Typ: makeAkteExt(t, "text/plain"),
	}

	if !text.Equals(text1) {
		t.Fatalf("expected %v to equal %v", text, text1)
	}
}
