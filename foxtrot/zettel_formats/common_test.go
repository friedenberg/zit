package zettel_formats

import (
	"strings"
	"testing"

	"github.com/friedenberg/zit/charlie/etikett"
)

func makeEtiketten(t *testing.T, vs ...string) (es _EtikettSet) {
	es = etikett.MakeSet()

	for _, v := range vs {
		if err := es.AddString(v); err != nil {
			t.Fatalf("%s", err)
		}
	}

	return
}

func makeAkteExt(t *testing.T, v string) (es _AkteExt) {
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

func (b stringBuilderCloser) Sha() _Sha {
	return _Sha{}
}

type akteWriterFactory struct {
	stringBuilderCloser
}

func (aw akteWriterFactory) AkteWriter() (_ObjekteWriter, error) {
	return aw, nil
}

func readFormat(t *testing.T, f _ZettelFormat, contents string) (z _Zettel, a string) {
	t.Helper()

	awf := akteWriterFactory{
		stringBuilderCloser{Builder: &strings.Builder{}},
	}

	c := _ZettelFormatContextRead{
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
