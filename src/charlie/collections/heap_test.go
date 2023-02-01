package collections

import (
	"reflect"
	"testing"

	"github.com/friedenberg/zit/src/bravo/int_value"
	"github.com/friedenberg/zit/src/bravo/test_logz"
)

func TestReset(t1 *testing.T) {
	t := test_logz.T{T: t1}

	els := []int_value.IntValue{
		int_value.Make(1),
		int_value.Make(0),
		int_value.Make(3),
		int_value.Make(4),
		int_value.Make(2),
	}

	sut := MakeHeapFromSlice(els)

	if sut.Len() != 5 {
		t.Fatalf("expected len 5 but got %d", sut.Len())
	}

	sut.Reset()

	if sut.Len() != 0 {
		t.Fatalf("expected len 0 but got %d", sut.Len())
	}
}

func TestSaveAndRestore(t1 *testing.T) {
	t := test_logz.T{T: t1}

	els := []int_value.IntValue{
		int_value.Make(1),
		int_value.Make(0),
		int_value.Make(3),
		int_value.Make(4),
		int_value.Make(2),
	}

	sut := MakeHeapFromSlice(els)

	checkAllElements := func() {
		defer sut.Restore()
		for i := 0; i < 5; i++ {
			el, ok := sut.PopAndSave()
			ex := int_value.Make(i)

			if !ok {
				t.Fatalf("expected pop and save to return an element but got nothing")
			}

			if !el.Equals(ex) {
				t.Fatalf("expected pop and save to return %s but got %s", ex, el)
			}
		}
	}

	checkAllElements()

	for i := 0; i < 5; i++ {
		el, ok := sut.Pop()
		ex := int_value.Make(i)

		if !ok {
			t.Fatalf("expected pop and save to return an element but got nothing. Idx: %d", i)
		}

		if !el.Equals(ex) {
			t.Fatalf("expected pop and save to return %s but got %s", ex, el)
		}
	}
}

func Test3Sorted(t1 *testing.T) {
	t := test_logz.T{T: t1}

	els := []int_value.IntValue{
		int_value.Make(1),
		int_value.Make(0),
		int_value.Make(3),
		int_value.Make(4),
		int_value.Make(2),
	}

	sut := MakeHeapFromSlice(els)
	sorted := sut.Sorted()

	expected := []int_value.IntValue{
		int_value.Make(0),
		int_value.Make(1),
		int_value.Make(2),
		int_value.Make(3),
		int_value.Make(4),
	}

	if !reflect.DeepEqual([]int_value.IntValue(sorted), expected) {
		t.Fatalf("expected %#v but got %#v", expected, sorted)
	}
}
