package zettel

import (
	"testing"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/test_logz"
	"code.linenisgreat.com/zit/go/zit/src/echo/descriptions"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/test_metadatei_io"
)

func makeTestTextFormat(
	af *test_metadatei_io.BlobIOFactory,
) metadatei.TextFormat {
	if af == nil {
		af = test_metadatei_io.FixtureFactoryReadWriteCloser(nil)
	}

	return metadatei.MakeTextFormat(
		af,
		nil,
	)
}

func TestReadWithoutBlob(t1 *testing.T) {
	t := test_logz.T{T: t1}
	af := test_metadatei_io.FixtureFactoryReadWriteCloser(nil)

	actual, blob := readFormat(
		t,
		makeTestTextFormat(af),
		af,
		`---
# the title
- tag1
- tag2
- tag3
! md
---
`,
	)

	expected := &metadatei.Metadatei{
		Bezeichnung: descriptions.Make("the title"),
		Typ:         makeBlobExt(t, "md"),
	}

	expected.SetEtiketten(makeEtiketten(t,
		"tag1",
		"tag2",
		"tag3",
	))

	if !actual.Equals(expected) {
		t.Fatalf("zettel:\nexpected: %#v\n  actual: %#v", expected, actual)
	}

	if blob != "" {
		t.Fatalf("blob:\nexpected empty but got %q", blob)
	}
}

func TestReadWithoutBlobWithMultilineBezeichnung(t1 *testing.T) {
	t := test_logz.T{T: t1}

	af := test_metadatei_io.FixtureFactoryReadWriteCloser(nil)

	actual, blob := readFormat(
		t,
		makeTestTextFormat(af),
		af,
		`---
# the title
# continues
- tag1
- tag2
- tag3
! md
---
`,
	)

	expected := &metadatei.Metadatei{
		Bezeichnung: descriptions.Make("the title\ncontinues"),
		Typ:         makeBlobExt(t, "md"),
	}

	expected.SetEtiketten(makeEtiketten(t,
		"tag1",
		"tag2",
		"tag3",
	))

	if !actual.Equals(expected) {
		t.Fatalf("zettel:\nexpected: %#v\n  actual: %#v", expected, actual)
	}

	if blob != "" {
		t.Fatalf("blob:\nexpected empty but got %q", blob)
	}
}

func TestReadWithBlob(t1 *testing.T) {
	t := test_logz.T{T: t1}

	af := test_metadatei_io.FixtureFactoryReadWriteCloser(
		map[string]string{
			"fa8242e99f48966ca514092b4233b446851f42b57ad5031bf133e1dd76787f3e": "the body\n",
			"036a8e44e472523c0306946f2712f372c234f8a24532e933f1509ae4db0da064": "the body",
		},
	)

	actual, blob := readFormat(
		t,
		makeTestTextFormat(af),
		af,
		`---
# the title
- tag1
- tag2
- tag3
! md
---

the body
`,
	)

	expected := &metadatei.Metadatei{
		Bezeichnung: descriptions.Make("the title"),
		Typ:         makeBlobExt(t, "md"),
	}

	errors.PanicIfError(expected.Akte.Set(
		"036a8e44e472523c0306946f2712f372c234f8a24532e933f1509ae4db0da064",
	))

	expected.SetEtiketten(makeEtiketten(t,
		"tag1",
		"tag2",
		"tag3",
	))

	if !actual.Equals(expected) {
		t.Fatalf("zettel:\nexpected: %#v\n  actual: %#v", expected, actual)
	}

	expectedBlob := "the body\n"

	if expectedBlob != blob {
		t.Fatalf("blob:\nexpected: %#v\n  actual: %#v", expectedBlob, blob)
	}
}
