package zettel

import (
	"strings"
	"testing"

	"github.com/friedenberg/zit/src/alfa/typ"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/etikett"
	"github.com/friedenberg/zit/src/echo/age_io"
)

func makeEtiketten(t *testing.T, vs ...string) (es etikett.Set) {
	es = etikett.MakeSet()

	for _, v := range vs {
		if err := es.AddString(v); err != nil {
			t.Fatalf("%s", err)
		}
	}

	return
}

func makeAkteExt(t *testing.T, v string) (es typ.Typ) {
	if err := es.Set(v); err != nil {
		t.Fatalf("%s", err)
	}

	return
}

type stringBuilderCloser struct {
	*strings.Builder
}

func (b stringBuilderCloser) Close() error {
	return nil
}

func (b stringBuilderCloser) Sha() sha.Sha {
	return sha.Sha{}
}

type akteWriterFactory struct {
	stringBuilderCloser
}

func (aw akteWriterFactory) AkteWriter() (age_io.Writer, error) {
	return aw, nil
}

func readFormat(t *testing.T, f Format, contents string) (z Zettel, a string) {
	t.Helper()

	awf := akteWriterFactory{
		stringBuilderCloser{Builder: &strings.Builder{}},
	}

	c := FormatContextRead{
		In:                strings.NewReader(contents),
		AkteWriterFactory: awf,
	}

	n, err := f.ReadFrom(&c)

	if err != nil {
		t.Fatalf("failed to read zettel format: %s", err)
	}

	if n != int64(len(contents)) {
		t.Fatalf("expected to read %d but only read %d", len(contents), n)
	}

	z = c.Zettel
	a = awf.String()

	return
}
