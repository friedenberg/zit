package zettel

import (
	"testing"
)

func TestReadWithoutAkte(t *testing.T) {
	actual, akte := readFormat(
		t,
		Text{},
		`---
# the title
- tag1
- tag2
- tag3
! md
---
`,
	)

	expected := Zettel{
		Bezeichnung: "the title",
		Etiketten: makeEtiketten(t,
			"tag1",
			"tag2",
			"tag3",
		),
		Typ: makeAkteExt(t, "md"),
	}

	if !actual.Equals(expected) {
		t.Fatalf("zettel:\nexpected: %#v\n  actual: %#v", expected, actual)
	}

	if akte != "" {
		t.Fatalf("akte:\nexpected empty but got %q", akte)
	}
}

func TestReadWithoutAkteWithMultilineBezeichnung(t *testing.T) {
	actual, akte := readFormat(
		t,
		Text{},
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

	expected := Zettel{
		Bezeichnung: "the title continues",
		Etiketten: makeEtiketten(t,
			"tag1",
			"tag2",
			"tag3",
		),
		Typ: makeAkteExt(t, "md"),
	}

	if !actual.Equals(expected) {
		t.Fatalf("zettel:\nexpected: %#v\n  actual: %#v", expected, actual)
	}

	if akte != "" {
		t.Fatalf("akte:\nexpected empty but got %q", akte)
	}
}

func TestReadWithAkte(t *testing.T) {
	actual, akte := readFormat(
		t,
		Text{},
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

	expected := Zettel{
		Bezeichnung: "the title",
		Etiketten: makeEtiketten(t,
			"tag1",
			"tag2",
			"tag3",
		),
		Typ: makeAkteExt(t, "md"),
	}

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
