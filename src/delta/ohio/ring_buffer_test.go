package ohio

import (
	"bytes"
	"testing"

	"github.com/friedenberg/zit/src/bravo/test_logz"
)

func TestRingBufferEmpty(t1 *testing.T) {
	t := test_logz.T{T: t1}
	sut := MakeRingBuffer(10)

	{
		actual := sut.Len()

		if sut.Len() != 0 {
			t.Errorf("expected %d but got %d", 0, actual)
		}
	}

	{
		n, err := sut.Write([]byte("test"))

		if n != 4 {
			t.Errorf("expected %d but got %d", 4, n)
		}

		if err != nil {
			t.Errorf("expected no error but got %s", err)
		}

		{
			expected := 4
			actual := sut.Len()

			if expected != actual {
				t.Errorf("expected %d but got %d", expected, actual)
			}
		}
	}

	{
		n := sut.PeekMatch([]byte("tes"))

		if n != 3 {
			t.Errorf("expected %d but got %d", 3, n)
		}
	}

	// {
	// 	start, end := sut.Find([]byte("t"))

	// 	if start != 0 {
	// 		t.Errorf("expected %d but got %d", 0, start)
	// 	}

	// 	if end != 0 {
	// 		t.Errorf("expected %d but got %d", 0, end)
	// 	}
	// }

	// {
	// 	start, end := sut.Find([]byte("test"))

	// 	if start != 0 {
	// 		t.Errorf("expected %d but got %d", 0, start)
	// 	}

	// 	if end != 3 {
	// 		t.Errorf("expected %d but got %d", 3, end)
	// 	}
	// }

	{
		n := sut.PeekMatch([]byte("test"))

		if n != 4 {
			t.Errorf("expected %d but got %d", 4, n)
		}
	}

	{
		n := sut.PeekMatch([]byte("testy"))

		if n != 4 {
			t.Errorf("expected %d but got %d", 4, n)
		}
	}

	{
		b := make([]byte, 4)
		n, err := sut.Read(b)

		if n != 4 {
			t.Errorf("expected %d but got %d", 4, n)
		}

		t.AssertEOF(err)

		actual := string(b)

		if actual != "test" {
			t.Errorf("expected %q but got %q", "test", actual)
		}
	}

	// {
	// 	t.Logf("%#v", sut)
	// 	start, end := sut.Find([]byte("t"))

	// 	if start != -1 {
	// 		t.Errorf("expected start %d but got %d", -1, start)
	// 	}

	// 	if end != -1 {
	// 		t.Errorf("expected end %d but got %d", -1, end)
	// 	}
	// }
}

func TestRingBufferEmptyTooBig(t1 *testing.T) {
	t := test_logz.T{T: t1}
	sut := MakeRingBuffer(5)

	for i := 0; i < 11; i++ {
		t.Logf("i: %d", i)
		{
			n, err := sut.Write([]byte("test"))

			if n != 4 {
				t.Errorf("expected %d but got %d", 4, n)
			}

			t.AssertNoError(err)
		}

		{
			b := make([]byte, 4)
			n, err := sut.Read(b)

			if n != 4 {
				t.Errorf("expected %d but got %d", 4, n)
			}

			t.AssertEOF(err)

			{
				actual := string(b[:n])
				expected := "test"

				if actual != expected {
					t.Errorf("expected %q but got %q", expected, actual)
				}
			}
		}
	}

	// {
	// 	sut.Write([]byte("test"))
	// 	start, end := sut.Find([]byte("test"))

	// 	if start != 48 {
	// 		t.Errorf("expected %d but got %d", 48, start)
	// 	}

	// 	if end != -1 {
	// 		t.Errorf("expected %d but got %d", -1, end)
	// 	}
	// }
}

func TestRingBufferEmptyTooSmall(t1 *testing.T) {
	t := test_logz.T{T: t1}
	sut := MakeRingBuffer(3)

	{
		n, err := sut.Write([]byte("teal"))

		if n != 3 {
			t.Errorf("expected %d but got %d", 4, n)
		}

		t.AssertEOF(err)

		if sut.Len() != 3 {
			t.Errorf("expected len 3 but got %d", sut.Len())
		}
	}

	{
		n, err := sut.Write([]byte("teal"))

		if n != 0 {
			t.Errorf("expected %d but got %d", 4, n)
		}

		t.AssertEOF(err)

		if sut.Len() != 3 {
			t.Errorf("expected len 3 but got %d", sut.Len())
		}
	}

	{
		b := make([]byte, 4)
		n, err := sut.Read(b)

		{
			expected := 3
			if n != expected {
				t.Errorf("expected %d but got %d", expected, n)
			}

			if sut.Len() != 0 {
				t.Errorf("expected len 0 but got %d", sut.Len())
			}
		}

		t.AssertEOF(err)

		actual := string(b[:n])
		expected := "tea"

		if actual != expected {
			t.Errorf("expected %q but got %q", expected, actual)
		}
	}

	{
		b := make([]byte, 4)
		n, err := sut.Read(b)

		{
			expected := 0
			if n != expected {
				t.Errorf("expected %d but got %d", expected, n)
			}
		}

		t.AssertEOF(err)
	}
}

func TestRingBufferDefault(t1 *testing.T) {
	t := test_logz.T{T: t1}
	t2 := t.Skip(1)
	sut := MakeRingBuffer(0)

	one_5 := make([]byte, 2730)
	half := make([]byte, 2048)

	l := 0

	write := func() {
		n, err := sut.Write(one_5)

		if n != len(one_5) {
			t2.Errorf("expected %d but got %d", len(one_5), n)
		}

		l += n

		t2.AssertNoError(err)

		if sut.Len() != l {
			t2.Errorf("expected len %d but got %d", l, sut.Len())
		}
	}

	read := func() {
		n, err := sut.Read(half)

		if n != len(half) {
			t2.Errorf("expected %d but got %d", len(half), n)
		}

		l -= n

		t2.AssertNoError(err)

		if sut.Len() != l {
			t2.Errorf("expected len %d but got %d", l, sut.Len())
		}
	}

	write()
	read()
	write()
	read()
	write()
	read()
}

func TestRingBufferDefaultReadFrom(t1 *testing.T) {
	t := test_logz.T{T: t1}
	sut := MakeRingBuffer(0)

	one_5 := bytes.NewBuffer(make([]byte, 2730))
	half := make([]byte, 2048)

	l := 0
	t2 := t.Skip(1)

	write := func() {
		n, err := sut.FillWith(one_5)
		one_5 = bytes.NewBuffer(make([]byte, 2730))

		if n != one_5.Len() {
			t2.Errorf("expected %d but got %d", one_5.Len(), n)
		}

		l += n

		t2.AssertNoError(err)

		if sut.Len() != l {
			t2.Errorf("expected len %d but got %d", l, sut.Len())
		}
	}

	read := func() {
		n, err := sut.Read(half)

		if n != len(half) {
			t2.Errorf("expected %d but got %d", len(half), n)
		}

		l -= n

		t2.AssertNoError(err)

		if sut.Len() != l {
			t2.Errorf("expected len %d but got %d", l, sut.Len())
		}
	}

	write()
	read()
	write()
	read()
	write()
	read()
}