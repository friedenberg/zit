package collections

import (
	"reflect"
	"sort"
	"testing"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/test_logz"
	"github.com/friedenberg/zit/src/bravo/values"
)

func makeStringValues(vs ...string) (out []values.String) {
	out = make([]values.String, len(vs))

	for i, v := range vs {
		out[i] = values.MakeString(v)
	}

	return
}

func assertSet(
	t test_logz.T,
	sut schnittstellen.SetLike[values.String],
	vals []values.String,
) {
	t.Helper()

	// Len() int
	if sut.Len() != len(vals) {
		t.Fatalf("expected len %d but got %d", len(vals), sut.Len())
	}

	// Key(T) string
	{
		v := "wow"
		k := sut.Key(values.MakeString(v))

		if k != v {
			t.Fatalf("expected key %s but got %s", v, k)
		}
	}

	// Get(string) (T, bool)
	{
		for _, v := range vals {
			v1, ok := sut.Get(v.String())

			if !ok {
				t.Fatalf("expected sut to contain %s", v)
			}

			if v1 != v {
				t.Fatalf("expected %s but got %s", v, v1)
			}
		}
	}

	// ContainsKey(string) bool
	{
		for _, v := range vals {
			ok := sut.ContainsKey(v.String())

			if !ok {
				t.Fatalf("expected sut to contain %s", v)
			}
		}
	}

	{
		ex := vals
		ac := sut.Elements()

		sort.Slice(ac, func(i, j int) bool { return ac[i].Less(ac[j]) })

		if !reflect.DeepEqual(ex, ac) {
			t.Fatalf("expected %s but got %s", ex, ac)
		}
	}

	// Contains(T) bool
	for _, v := range vals {
		if !sut.Contains(v) {
			t.Fatalf("expected %s to contain %s", sut, v)
		}
	}

	// Copy
	{
		sutCopy := sut.CloneSetLike()

		if !sut.EqualsSetLike(sutCopy) {
			t.Fatalf("expected copy to equal original")
		}
	}

	// MutableCopy
	{
		sutCopy := sut.CloneMutableSetLike()

		if !sut.EqualsSetLike(sutCopy) {
			t.Fatalf("expected mutable copy to equal original")
		}

		sutCopy.Reset()

		if sut.EqualsSetLike(sutCopy) {
			t.Fatalf("expected reset mutable copy to not equal original")
		}
	}

	// Each(WriterFunc[T]) error
	// EachKey(WriterFuncKey) error
}

func TestSet(t1 *testing.T) {
	t := test_logz.T{T: t1}

	{
		vals := makeStringValues(
			"1 one",
			"2 two",
			"3 three",
		)

		sut := MakeSet[values.String](
			func(v values.String) string {
				return v.String()
			},
			vals...,
		)

		assertSet(t, sut, vals)
	}

	{
		vals := makeStringValues(
			"1 one",
			"2 two",
			"3 three",
		)

		sut := MakeMutableSet[values.String](
			func(v values.String) string {
				return v.String()
			},
			vals...,
		)

		assertSet(t, sut, vals)
	}

	{
		vals := makeStringValues(
			"1 one",
			"2 two",
			"3 three",
		)

		sut := MakeSet[values.String](
			(values.String).String,
			vals...,
		)

		assertSet(t, sut, vals)
	}

	{
		vals := makeStringValues(
			"1 one",
			"2 two",
			"3 three",
		)

		sut := MakeMutableSet[values.String](
			(values.String).String,
			vals...,
		)

		assertSet(t, sut, vals)
	}
}
