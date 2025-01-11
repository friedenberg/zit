package object_metadata

import (
	"bytes"
	"strings"
	"testing"

	"code.linenisgreat.com/zit/go/zit/src/bravo/test_logz"
	"code.linenisgreat.com/zit/go/zit/src/triple_hyphen_io"
)

func Test1(t1 *testing.T) {
	t := test_logz.T{
		T: t1,
	}

	in := `---
metadatei
---

body
`
	mExpected := "metadatei\n"
	bExpected := "body\n"
	nExpected := int64(len(in))

	mr := &bytes.Buffer{}
	ar := &bytes.Buffer{}

	r := triple_hyphen_io.Reader{
		Metadata: mr,
		Blob:     ar,
	}

	var n int64
	var err error

	n, err = r.ReadFrom(strings.NewReader(in))

	if n != nExpected {
		t.Errorf("expected to read %d but read %d", nExpected, n)
	}

	if err != nil {
		t.Errorf("expected no error but got %s", err)
	}

	mActual := string(mr.Bytes())

	if mActual != mExpected {
		t.Errorf("expected %q but got %q", mExpected, mActual)
	}

	bActual := string(ar.Bytes())

	if bActual != bExpected {
		t.Errorf("expected %q but got %q", bExpected, bActual)
	}
}

func Test2(t1 *testing.T) {
	t := test_logz.T{
		T: t1,
	}

	in := `---
metadatei
---
`
	mExpected := "metadatei\n"
	bExpected := ""
	nExpected := int64(len(in))

	mr := &bytes.Buffer{}
	ar := &bytes.Buffer{}

	r := triple_hyphen_io.Reader{
		Metadata: mr,
		Blob:     ar,
	}

	var n int64
	var err error

	n, err = r.ReadFrom(strings.NewReader(in))

	if n != nExpected {
		t.Errorf("expected to read %d but read %d", nExpected, n)
	}

	if err != nil {
		t.Errorf("expected no error but got %s", err)
	}

	mActual := string(mr.Bytes())

	if mActual != mExpected {
		t.Errorf("expected %q but got %q", mExpected, mActual)
	}

	bActual := string(ar.Bytes())

	if bActual != bExpected {
		t.Errorf("expected %q but got %q", bExpected, bActual)
	}
}
