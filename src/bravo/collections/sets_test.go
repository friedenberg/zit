package collections

import (
	"os"
	"reflect"
	"sort"
	"testing"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/test_logz"
)

func TestMain(m *testing.M) {
	errors.SetTesting()
	code := m.Run()
	os.Exit(code)
}

type stringValue string

func (s stringValue) String() string {
	return string(s)
}

func makeStringValues(vs ...string) (out []stringValue) {
	out = make([]stringValue, len(vs))

	for i, v := range vs {
		out[i] = stringValue(v)
	}

	return
}

func (s *stringValue) Set(v string) error {
	*s = stringValue(v)
	return nil
}

func assertSet(t test_logz.T, sut Set[stringValue], vals []stringValue) {
	t.Helper()

	//Len() int
	if sut.Len() != len(vals) {
		t.Fatalf("expected len %d but got %d", len(vals), sut.Len())
	}

	//Key(T) string
	{
		v := "wow"
		k := sut.Key(stringValue(v))

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

		sort.Slice(ac, func(i, j int) bool { return ac[i] < ac[j] })

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
		sutCopy := sut.Copy()

		if !sut.Equals(sutCopy) {
			t.Fatalf("expected copy to equal original")
		}
	}

	// MutableCopy
	{
		sutCopy := sut.MutableCopy()

		if !sut.Equals(sutCopy) {
			t.Fatalf("expected mutable copy to equal original")
		}

		sutCopy.Reset(nil)

		if sut.Equals(sutCopy) {
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

		sut := MakeSetGeneric[stringValue](
			func(v stringValue) string {
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

		sut := MakeMutableSetGeneric[stringValue](
			func(v stringValue) string {
				return v.String()
			},
			vals...,
		)

		assertSet(t, Set[stringValue]{SetLike: sut}, vals)
	}

	{
		vals := makeStringValues(
			"1 one",
			"2 two",
			"3 three",
		)

		sut := MakeValueSet[stringValue](
			vals...,
		)

		assertSet(t, Set[stringValue]{SetLike: sut}, vals)
	}

	{
		vals := makeStringValues(
			"1 one",
			"2 two",
			"3 three",
		)

		sut := MakeMutableValueSet[stringValue](
			vals...,
		)

		assertSet(t, Set[stringValue]{SetLike: sut}, vals)
	}
}
