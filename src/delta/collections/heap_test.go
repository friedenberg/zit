package collections

import (
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
		for i := 0; i < 5; i++ {
			el, ok := sut.PopAndSave()
			ex := int_value.Make(i)

			if !ok {
				t.Fatalf("expected pop and save to return an element but got nothing")
			}

			if !el.Equals(&ex) {
				t.Fatalf("expected pop and save to return %s but got %s", ex, el)
			}
		}
	}

	checkAllElements()
	sut.Restore()
	checkAllElements()
}
