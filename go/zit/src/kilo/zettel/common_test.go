package zettel

import (
	"strings"

	"code.linenisgreat.com/zit-go/src/bravo/test_logz"
	"code.linenisgreat.com/zit-go/src/charlie/collections_ptr"
	"code.linenisgreat.com/zit-go/src/echo/kennung"
	"code.linenisgreat.com/zit-go/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit-go/src/foxtrot/test_metadatei_io"
	"code.linenisgreat.com/zit-go/src/hotel/sku"
)

type inlineTypChecker struct {
	answer bool
}

func (t inlineTypChecker) IsInlineTyp(k kennung.Typ) bool {
	return t.answer
}

func makeEtiketten(t test_logz.T, vs ...string) (es kennung.EtikettSet) {
	var err error

	if es, err = collections_ptr.MakeValueSetString[kennung.Etikett, *kennung.Etikett](nil, vs...); err != nil {
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
) (z *metadatei.Metadatei, a string) {
	var zt sku.Transacted

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

	z = zt.GetMetadatei()
	a = af.CurrentBufferString()

	return
}
