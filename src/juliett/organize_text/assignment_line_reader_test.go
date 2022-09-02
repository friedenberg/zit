package organize_text

import (
	"os"
	"strings"
	"testing"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/test_logz"
	"github.com/friedenberg/zit/src/delta/etikett"
)

func TestMain(m *testing.M) {
	errors.SetTesting()
	code := m.Run()
	os.Exit(code)
}

func TestAssignmentLineReaderOneHeadingNoZettels(t *testing.T) {
	input :=
		`# wow
    `

	sr := strings.NewReader(input)
	sub := assignmentLineReader{}

	n, err := sub.ReadFrom(sr)

	if n == 0 {
		test_logz.Errorf(t, "expected read amount to be greater than 0")
	}

	if err != nil {
		test_logz.Errorf(t, "expected no error but got %q", err)
	}

	{
		expected := etikett.NewSet(etikett.Etikett{Value: "wow"})
		actual := sub.root.children[0].etiketten

		if !actual.Equals(*expected) {
			test_logz.Errorf(t, "\nexpected: %s\n  actual: %s", expected, actual)
		}
	}
}

func TestAssignmentLineReader2Heading2Zettels(t *testing.T) {
	input :=
		`# wow

    - [one] uno
    - [dos] two
    `

	sr := strings.NewReader(input)
	sub := assignmentLineReader{}

	n, err := sub.ReadFrom(sr)

	if n == 0 {
		test_logz.Errorf(t, "expected read amount to be greater than 0")
	}

	if err != nil {
		test_logz.Errorf(t, "expected no error but got %q", err)
	}

	{
		expected := etikett.NewSet(etikett.Etikett{Value: "wow"})
		actual := sub.root.children[0].etiketten

		if !actual.Equals(*expected) {
			test_logz.Errorf(t, "\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	{
		expected := makeZettelSet()
		expected.Add(zettel{
			Hinweis:     "one",
			Bezeichnung: "uno",
		})

		expected.Add(zettel{
			Hinweis:     "dos",
			Bezeichnung: "two",
		})

		actual := sub.root.children[0].named

		if !actual.Equals(expected) {
			test_logz.Errorf(t, "\nexpected: %s\n  actual: %s", expected, actual)
		}
	}
}

func TestAssignmentLineReader1_1Heading2_2Zettels(t *testing.T) {
	input :=
		`# wow

    ## sub-wow

    - [one] uno
    - [dos] two
    `

	sr := strings.NewReader(input)
	sub := assignmentLineReader{}

	n, err := sub.ReadFrom(sr)

	if n == 0 {
		test_logz.Errorf(t, "expected read amount to be greater than 0")
	}

	if err != nil {
		test_logz.Errorf(t, "expected no error but got %q", err)
	}

	{
		expected := etikett.NewSet(etikett.Etikett{Value: "wow"})
		actual := sub.root.children[0].etiketten

		if !actual.Equals(*expected) {
			test_logz.Errorf(t, "\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	{
		expected := etikett.NewSet(etikett.Etikett{Value: "sub-wow"})

		if sub.root != sub.root.children[0].parent {
			test_logz.Fatalf(t, "%v, %v", sub.root, sub.root.children[0].parent)
		}

		l := len(sub.root.children[0].children)

		if l != 1 {
			test_logz.Fatalf(t, "\nexpected: %d\n  actual: %d", 1, l)
		}

		actual := sub.root.children[0].children[0].etiketten

		if !actual.Equals(*expected) {
			test_logz.Errorf(t, "\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	{
		expected := makeZettelSet()
		expected.Add(zettel{
			Hinweis:     "one",
			Bezeichnung: "uno",
		})

		expected.Add(zettel{
			Hinweis:     "dos",
			Bezeichnung: "two",
		})

		actual := sub.root.children[0].children[0].named

		if !actual.Equals(expected) {
			test_logz.Errorf(t, "\nexpected: %s\n  actual: %s", expected, actual)
		}
	}
}

func TestAssignmentLineReader2_1Heading2_2_2Zettels(t *testing.T) {
	input :=
		`# wow

    - [one] uno
    - [dos] two

    ## sub-wow

    - [three] tres
    - [four] quatro

    # cow

    - [one] uno
    - [dos] two
    `

	sr := strings.NewReader(input)
	sub := assignmentLineReader{}

	n, err := sub.ReadFrom(sr)

	if n == 0 {
		test_logz.Errorf(t, "expected read amount to be greater than 0")
	}

	if err != nil {
		test_logz.Errorf(t, "expected no error but got %q", err)
	}

	{
		expected := etikett.NewSet(etikett.Etikett{Value: "wow"})
		actual := sub.root.children[0].etiketten

		if !actual.Equals(*expected) {
			test_logz.Errorf(t, "\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	{
		expected := etikett.NewSet(etikett.Etikett{Value: "sub-wow"})

		l := len(sub.root.children[0].children)
		if l != 1 {
			test_logz.Fatalf(t, "\nexpected: %d\n  actual: %d", 1, l)
		}

		actual := sub.root.children[0].children[0].etiketten

		if !actual.Equals(*expected) {
			test_logz.Errorf(t, "\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	{
		expected := etikett.NewSet(etikett.Etikett{Value: "cow"})
		actual := sub.root.children[1].etiketten

		if !actual.Equals(*expected) {
			test_logz.Errorf(t, "\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	{
		expected := makeZettelSet()
		expected.Add(zettel{
			Hinweis:     "one",
			Bezeichnung: "uno",
		})

		expected.Add(zettel{
			Hinweis:     "dos",
			Bezeichnung: "two",
		})

		actual := sub.root.children[0].named

		if !actual.Equals(expected) {
			test_logz.Errorf(t, "\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	{
		expected := makeZettelSet()
		expected.Add(zettel{
			Hinweis:     "one",
			Bezeichnung: "uno",
		})

		expected.Add(zettel{
			Hinweis:     "dos",
			Bezeichnung: "two",
		})

		actual := sub.root.children[1].named

		if !actual.Equals(expected) {
			test_logz.Errorf(t, "\nexpected: %s\n  actual: %s", expected, actual)
		}
	}
}

func TestAssignmentLineReaderBigCheese(t *testing.T) {
	input :=
		`# task
    - [one] uno
    - [two] dos
    ## priority-1
    ### w-2022-07-09
    - [three] tres
    ###
    - [four] quatro
    ## priority-2
    - [five] cinco
    - [six] seis
    `

	sr := strings.NewReader(input)
	sub := assignmentLineReader{}

	n, err := sub.ReadFrom(sr)

	if n == 0 {
		test_logz.Errorf(t, "expected read amount to be greater than 0")
	}

	if err != nil {
		test_logz.Fatalf(t, "expected no error but got %q", err)
	}

	// `# task
	{
		expected := etikett.NewSet(etikett.Etikett{Value: "task"})
		actual := sub.root.children[0].etiketten

		if !actual.Equals(*expected) {
			test_logz.Errorf(t, "\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	// - [one] uno
	// - [two] dos
	{
		expected := makeZettelSet()
		expected.Add(zettel{
			Hinweis:     "one",
			Bezeichnung: "uno",
		})

		expected.Add(zettel{
			Hinweis:     "two",
			Bezeichnung: "dos",
		})

		actual := sub.root.children[0].named

		if !actual.Equals(expected) {
			test_logz.Errorf(t, "\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	// ## priority-1
	{
		expected := etikett.NewSet(etikett.Etikett{Value: "priority-1"})

		e := 2
		l := len(sub.root.children[0].children)
		if l != e {
			test_logz.Fatalf(t, "\nexpected: %d\n  actual: %d", e, l)
		}

		actual := sub.root.children[0].children[0].etiketten

		if !actual.Equals(*expected) {
			test_logz.Errorf(t, "\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	// ### w-2022-07-09
	{
		expected := etikett.NewSet(etikett.Etikett{Value: "w-2022-07-09"})
		actual := sub.root.children[0].children[0].children[0].etiketten

		if !actual.Equals(*expected) {
			test_logz.Errorf(t, "\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	// - [three] tres
	{
		expected := makeZettelSet()
		expected.Add(zettel{
			Hinweis:     "three",
			Bezeichnung: "tres",
		})
		actual := sub.root.children[0].children[0].children[0].named

		if !actual.Equals(expected) {
			test_logz.Errorf(t, "\nexpected: %q\n  actual: %q", expected, actual)
		}
	}

	// ##
	// - [four] quatro
	{
		expected := makeZettelSet()
		expected.Add(zettel{
			Hinweis:     "four",
			Bezeichnung: "quatro",
		})
		actual := sub.root.children[0].children[0].named

		if !actual.Equals(expected) {
			test_logz.Errorf(t, "\nexpected: %q\n  actual: %q", expected, actual)
		}
	}

	// ## priority-2
	// - [five] cinco
	// - [six] seis
	// `
	{
		expected := makeZettelSet()
		expected.Add(zettel{
			Hinweis:     "five",
			Bezeichnung: "cinco",
		})
		expected.Add(zettel{
			Hinweis:     "six",
			Bezeichnung: "seis",
		})
		actual := sub.root.children[0].children[1].named

		if !actual.Equals(expected) {
			test_logz.Errorf(t, "\nexpected: %q\n  actual: %q", expected, actual)
		}
	}
}
