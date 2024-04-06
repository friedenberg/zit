package heap

import (
	"reflect"
	"testing"

	"code.linenisgreat.com/zit/src/bravo/test_logz"
	"code.linenisgreat.com/zit/src/bravo/values"
)

func TestReset(t1 *testing.T) {
	t := test_logz.T{T: t1}

	els := []*values.Int{
		values.MakeInt(1),
		values.MakeInt(0),
		values.MakeInt(3),
		values.MakeInt(4),
		values.MakeInt(2),
	}

	sut := MakeHeapFromSliceUnsorted[values.Int, *values.Int](
		values.IntEqualer{},
		values.IntLessor{},
		values.IntResetter{},
		els,
	)

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

	els := []*values.Int{
		values.MakeInt(1),
		values.MakeInt(0),
		values.MakeInt(3),
		values.MakeInt(4),
		values.MakeInt(2),
	}

	eql := values.IntEqualer{}

	sut := MakeHeapFromSliceUnsorted[values.Int, *values.Int](
		eql,
		values.IntLessor{},
		values.IntResetter{},
		els,
	)

	checkAllElements := func() {
		defer sut.Restore()
		for i := 0; i < 5; i++ {
			el, ok := sut.PopAndSave()
			ex := values.MakeInt(i)

			if !ok {
				t.Fatalf(
					"expected pop and save to return an element but got nothing",
				)
			}

			if !eql.EqualsPtr(el, ex) {
				t.Fatalf(
					"expected pop and save to return %s but got %s",
					ex,
					el,
				)
			}
		}
	}

	checkAllElements()

	for i := 0; i < 5; i++ {
		el, ok := sut.Pop()
		ex := values.MakeInt(i)

		if !ok {
			t.Fatalf(
				"expected pop and save to return an element but got nothing. Idx: %d",
				i,
			)
		}

		if !eql.EqualsPtr(el, ex) {
			t.Fatalf("expected pop and save to return %s but got %s", ex, el)
		}
	}
}

func Test3Sorted(t1 *testing.T) {
	t := test_logz.T{T: t1}

	els := []*values.Int{
		values.MakeInt(1),
		values.MakeInt(0),
		values.MakeInt(3),
		values.MakeInt(4),
		values.MakeInt(2),
	}

	eql := values.IntEqualer{}

	sut := MakeHeapFromSliceUnsorted[values.Int, *values.Int](
		eql,
		values.IntLessor{},
		values.IntResetter{},
		els,
	)

	sorted := sut.Sorted()

	expected := []*values.Int{
		values.MakeInt(0),
		values.MakeInt(1),
		values.MakeInt(2),
		values.MakeInt(3),
		values.MakeInt(4),
	}

	if !reflect.DeepEqual([]*values.Int(sorted), expected) {
		t.Fatalf("expected %#v but got %#v", expected, sorted)
	}
}

func TestDupes(t1 *testing.T) {
	t := test_logz.T{T: t1}
	t.SkipTest()

	els := []*values.Int{
		values.MakeInt(1),
		values.MakeInt(0),
		values.MakeInt(3),
		values.MakeInt(3),
		values.MakeInt(4),
		values.MakeInt(2),
	}

	eql := values.IntEqualer{}

	sut := MakeHeapFromSliceUnsorted[values.Int, *values.Int](
		eql,
		values.IntLessor{},
		values.IntResetter{},
		els,
	)

	sorted := sut.Sorted()

	expected := []*values.Int{
		values.MakeInt(0),
		values.MakeInt(1),
		values.MakeInt(2),
		values.MakeInt(3),
		values.MakeInt(4),
	}

	if !reflect.DeepEqual([]*values.Int(sorted), expected) {
		t.Fatalf("expected %s but got %s", expected, sorted)
	}
}

func TestMerge(t1 *testing.T) {
	t := test_logz.T{T: t1}

	eql := values.IntEqualer{}
	llr := values.IntLessor{}

	els := []*values.Int{
		values.MakeInt(1),
		values.MakeInt(0),
		values.MakeInt(3),
		values.MakeInt(4),
		values.MakeInt(2),
	}

	otherStream := MakeHeapFromSliceUnsorted[values.Int, *values.Int](
		eql,
		llr,
		values.IntResetter{},
		[]*values.Int{
			values.MakeInt(8),
			values.MakeInt(9),
			values.MakeInt(3),
			values.MakeInt(7),
			values.MakeInt(6),
		},
	)

	expected := []*values.Int{
		values.MakeInt(0),
		values.MakeInt(1),
		values.MakeInt(2),
		values.MakeInt(3),
		values.MakeInt(4),
		values.MakeInt(6),
		values.MakeInt(7),
		values.MakeInt(8),
		values.MakeInt(9),
	}

	sut := MakeHeapFromSliceUnsorted[values.Int, *values.Int](
		eql,
		llr,
		values.IntResetter{},
		els,
	)

	actual := make([]*values.Int, 0)

	err := MergeStream(
		sut,
		otherStream.PopError,
		func(v *values.Int) (err error) {
			actual = append(actual, v)
			return
		},
	)
	if err != nil {
		t.AssertNoError(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected %q but got %q", expected, actual)
	}
}
