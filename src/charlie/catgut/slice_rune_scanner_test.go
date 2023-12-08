package catgut

import (
	"testing"
	"unicode/utf8"

	"github.com/friedenberg/zit/src/bravo/test_logz"
)

func testSliceRuneScannerDataValid() []Slice {
	return []Slice{
		{
			data: [2][]byte{[]byte("string")},
		},
		{
			data: [2][]byte{[]byte("\u2318")},
		},
	}
}

func TestSliceRuneScannerValid(t1 *testing.T) {
	for _, expected := range testSliceRuneScannerDataValid() {
		t1.Run(
			expected.String(),
			func(t2 *testing.T) {
				t := test_logz.T{T: t2}

				sut, err := MakeSliceRuneScanner(expected)
				t.AssertNoError(err)

				for _, rEx := range []rune(expected.String()) {
					widthEx := utf8.RuneLen(rEx)
					r, width, ok := sut.Scan()

					if !ok {
						t.Errorf("expected successful scan")
					}

					if r != rEx {
						t.Errorf("expected %c but got %c", rEx, r)
					}

					if width != widthEx {
						t.Errorf("expected %d but got %d", widthEx, width)
					}
				}

				_, _, ok := sut.Scan()

				if ok {
					t.Errorf("expected unsuccessful scan")
				}
			},
		)
	}
}

func testSliceRuneScannerDataInvalid() []Slice {
	return []Slice{
		{
			data: [2][]byte{[]byte("\xbd\xb2\x3d\xbc\x20\xe2\x8c\x98")},
		},
	}
}

func TestSliceRuneScannerInvalid(t1 *testing.T) {
	for _, expected := range testSliceRuneScannerDataInvalid() {
		t1.Run(
			expected.String(),
			func(t2 *testing.T) {
				t := test_logz.T{T: t2}

				sut, err := MakeSliceRuneScanner(expected)
				t.AssertNoError(err)

				_, _, ok := sut.Scan()

				if ok {
					t.Errorf("expected unsuccessful scan")
				}
			},
		)
	}
}
