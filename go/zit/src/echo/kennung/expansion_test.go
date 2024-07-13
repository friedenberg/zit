package kennung

import (
	"testing"

	"code.linenisgreat.com/zit/go/zit/src/bravo/expansion"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/test_logz"
)

func stringSliceEquals(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func TestStringSliceUnequal(t *testing.T) {
	expected := []string{
		"this",
		"is",
		"a",
	}

	actual := []string{
		"this",
		"is",
		"a",
		"tag",
	}

	if stringSliceEquals(expected, actual) {
		t.Errorf("expected unequal slices")
	}
}

func TestStringSliceEquals(t *testing.T) {
	expected := []string{
		"this",
		"is",
		"a",
		"tag",
	}

	actual := []string{
		"this",
		"is",
		"a",
		"tag",
	}

	if !stringSliceEquals(expected, actual) {
		t.Errorf("expected equal slices")
	}
}

func TestExpansionAll(t1 *testing.T) {
	t := test_logz.T{T: t1}
	e := MustTag("this-is-a-tag")
	ex := ExpandOne(&e, expansion.ExpanderAll)
	expected := []string{
		"a",
		"a-tag",
		"is",
		"is-a-tag",
		"tag",
		"this",
		"this-is",
		"this-is-a",
		"this-is-a-tag",
	}

	actual := iter.SortedStrings[Tag](ex)

	if !stringSliceEquals(actual, expected) {
		t.Errorf(
			"expanded tags don't match:\nexpected: %q\n  actual: %q",
			expected,
			actual,
		)
	}
}

func TestExpansionRight(t *testing.T) {
	e := MustTag("this-is-a-tag")
	ex := ExpandOne(&e, expansion.ExpanderRight)
	expected := []string{
		"this",
		"this-is",
		"this-is-a",
		"this-is-a-tag",
	}

	actual := iter.SortedStrings[Tag](ex)

	if !stringSliceEquals(actual, expected) {
		t.Errorf(
			"expanded tags don't match:\nexpected: %q\n  actual: %q",
			expected,
			actual,
		)
	}
}
