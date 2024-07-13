package zettel

import (
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/bravo/test_logz"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_ptr"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/test_object_metadata_io"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type inlineTypChecker struct {
	answer bool
}

func (t inlineTypChecker) IsInlineTyp(k ids.Type) bool {
	return t.answer
}

func makeEtiketten(t test_logz.T, vs ...string) (es ids.TagSet) {
	var err error

	if es, err = collections_ptr.MakeValueSetString[ids.Tag, *ids.Tag](nil, vs...); err != nil {
		t.Fatalf("%s", err)
	}

	return
}

func makeBlobExt(t test_logz.T, v string) (es ids.Type) {
	if err := es.Set(v); err != nil {
		t.Fatalf("%s", err)
	}

	return
}

func readFormat(
	t1 test_logz.T,
	f object_metadata.TextFormat,
	af *test_object_metadata_io.BlobIOFactory,
	contents string,
) (z *object_metadata.Metadata, a string) {
	var zt sku.Transacted

	t := t1.Skip(1)

	n, err := f.ParseMetadata(
		strings.NewReader(contents),
		&zt,
	)
	if err != nil {
		t.Fatalf("failed to read zettel format: %s", err)
	}

	if n != int64(len(contents)) {
		t.Fatalf("expected to read %d but only read %d", len(contents), n)
	}

	z = zt.GetMetadata()
	a = af.CurrentBufferString()

	return
}
