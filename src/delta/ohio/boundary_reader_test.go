package ohio

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/DataDog/zstd"
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/test_logz"
)

func TestBoundaryReaderEmpty(t1 *testing.T) {
	t := test_logz.T{T: t1}

	data := ``
	r := strings.NewReader(data)
	sut := MakeBoundaryReader(r, "\n")

	b := &strings.Builder{}
	n, err := io.Copy(b, sut)
	if err != nil {
		t.Errorf("expected no error but got %s", err)
	}

	if n != 0 {
		t.Errorf("expected no bytes read but got %d", n)
	}

	actual := b.String()

	if data != actual {
		t.Errorf("expected %q but got %q", data, actual)
	}
}

func TestBoundaryReaderContainsBoundary(t1 *testing.T) {
	t := test_logz.T{T: t1}

	data := "\n"
	r := strings.NewReader(data)
	sut := MakeBoundaryReader(r, "\n")

	b := &strings.Builder{}

	var n1 int
	var err error
	n1, err = sut.ReadBoundary()

	if err != nil {
		t.Errorf("expected no error but got %q", err)
	}

	if n1 != 1 {
		t.Errorf("expected no bytes read but got %d", n1)
	}

	var n int64
	n, err = io.Copy(b, sut)
	if err != nil {
		t.Errorf("expected no error but got %s", err)
	}

	if n != 0 {
		t.Errorf("expected no bytes read but got %d", n)
	}

	actual := b.String()

	if "" != actual {
		t.Errorf("expected %q but got %q", "", actual)
	}
}

func TestBoundaryReaderEmptyReadBoundary(t1 *testing.T) {
	t := test_logz.T{T: t1}

	data := ""
	r := strings.NewReader(data)
	sut := MakeBoundaryReader(r, "\n")

	var n1 int
	var err error
	n1, err = sut.ReadBoundary()

	if err == nil {
		t.Errorf("expected error but got none")
	}

	if n1 != 0 {
		t.Errorf("expected no bytes read but got %d", n1)
	}
}

func TestBoundaryReaderSandwich(t1 *testing.T) {
	t := test_logz.T{T: t1}

	data := `---
content
content
content
---
`
	r := strings.NewReader(data)
	sut := MakeBoundaryReader(r, "---\n")

	var n1 int
	var err error
	n1, err = sut.ReadBoundary()

	t.AssertNoError(err)

	if n1 != 4 {
		t.Errorf("expected 4 bytes read but got %d", n1)
	}

	b := strings.Builder{}

	var n int64
	n, err = io.Copy(&b, sut)

	t.AssertNoError(err)

	actual := b.String()
	expected := "content\ncontent\ncontent\n"

	if actual != expected {
		t.Errorf("expected %q but got %q", expected, actual)
	}

	if n != int64(len(expected)) {
		t.Errorf("expected %d bytes read but got %d", len(expected), n1)
	}
}

func TestBoundaryReaderSandwich2(t1 *testing.T) {
	t := test_logz.T{T: t1}

	data := `---
1 blob
---
2 blob
---
3 blob
---
`
	r := strings.NewReader(data)
	sut := MakeBoundaryReader(r, "---\n")

	readBoundary := func(last bool) {
		var n1 int
		var err error
		n1, err = sut.ReadBoundary()

		if last {
			if !errors.IsEOF(err) {
				t.AssertError(io.EOF)
			}
		} else {
			t.AssertNoError(err)
		}

		if n1 != 4 {
			t.Errorf("expected 4 bytes read but got %d", n1)
		}
	}

	readBoundary(false)

	for i := 0; i < 3; i++ {
		b := strings.Builder{}

		var n int64
		n, err := io.Copy(&b, sut)

		t.AssertNoError(err)

		actual := b.String()
		expected := fmt.Sprintf("%d blob\n", i+1)

		if actual != expected {
			t.Errorf("expected %q but got %q", expected, actual)
		}

		if n != int64(len(expected)) {
			t.Errorf("expected %d bytes read but got %d", len(expected), n)
		}

		end := false

		if i == 3 {
			end = true
		}

		readBoundary(end)
	}
}

func TestBoundaryReaderSandwich3(t1 *testing.T) {
	t := test_logz.T{T: t1}

	dataRaw := `---
1 blob
---
2 blob
---
3 blob
---
`

	dataComp := new(bytes.Buffer)

	comp := zstd.NewWriter(dataComp)

	io.Copy(comp, strings.NewReader(dataRaw))

	comp.Flush()

	sut := MakeBoundaryReader(zstd.NewReader(dataComp), "---\n")

	readBoundary := func(last bool) {
		var n1 int
		var err error
		n1, err = sut.ReadBoundary()

		if last {
			if !errors.IsEOF(err) {
				t.AssertError(io.EOF)
			}
		} else {
			t.AssertNoError(err)
		}

		if n1 != 4 {
			t.Errorf("expected 4 bytes read but got %d", n1)
		}
	}

	readBoundary(false)

	for i := 0; i < 3; i++ {
		b := strings.Builder{}

		var n int64
		n, err := io.Copy(&b, sut)

		t.AssertNoError(err)

		actual := b.String()
		expected := fmt.Sprintf("%d blob\n", i+1)

		if actual != expected {
			t.Errorf("expected %q but got %q", expected, actual)
		}

		if n != int64(len(expected)) {
			t.Errorf("expected %d bytes read but got %d", len(expected), n)
		}

		end := false

		if i == 3 {
			end = true
		}

		readBoundary(end)
	}
}
