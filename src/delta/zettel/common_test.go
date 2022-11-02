package zettel

import (
	"strings"

	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/bravo/test_logz"
	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/charlie/typ"
)

func makeEtiketten(t test_logz.T, vs ...string) (es etikett.Set) {
	var err error

	if es, err = etikett.MakeSetStrings(vs...); err != nil {
		t.Fatalf("%s", err)
	}

	return
}

func makeAkteExt(t test_logz.T, v string) (es typ.Typ) {
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

func (aw akteWriterFactory) AkteWriter() (sha.WriteCloser, error) {
	return aw, nil
}

func readFormat(t test_logz.T, f Format, contents string) (z Zettel, a string) {
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
