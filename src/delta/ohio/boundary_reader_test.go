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

	if !errors.IsEOF(err) {
		t.Errorf("expected %q but got %q", io.EOF, err)
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
		t.Logf("%d", i)
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

		if i == 2 {
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

		if i == 2 {
			end = true
		}

		readBoundary(end)
	}
}

func TestBigMac(t1 *testing.T) {
	t := test_logz.T{T: t1}

	dataRaw := getRawData()

	dataComp := new(bytes.Buffer)

	comp := zstd.NewWriter(dataComp)

	n, err := io.Copy(comp, strings.NewReader(dataRaw))

	if n == 0 {
		t.Fatalf("read 0")
	}

	t.AssertNoError(err)

	comp.Flush()

	sut := MakeBoundaryReader(zstd.NewReader(dataComp), "---\n")

	var n1 int
	n1, err = sut.ReadBoundary()

	if n1 != 4 {
		t.Errorf("expected 4 but got %d", n1)
	}

	t.AssertNoError(err)

	i := 0

	for {
		t.Logf("i: %d", i)

		n, err = io.Copy(io.Discard, sut)

		if n == 0 {
			t.Fatalf("read 0")
		}

		t.AssertNoError(err)

		n, err := sut.ReadBoundary()

		if n != 4 {
			t.Errorf("expected 4 but got %d", n)
		}

		i++

		if errors.IsEOF(err) {
			break
		}

		t.AssertNoError(err)
	}

	expected := 55

	if i != expected {
		t.Errorf("expected %d entries but got %d", expected, i)
	}
}

func getRawData() string {
	return `---
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Etikett XXXX
Etikett XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ png
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Etikett XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Etikett XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ png
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Etikett XXXX
Etikett XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ jpg
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Bezeichnung XXXX
Etikett XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Etikett XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Etikett XXXX
Etikett XXXX
Etikett XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2039724665.0
Typ pdf
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Etikett XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2039724665.0
Typ txt
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2039913141.0
Typ gpg
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Etikett XXXX
Etikett XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2040580662.0
Typ pdf
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2041010005.0
Typ jpeg
---
Bezeichnung XXXX
Etikett XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2041024286.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2041119818.0
Typ ttf
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2041120219.0
Typ pdf
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2041120264.0
Typ pdf
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Etikett XXXX
Etikett XXXX
Etikett XXXX
Etikett XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2041761909.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2042394778.0
Typ m4r
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2043145169.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2045916290.0
Typ toml-bookmark
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Etikett XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2049138420.0
Typ ttf
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2049673708.0
Typ png
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Etikett XXXX
Etikett XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2054347577.86051
Typ pdf
---
Bezeichnung XXXX
Etikett XXXX
Etikett XXXX
Etikett XXXX
Etikett XXXX
Etikett XXXX
Etikett XXXX
Etikett XXXX
Etikett XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2070274012.190275
Typ task
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2070624013.945645
Typ toml-bookmark
---
`
}
