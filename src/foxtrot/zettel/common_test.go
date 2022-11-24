package zettel

import (
	"bytes"
	"strings"

	"github.com/friedenberg/zit/src/bravo/test_logz"
	"github.com/friedenberg/zit/src/charlie/kennung"
	"github.com/friedenberg/zit/src/delta/test_metadatei_io"
)

func makeEtiketten(t test_logz.T, vs ...string) (es kennung.EtikettSet) {
	var err error

	if es, err = kennung.MakeSetStrings(vs...); err != nil {
		t.Fatalf("%s", err)
	}

	return
}

func makeAkteExt(t test_logz.T, v string) (es kennung.Typ) {
	if err := es.Set(v); err != nil {
		t.Fatalf("%s", err)
	}

	return
}

func readFormat(t test_logz.T, f Format, contents string) (z Zettel, a string) {
	t.Helper()

	awf := test_metadatei_io.NopFactoryReadWriter(bytes.NewBuffer(nil))

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
