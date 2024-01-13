package kennung

import (
	"os"
	"testing"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/test_logz"
)

func TestMain(m *testing.M) {
	errors.SetTesting()
	code := m.Run()
	os.Exit(code)
}

func assertListHasElement(t test_logz.T, list []string, exIdx int, el rune) {
	idx, ok := BinarySearchForRuneInEtikettenSortedStringSlice(list, el)

	if !ok {
		t.Fatalf("expected found value for index '%d' and rune '%c'", exIdx, el)
	}

	if idx != exIdx {
		t.Fatalf("expected index '%d' for rune '%c' but got '%d'", exIdx, el, idx)
	}
}

func TestBinarySearch(t1 *testing.T) {
	list1 := []string{
		"123",
		"2234",
		"3333",
		"444",
		"555",
	}

	type tc struct {
		description string
		list        []string
		exIdx       int
		char        rune
	}

	testCases := []tc{
		{
			list:  list1,
			exIdx: 0,
			char:  '1',
		},
		{
			list:  list1,
			exIdx: 1,
			char:  '2',
		},
		{
			list:  list1,
			exIdx: 2,
			char:  '3',
		},
		{
			list:  list1,
			exIdx: 3,
			char:  '4',
		},
		{
			list:  list1,
			exIdx: 4,
			char:  '5',
		},
	}

	for _, v := range testCases {
		t1.Run(
			v.description,
			func(t2 *testing.T) {
				t := test_logz.T{T: t2}
				assertListHasElement(t, v.list, v.exIdx, v.char)
			},
		)
	}
}
