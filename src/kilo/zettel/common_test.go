package zettel

import (
	"strings"

	"github.com/friedenberg/zit/src/bravo/test_logz"
	"github.com/friedenberg/zit/src/foxtrot/kennung"
	"github.com/friedenberg/zit/src/foxtrot/test_metadatei_io"
)

type alwaysInlineTypChecker struct{}

func (_ alwaysInlineTypChecker) IsInlineTyp(k kennung.Typ) bool {
	return true
}

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

func readFormat(
	t test_logz.T,
	f Format,
	af *test_metadatei_io.AkteIOFactory,
	contents string,
) (z Objekte, a string) {
	t.Helper()

	c := FormatContextRead{
		In: strings.NewReader(contents),
	}

	n, err := f.ReadFrom(&c)

	if err != nil {
		t.Fatalf("failed to read zettel format: %s", err)
	}

	if n != int64(len(contents)) {
		t.Fatalf("expected to read %d but only read %d", len(contents), n)
	}

	z = c.Zettel
	a = af.CurrentBufferString()

	return
}
