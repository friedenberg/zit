package zettel

import (
	"reflect"
	"testing"

	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/bravo/test_logz"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/echo/bezeichnung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
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

	sut2 := sut.CloneSetLike()

	if sut2.Len() != 3 {
		t.Fatalf("expected len 3 but got %d", sut2.Len())
	}

	{
		ac := iter.SortedStrings[kennung.Etikett](sut)

		if !reflect.DeepEqual(ac, vs) {
			t.Fatalf("expected %q but got %q", vs, ac)
		}
	}

	{
		ex := "tag1, tag2, tag3"
		ac := iter.StringCommaSeparated[kennung.Etikett](sut)

		if ac != ex {
			t.Fatalf("expected %q but got %q", ex, ac)
		}
	}

	{
		ex := "tag1, tag2, tag3"
		ac := iter.StringCommaSeparated[kennung.Etikett](
			sut.CloneSetLike(),
		)

		if ac != ex {
			t.Fatalf("expected %q but got %q", ex, ac)
		}
	}
}

func TestEqualitySelf(t1 *testing.T) {
	t := test_logz.T{T: t1}

	text := metadatei.Metadatei{
		Bezeichnung: bezeichnung.Make("the title"),
		Etiketten: makeEtiketten(t,
			"tag1",
			"tag2",
			"tag3",
		).CloneSetPtrLike(),
		Typ: makeAkteExt(t, "text"),
	}

	if !text.Equals(text) {
		t.Fatalf("expected %v to equal itself", text)
	}
}

func TestEqualityNotSelf(t1 *testing.T) {
	t := test_logz.T{T: t1}

	text := metadatei.Metadatei{
		Bezeichnung: bezeichnung.Make("the title"),
		Etiketten: makeEtiketten(t,
			"tag1",
			"tag2",
			"tag3",
		).CloneSetPtrLike(),
		Typ: makeAkteExt(t, "text"),
	}

	text1 := metadatei.Metadatei{
		Bezeichnung: bezeichnung.Make("the title"),
		Etiketten: makeEtiketten(t,
			"tag1",
			"tag2",
			"tag3",
		).CloneSetPtrLike(),
		Typ: makeAkteExt(t, "text"),
	}

	if !text.Equals(text1) {
		t.Fatalf("expected %v to equal %v", text, text1)
	}
}
