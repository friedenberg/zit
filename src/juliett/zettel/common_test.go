package zettel

import (
	"strings"

	"github.com/friedenberg/zit/src/bravo/test_logz"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/echo/test_metadatei_io"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
)

type inlineTypChecker struct {
	answer bool
}

func (t inlineTypChecker) IsInlineTyp(k kennung.Typ) bool {
	return t.answer
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
	t1 test_logz.T,
	f metadatei.TextFormat,
	af *test_metadatei_io.AkteIOFactory,
	contents string,
) (z Objekte, a string) {
	var zt Transacted

	t := t1.Skip(1)

	n, err := f.ParseMetadatei(
		strings.NewReader(contents),
		&zt,
	)
	if err != nil {
		t.Fatalf("failed to read zettel format: %s", err)
	}

	if n != int64(len(contents)) {
		t.Fatalf("expected to read %d but only read %d", len(contents), n)
	}

	z = zt.Objekte
	a = af.CurrentBufferString()

	return
}
