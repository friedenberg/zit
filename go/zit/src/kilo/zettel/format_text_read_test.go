package zettel

import (
	"testing"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/test_logz"
	"code.linenisgreat.com/zit/go/zit/src/echo/bezeichnung"
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

func TestReadWithoutAkte(t1 *testing.T) {
	t := test_logz.T{T: t1}
	af := test_metadatei_io.FixtureFactoryReadWriteCloser(nil)

	actual, akte := readFormat(
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
		Bezeichnung: bezeichnung.Make("the title"),
		Typ:         makeAkteExt(t, "md"),
	}

	expected.SetEtiketten(makeEtiketten(t,
		"tag1",
		"tag2",
		"tag3",
	))

	if !actual.Equals(expected) {
		t.Fatalf("zettel:\nexpected: %#v\n  actual: %#v", expected, actual)
	}

	if akte != "" {
		t.Fatalf("akte:\nexpected empty but got %q", akte)
	}
}

func TestReadWithoutAkteWithMultilineBezeichnung(t1 *testing.T) {
	t := test_logz.T{T: t1}

	af := test_metadatei_io.FixtureFactoryReadWriteCloser(nil)

	actual, akte := readFormat(
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
		Bezeichnung: bezeichnung.Make("the title\ncontinues"),
		Typ:         makeAkteExt(t, "md"),
	}

	expected.SetEtiketten(makeEtiketten(t,
		"tag1",
		"tag2",
		"tag3",
	))

	if !actual.Equals(expected) {
		t.Fatalf("zettel:\nexpected: %#v\n  actual: %#v", expected, actual)
	}

	if akte != "" {
		t.Fatalf("akte:\nexpected empty but got %q", akte)
	}
}

func TestReadWithAkte(t1 *testing.T) {
	t := test_logz.T{T: t1}

	af := test_metadatei_io.FixtureFactoryReadWriteCloser(
		map[string]string{
			"fa8242e99f48966ca514092b4233b446851f42b57ad5031bf133e1dd76787f3e": "the body\n",
			"036a8e44e472523c0306946f2712f372c234f8a24532e933f1509ae4db0da064": "the body",
		},
	)

	actual, akte := readFormat(
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
		Bezeichnung: bezeichnung.Make("the title"),
		Typ:         makeAkteExt(t, "md"),
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

	expectedAkte := "the body\n"

	if expectedAkte != akte {
		t.Fatalf("akte:\nexpected: %#v\n  actual: %#v", expectedAkte, akte)
	}
}

// func TestReadMultilineBezeichnung(t *testing.T) {
// 	zt := makeText(
// 		t,
// 		`---
// # the title
//   continues here
// - tag1
// - tag2
// - tag3
// ! text/plain
// ---

// the body`,
// 	)

// 	expected := Text{
// 		Metadatei: Metadatei{
// 			Bezeichnung: "the title continues here",
// 			Etiketten: []etikett.Etikett{
// 				etikett.Etikett{Value: "tag1"},
// 				etikett.Etikett{Value: "tag2"},
// 				etikett.Etikett{Value: "tag3"},
// 			},
// 			AkteExt: "text/plain",
// 		},
// 		Akte: akte{
// 			buffer: bytes.NewBufferString("the body"),
// 		},
// 	}

// 	if !zt.Equals(expected) {
// 		t.Fatalf("\nexpected: %#v\nactual: %#v", expected, zt)
// 	}
// }

// func TestReadImplicitExt(t *testing.T) {
// 	zt := makeText(
// 		t,
// 		`---
// # the title
// - tag1
// - tag2
// - tag3
// ! the_file.png
// ---`,
// 	)

// 	panic(errors.Errorf("%#v", zt))

// 	expected := Text{
// 		Metadatei: Metadatei{
// 			Bezeichnung: "the title",
// 			Etiketten: []etikett.Etikett{
// 				etikett.Etikett{Value: "tag1"},
// 				etikett.Etikett{Value: "tag2"},
// 				etikett.Etikett{Value: "tag3"},
// 			},
// 			AkteExt: "the_file.png",
// 		},
// 	}

// 	if !zt.Equals(expected) {
// 		t.Fatalf("\nexpected: %#v\nactual: %#v", expected, zt)
// 	}
// }

// func TestWrite(t *testing.T) {
// 	zt := &Text{
// 		Metadatei: Metadatei{
// 			Bezeichnung: "the title",
// 			Etiketten: []etikett.Etikett{
// 				etikett.Etikett{Value: "tag1"},
// 				etikett.Etikett{Value: "tag2"},
// 				etikett.Etikett{Value: "tag3"},
// 			},
// 			AkteExt: "text/plain",
// 		},
// 		Akte: akte{
// 			buffer: bytes.NewBufferString("the body"),
// 		},
// 	}

// 	expected := `---
// # the title
// - tag1
// - tag2
// - tag3
// ! text/plain
// ---

// the body`

// 	var err error

// 	actual := &strings.Builder{}

// 	if _, err = zt.WriteTo(actual); err != nil {
// 		t.Fatalf("%s", err)
// 		return
// 	}

// 	if expected != actual.String() {
// 		t.Fatalf("\nexpected: %q\nactual:   %q", expected, actual.String())
// 	}
// }
