package organize_text

import (
	"os"
	"strings"
	"testing"

	"github.com/friedenberg/zit/src/alfa/bezeichnung"
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/bravo/test_logz"
	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/charlie/hinweis"
)

func TestMain(m *testing.M) {
	errors.SetTesting()
	code := m.Run()
	os.Exit(code)
}

func makeHinweis(t *testing.T, v string) (h hinweis.Hinweis) {
	var err error

	if err = h.Set(v); err != nil {
		t.Fatalf("%s", err)
	}

	return
}

func makeBez(t *testing.T, v string) (b bezeichnung.Bezeichnung) {
	var err error

	if err = b.Set(v); err != nil {
		t.Fatalf("%s", err)
	}

	return
}

func TestAssignmentLineReaderOneHeadingNoZettels(t1 *testing.T) {
	t := test_logz.T{T: t1}
	input :=
		`# wow
    `

	sr := strings.NewReader(input)
	sub := assignmentLineReader{}

	n, err := sub.ReadFrom(sr)

	if n == 0 {
		t.Errorf("expected read amount to be greater than 0")
	}

	if err != nil {
		t.Fatalf("expected no error but got %q", err)
	}

	{
		expected := etikett.MakeSet(etikett.Etikett{Value: "wow"})

		if len(sub.root.children) < 1 {
			t.Fatalf("expected exactly 1 child")
		}

		actual := sub.root.children[0].etiketten

		if !actual.Equals(expected) {
			t.Errorf("\nexpected: %s\n  actual: %s", expected, actual)
		}
	}
}

func TestAssignmentLineReader2Heading2Zettels(t *testing.T) {
	t1 := test_logz.T{T: t}

	input :=
		`# wow

    - [one/wow] uno
    - [dos/wow] two/wow
    `

	sr := strings.NewReader(input)
	sub := assignmentLineReader{}

	n, err := sub.ReadFrom(sr)

	if n == 0 {
		t1.Errorf("expected read amount to be greater than 0")
	}

	if err != nil {
		t1.Fatalf("expected no error but got %q", err)
	}

	{
		expected := etikett.MakeSet(etikett.Etikett{Value: "wow"})
		actual := sub.root.children[0].etiketten

		if !actual.Equals(expected) {
			t1.Errorf("\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	{
		expected := collections.MakeMutableValueSet[zettel, *zettel]()
		expected.Add(zettel{
			Hinweis:     makeHinweis(t, "one/wow"),
			Bezeichnung: makeBez(t, "uno"),
		})

		expected.Add(zettel{
			Hinweis:     makeHinweis(t, "dos/wow"),
			Bezeichnung: makeBez(t, "two/wow"),
		})

		actual := sub.root.children[0].named

		if !actual.Equals(expected) {
			t1.Errorf("\nexpected: %s\n  actual: %s", expected, actual)
		}
	}
}

func TestAssignmentLineReader1_1Heading2_2Zettels(t1 *testing.T) {
	t := test_logz.T{T: t1}

	input :=
		`# wow

    ## sub-wow

    - [one/wow] uno
    - [dos/wow] two/wow
    `

	sr := strings.NewReader(input)
	sub := assignmentLineReader{}

	n, err := sub.ReadFrom(sr)

	if n == 0 {
		t.Errorf("expected read amount to be greater than 0")
	}

	if err != nil {
		t.Fatalf("expected no error but got %q", err)
	}

	{
		expected := etikett.MakeSet(etikett.Etikett{Value: "wow"})
		actual := sub.root.children[0].etiketten

		if !actual.Equals(expected) {
			t.Errorf("\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	{
		expected := etikett.MakeSet(etikett.Etikett{Value: "sub-wow"})

		if sub.root != sub.root.children[0].parent {
			t.Fatalf("%v, %v", sub.root, sub.root.children[0].parent)
		}

		l := len(sub.root.children[0].children)

		if l != 1 {
			t.Fatalf("\nexpected: %d\n  actual: %d", 1, l)
		}

		actual := sub.root.children[0].children[0].etiketten

		if !actual.Equals(expected) {
			t.Errorf("\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	{
		expected := collections.MakeMutableValueSet[zettel, *zettel]()
		expected.Add(zettel{
			Hinweis:     makeHinweis(t1, "one/wow"),
			Bezeichnung: makeBez(t1, "uno"),
		})

		expected.Add(zettel{
			Hinweis:     makeHinweis(t1, "dos/wow"),
			Bezeichnung: makeBez(t1, "two/wow"),
		})

		actual := sub.root.children[0].children[0].named

		if !actual.Equals(expected) {
			t1.Errorf("\nexpected: %v\n  actual: %v", expected, actual)
		}
	}
}

func TestAssignmentLineReader2_1Heading2_2_2Zettels(t *testing.T) {
	t1 := test_logz.T{T: t}

	input :=
		`# wow

    - [one/wow] uno
    - [dos/wow] two/wow

    ## sub-wow

    - [three/wow] tres
    - [four/wow] quatro

    # cow

    - [one/wow] uno
    - [dos/wow] two/wow
    `

	sr := strings.NewReader(input)
	sub := assignmentLineReader{}

	n, err := sub.ReadFrom(sr)

	if n == 0 {
		t1.Errorf("expected read amount to be greater than 0")
	}

	if err != nil {
		t1.Fatalf("expected no error but got %q", err)
	}

	{
		expected := etikett.MakeSet(etikett.Etikett{Value: "wow"})
		actual := sub.root.children[0].etiketten

		if !actual.Equals(expected) {
			t1.Errorf("\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	{
		expected := etikett.MakeSet(etikett.Etikett{Value: "sub-wow"})

		l := len(sub.root.children[0].children)
		if l != 1 {
			t1.Fatalf("\nexpected: %d\n  actual: %d", 1, l)
		}

		actual := sub.root.children[0].children[0].etiketten

		if !actual.Equals(expected) {
			t1.Errorf("\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	{
		expected := etikett.MakeSet(etikett.Etikett{Value: "cow"})
		actual := sub.root.children[1].etiketten

		if !actual.Equals(expected) {
			t1.Errorf("\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	{
		expected := collections.MakeMutableValueSet[zettel, *zettel]()
		expected.Add(zettel{
			Hinweis:     makeHinweis(t, "one/wow"),
			Bezeichnung: makeBez(t, "uno"),
		})

		expected.Add(zettel{
			Hinweis:     makeHinweis(t, "dos/wow"),
			Bezeichnung: makeBez(t, "two/wow"),
		})

		actual := sub.root.children[0].named

		if !actual.Equals(expected) {
			t1.Errorf("\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	{
		expected := collections.MakeMutableValueSet[zettel, *zettel]()
		expected.Add(zettel{
			Hinweis:     makeHinweis(t, "one/wow"),
			Bezeichnung: makeBez(t, "uno"),
		})

		expected.Add(zettel{
			Hinweis:     makeHinweis(t, "dos/wow"),
			Bezeichnung: makeBez(t, "two/wow"),
		})

		actual := sub.root.children[1].named

		if !actual.Equals(expected) {
			t1.Errorf("\nexpected: %s\n  actual: %s", expected, actual)
		}
	}
}

func TestAssignmentLineReader2_1Heading2_2_2ZettelsOffset(t *testing.T) {
	t1 := test_logz.T{T: t}

	input :=
		`
    - [one/wow] uno
    - [dos/wow] two/wow

    ## sub-wow

    - [three/wow] tres
    - [four/wow] quatro

    ## sub-cow

    - [one/wow] uno
    - [dos/wow] two/wow
    `

	sr := strings.NewReader(input)
	sub := assignmentLineReader{}

	n, err := sub.ReadFrom(sr)

	if n == 0 {
		t1.Errorf("expected read amount to be greater than 0")
	}

	if err != nil {
		t1.Fatalf("expected no error but got %q", err)
	}

	{
		expected := etikett.MakeSet(etikett.Etikett{Value: "sub-wow"})
		actual := sub.root.children[0].etiketten

		if !actual.Equals(expected) {
			t1.Errorf("\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	{
		expected := etikett.MakeSet(etikett.Etikett{Value: "sub-cow"})

		l := len(sub.root.children)
		expLen := 2
		if l != expLen {
			t1.Fatalf("\nexpected: %d\n  actual: %d", expLen, l)
		}

		actual := sub.root.children[1].etiketten

		if !actual.Equals(expected) {
			t1.Errorf("\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	{
		expected := collections.MakeMutableValueSet[zettel, *zettel]()
		expected.Add(zettel{
			Hinweis:     makeHinweis(t, "four/wow"),
			Bezeichnung: makeBez(t, "quatro"),
		})

		expected.Add(zettel{
			Hinweis:     makeHinweis(t, "three/wow"),
			Bezeichnung: makeBez(t, "tres"),
		})

		actual := sub.root.children[0].named

		if !actual.Equals(expected) {
			t1.Errorf("\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	{
		expected := collections.MakeMutableValueSet[zettel, *zettel]()
		expected.Add(zettel{
			Hinweis:     makeHinweis(t, "one/wow"),
			Bezeichnung: makeBez(t, "uno"),
		})

		expected.Add(zettel{
			Hinweis:     makeHinweis(t, "dos/wow"),
			Bezeichnung: makeBez(t, "two/wow"),
		})

		actual := sub.root.children[1].named

		if !actual.Equals(expected) {
			t1.Errorf("\nexpected: %s\n  actual: %s", expected, actual)
		}
	}
}

func TestAssignmentLineReaderBigCheese(t *testing.T) {
	t1 := test_logz.T{T: t}

	input :=
		`# task
    - [one/wow] uno
    - [two/wow] dos/wow
    ## priority-1
    ### w-2022-07-09
    - [three/wow] tres
    ###
    - [four/wow] quatro
    ## priority-2
    - [five/wow] cinco
    - [six/wow] seis
    `

	sr := strings.NewReader(input)
	sub := assignmentLineReader{}

	n, err := sub.ReadFrom(sr)

	if n == 0 {
		t1.Errorf("expected read amount to be greater than 0")
	}

	if err != nil {
		t1.Fatalf("expected no error but got %q", err)
	}

	// `# task
	{
		expected := etikett.MakeSet(etikett.Etikett{Value: "task"})
		actual := sub.root.children[0].etiketten

		if !actual.Equals(expected) {
			t1.Errorf("\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	// - [one/wow] uno
	// - [two/wow] dos/wow
	{
		expected := collections.MakeMutableValueSet[zettel, *zettel]()
		expected.Add(zettel{
			Hinweis:     makeHinweis(t, "one/wow"),
			Bezeichnung: makeBez(t, "uno"),
		})

		expected.Add(zettel{
			Hinweis:     makeHinweis(t, "two/wow"),
			Bezeichnung: makeBez(t, "dos/wow"),
		})

		actual := sub.root.children[0].named

		if !actual.Equals(expected) {
			t1.Errorf("\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	// ## priority-1
	{
		expected := etikett.MakeSet(etikett.Etikett{Value: "priority-1"})

		e := 2
		l := len(sub.root.children[0].children)
		if l != e {
			t1.Fatalf("\nexpected: %d\n  actual: %d", e, l)
		}

		actual := sub.root.children[0].children[0].etiketten

		if !actual.Equals(expected) {
			t1.Errorf("\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	// ### w-2022-07-09
	{
		expected := etikett.MakeSet(etikett.Etikett{Value: "w-2022-07-09"})
		actual := sub.root.children[0].children[0].children[0].etiketten

		if !actual.Equals(expected) {
			t1.Errorf("\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	// - [three/wow] tres
	{
		expected := collections.MakeMutableValueSet[zettel, *zettel]()
		expected.Add(zettel{
			Hinweis:     makeHinweis(t, "three/wow"),
			Bezeichnung: makeBez(t, "tres"),
		})
		actual := sub.root.children[0].children[0].children[0].named

		if !actual.Equals(expected) {
			t1.Errorf("\nexpected: %q\n  actual: %q", expected, actual)
		}
	}

	// ##
	// - [four/wow] quatro
	{
		expected := collections.MakeMutableValueSet[zettel, *zettel]()
		expected.Add(zettel{
			Hinweis:     makeHinweis(t, "four/wow"),
			Bezeichnung: makeBez(t, "quatro"),
		})
		actual := sub.root.children[0].children[0].named

		if !actual.Equals(expected) {
			t1.Errorf("\nexpected: %q\n  actual: %q", expected, actual)
		}
	}

	// ## priority-2
	// - [five/wow] cinco
	// - [six/wow] seis
	// `
	{
		expected := collections.MakeMutableValueSet[zettel, *zettel]()
		expected.Add(zettel{
			Hinweis:     makeHinweis(t, "five/wow"),
			Bezeichnung: makeBez(t, "cinco"),
		})
		expected.Add(zettel{
			Hinweis:     makeHinweis(t, "six/wow"),
			Bezeichnung: makeBez(t, "seis"),
		})
		actual := sub.root.children[0].children[1].named

		if !actual.Equals(expected) {
			t1.Errorf("\nexpected: %q\n  actual: %q", expected, actual)
		}
	}
}
