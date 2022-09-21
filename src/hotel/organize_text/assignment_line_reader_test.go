package organize_text

import (
	"os"
	"strings"
	"testing"

	"github.com/friedenberg/zit/src/alfa/bezeichnung"
	"github.com/friedenberg/zit/src/alfa/errors"
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

func TestAssignmentLineReaderOneHeadingNoZettels(t *testing.T) {
	input :=
		`# wow
    `

	sr := strings.NewReader(input)
	sub := assignmentLineReader{}

	n, err := sub.ReadFrom(sr)

	if n == 0 {
		test_logz.Errorf(test_logz.T{T: t}, "expected read amount to be greater than 0")
	}

	if err != nil {
		test_logz.Errorf(test_logz.T{T: t}, "expected no error but got %q", err)
	}

	{
		expected := etikett.NewSet(etikett.Etikett{Value: "wow"})
		actual := sub.root.children[0].etiketten

		if !actual.Equals(*expected) {
			test_logz.Errorf(test_logz.T{T: t}, "\nexpected: %s\n  actual: %s", expected, actual)
		}
	}
}

func TestAssignmentLineReader2Heading2Zettels(t *testing.T) {
	input :=
		`# wow

    - [one/wow] uno
    - [dos/wow] two/wow
    `

	sr := strings.NewReader(input)
	sub := assignmentLineReader{}

	n, err := sub.ReadFrom(sr)

	if n == 0 {
		test_logz.Errorf(test_logz.T{T: t}, "expected read amount to be greater than 0")
	}

	if err != nil {
		test_logz.Errorf(test_logz.T{T: t}, "expected no error but got %q", err)
	}

	{
		expected := etikett.NewSet(etikett.Etikett{Value: "wow"})
		actual := sub.root.children[0].etiketten

		if !actual.Equals(*expected) {
			test_logz.Errorf(test_logz.T{T: t}, "\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	{
		expected := makeZettelSet()
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
			test_logz.Errorf(test_logz.T{T: t}, "\nexpected: %s\n  actual: %s", expected, actual)
		}
	}
}

func TestAssignmentLineReader1_1Heading2_2Zettels(t *testing.T) {
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
		test_logz.Errorf(test_logz.T{T: t}, "expected read amount to be greater than 0")
	}

	if err != nil {
		test_logz.Errorf(test_logz.T{T: t}, "expected no error but got %q", err)
	}

	{
		expected := etikett.NewSet(etikett.Etikett{Value: "wow"})
		actual := sub.root.children[0].etiketten

		if !actual.Equals(*expected) {
			test_logz.Errorf(test_logz.T{T: t}, "\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	{
		expected := etikett.NewSet(etikett.Etikett{Value: "sub-wow"})

		if sub.root != sub.root.children[0].parent {
			test_logz.Fatalf(test_logz.T{T: t}, "%v, %v", sub.root, sub.root.children[0].parent)
		}

		l := len(sub.root.children[0].children)

		if l != 1 {
			test_logz.Fatalf(test_logz.T{T: t}, "\nexpected: %d\n  actual: %d", 1, l)
		}

		actual := sub.root.children[0].children[0].etiketten

		if !actual.Equals(*expected) {
			test_logz.Errorf(test_logz.T{T: t}, "\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	{
		expected := makeZettelSet()
		expected.Add(zettel{
			Hinweis:     makeHinweis(t, "one/wow"),
			Bezeichnung: makeBez(t, "uno"),
		})

		expected.Add(zettel{
			Hinweis:     makeHinweis(t, "dos/wow"),
			Bezeichnung: makeBez(t, "two/wow"),
		})

		actual := sub.root.children[0].children[0].named

		if !actual.Equals(expected) {
			test_logz.Errorf(test_logz.T{T: t}, "\nexpected: %s\n  actual: %s", expected, actual)
		}
	}
}

func TestAssignmentLineReader2_1Heading2_2_2Zettels(t *testing.T) {
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
		test_logz.Errorf(test_logz.T{T: t}, "expected read amount to be greater than 0")
	}

	if err != nil {
		test_logz.Errorf(test_logz.T{T: t}, "expected no error but got %q", err)
	}

	{
		expected := etikett.NewSet(etikett.Etikett{Value: "wow"})
		actual := sub.root.children[0].etiketten

		if !actual.Equals(*expected) {
			test_logz.Errorf(test_logz.T{T: t}, "\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	{
		expected := etikett.NewSet(etikett.Etikett{Value: "sub-wow"})

		l := len(sub.root.children[0].children)
		if l != 1 {
			test_logz.Fatalf(test_logz.T{T: t}, "\nexpected: %d\n  actual: %d", 1, l)
		}

		actual := sub.root.children[0].children[0].etiketten

		if !actual.Equals(*expected) {
			test_logz.Errorf(test_logz.T{T: t}, "\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	{
		expected := etikett.NewSet(etikett.Etikett{Value: "cow"})
		actual := sub.root.children[1].etiketten

		if !actual.Equals(*expected) {
			test_logz.Errorf(test_logz.T{T: t}, "\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	{
		expected := makeZettelSet()
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
			test_logz.Errorf(test_logz.T{T: t}, "\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	{
		expected := makeZettelSet()
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
			test_logz.Errorf(test_logz.T{T: t}, "\nexpected: %s\n  actual: %s", expected, actual)
		}
	}
}

func TestAssignmentLineReader2_1Heading2_2_2ZettelsOffset(t *testing.T) {
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
		test_logz.Errorf(test_logz.T{T: t}, "expected read amount to be greater than 0")
	}

	if err != nil {
		test_logz.Errorf(test_logz.T{T: t}, "expected no error but got %q", err)
	}

	{
		expected := etikett.NewSet(etikett.Etikett{Value: "sub-wow"})
		actual := sub.root.children[0].etiketten

		if !actual.Equals(*expected) {
			test_logz.Errorf(test_logz.T{T: t}, "\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	{
		expected := etikett.NewSet(etikett.Etikett{Value: "sub-cow"})

		l := len(sub.root.children)
    expLen := 2
		if l != expLen {
			test_logz.Fatalf(test_logz.T{T: t}, "\nexpected: %d\n  actual: %d", expLen, l)
		}

		actual := sub.root.children[1].etiketten

		if !actual.Equals(*expected) {
			test_logz.Errorf(test_logz.T{T: t}, "\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	{
		expected := makeZettelSet()
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
			test_logz.Errorf(test_logz.T{T: t}, "\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	{
		expected := makeZettelSet()
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
			test_logz.Errorf(test_logz.T{T: t}, "\nexpected: %s\n  actual: %s", expected, actual)
		}
	}
}

func TestAssignmentLineReaderBigCheese(t *testing.T) {
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
		test_logz.Errorf(test_logz.T{T: t}, "expected read amount to be greater than 0")
	}

	if err != nil {
		test_logz.Fatalf(test_logz.T{T: t}, "expected no error but got %q", err)
	}

	// `# task
	{
		expected := etikett.NewSet(etikett.Etikett{Value: "task"})
		actual := sub.root.children[0].etiketten

		if !actual.Equals(*expected) {
			test_logz.Errorf(test_logz.T{T: t}, "\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	// - [one/wow] uno
	// - [two/wow] dos/wow
	{
		expected := makeZettelSet()
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
			test_logz.Errorf(test_logz.T{T: t}, "\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	// ## priority-1
	{
		expected := etikett.NewSet(etikett.Etikett{Value: "priority-1"})

		e := 2
		l := len(sub.root.children[0].children)
		if l != e {
			test_logz.Fatalf(test_logz.T{T: t}, "\nexpected: %d\n  actual: %d", e, l)
		}

		actual := sub.root.children[0].children[0].etiketten

		if !actual.Equals(*expected) {
			test_logz.Errorf(test_logz.T{T: t}, "\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	// ### w-2022-07-09
	{
		expected := etikett.NewSet(etikett.Etikett{Value: "w-2022-07-09"})
		actual := sub.root.children[0].children[0].children[0].etiketten

		if !actual.Equals(*expected) {
			test_logz.Errorf(test_logz.T{T: t}, "\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	// - [three/wow] tres
	{
		expected := makeZettelSet()
		expected.Add(zettel{
			Hinweis:     makeHinweis(t, "three/wow"),
			Bezeichnung: makeBez(t, "tres"),
		})
		actual := sub.root.children[0].children[0].children[0].named

		if !actual.Equals(expected) {
			test_logz.Errorf(test_logz.T{T: t}, "\nexpected: %q\n  actual: %q", expected, actual)
		}
	}

	// ##
	// - [four/wow] quatro
	{
		expected := makeZettelSet()
		expected.Add(zettel{
			Hinweis:     makeHinweis(t, "four/wow"),
			Bezeichnung: makeBez(t, "quatro"),
		})
		actual := sub.root.children[0].children[0].named

		if !actual.Equals(expected) {
			test_logz.Errorf(test_logz.T{T: t}, "\nexpected: %q\n  actual: %q", expected, actual)
		}
	}

	// ## priority-2
	// - [five/wow] cinco
	// - [six/wow] seis
	// `
	{
		expected := makeZettelSet()
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
			test_logz.Errorf(test_logz.T{T: t}, "\nexpected: %q\n  actual: %q", expected, actual)
		}
	}
}
