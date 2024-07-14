package zettel

import (
	"reflect"
	"testing"

	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/test_logz"
	"code.linenisgreat.com/zit/go/zit/src/echo/descriptions"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
)

func TestMakeEtiketten(t1 *testing.T) {
	t := test_logz.T{T: t1}

	vs := []string{
		"tag1",
		"tag2",
		"tag3",
	}

	var sut ids.TagSet
	var err error

	if sut, err = ids.MakeTagSetStrings(vs...); err != nil {
		t.Fatalf("%s", err)
	}

	if sut.Len() != 3 {
		t.Fatalf("expected len 3 but got %d", sut.Len())
	}

	{
		ac := sut.Len()

		if ac != 3 {
			t.Fatalf("expected len 3 but got %d", ac)
		}
	}

	sut2 := sut.CloneSetLike()

	if sut2.Len() != 3 {
		t.Fatalf("expected len 3 but got %d", sut2.Len())
	}

	{
		ac := iter.SortedStrings[ids.Tag](sut)

		if !reflect.DeepEqual(ac, vs) {
			t.Fatalf("expected %q but got %q", vs, ac)
		}
	}

	{
		ex := "tag1, tag2, tag3"
		ac := iter.StringCommaSeparated[ids.Tag](sut)

		if ac != ex {
			t.Fatalf("expected %q but got %q", ex, ac)
		}
	}

	{
		ex := "tag1, tag2, tag3"
		ac := iter.StringCommaSeparated[ids.Tag](
			sut.CloneSetLike(),
		)

		if ac != ex {
			t.Fatalf("expected %q but got %q", ex, ac)
		}
	}
}

func TestEqualitySelf(t1 *testing.T) {
	t := test_logz.T{T: t1}

	text := &object_metadata.Metadata{
		Description: descriptions.Make("the title"),
		Type:        makeBlobExt(t, "text"),
	}

	text.SetTags(makeEtiketten(t,
		"tag1",
		"tag2",
		"tag3",
	))

	if !text.Equals(text) {
		t.Fatalf("expected %v to equal itself", text)
	}
}

func TestEqualityNotSelf(t1 *testing.T) {
	t := test_logz.T{T: t1}

	text := object_metadata.Metadata{
		Description: descriptions.Make("the title"),
		Type:        makeBlobExt(t, "text"),
	}

	text.SetTags(makeEtiketten(t,
		"tag1",
		"tag2",
		"tag3",
	))

	text1 := &object_metadata.Metadata{
		Description: descriptions.Make("the title"),
		Type:        makeBlobExt(t, "text"),
	}

	text1.SetTags(makeEtiketten(t,
		"tag1",
		"tag2",
		"tag3",
	))

	if !text.Equals(text1) {
		t.Fatalf("expected %v to equal %v", text, text1)
	}
}
