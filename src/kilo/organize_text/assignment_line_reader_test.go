package organize_text

import (
	"os"
	"strings"
	"testing"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/bravo/test_logz"
	"github.com/friedenberg/zit/src/charlie/collections_value"
	"github.com/friedenberg/zit/src/echo/bezeichnung"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/hotel/sku"
)

func TestMain(m *testing.M) {
	errors.SetTesting()
	code := m.Run()
	os.Exit(code)
}

func makeHinweis(t *testing.T, v string) (k *kennung.Kennung2) {
	var err error

	var h kennung.Hinweis

	if err = h.Set(v); err != nil {
		t.Fatalf("%s", err)
	}

	return kennung.MustKennung2(h)
}

func makeBez(t *testing.T, v string) (b bezeichnung.Bezeichnung) {
	var err error

	if err = b.Set(v); err != nil {
		t.Fatalf("%s", err)
	}

	return
}

func makeObjWithHinAndBez(t *testing.T, hin string, bez string) (o *obj) {
	o = &obj{
		Sku: sku.Transacted{
			Metadatei: metadatei.Metadatei{
				Bezeichnung: makeBez(t, bez),
			},
		},
	}

	o.Sku.Kennung.SetWithKennung(makeHinweis(t, hin))

	return
}

func TestAssignmentLineReaderOneHeadingNoZettels(t1 *testing.T) {
	t := test_logz.T{T: t1}
	input := `# wow
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
		expected := kennung.MakeEtikettSet(kennung.MustEtikett("wow"))

		if len(sub.root.children) < 1 {
			t.Fatalf("expected exactly 1 child")
		}

		actual := sub.root.children[0].etiketten

		if !kennung.EtikettSetEquals(
			actual,
			expected,
		) {
			t.Errorf("\nexpected: %s\n  actual: %s", expected, actual)
		}
	}
}

func TestAssignmentLineReader2Heading2Zettels(t *testing.T) {
	t1 := test_logz.T{T: t}

	input := `# wow

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
		expected := kennung.MakeEtikettSet(kennung.MustEtikett("wow"))
		actual := sub.root.children[0].etiketten

		if !kennung.EtikettSetEquals(
			actual,
			expected,
		) {
			t1.Errorf("\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	{
		expected := collections_value.MakeMutableValueSet[*obj](nil)
		expected.Add(makeObjWithHinAndBez(t, "one/wow", "uno"))
		expected.Add(makeObjWithHinAndBez(t, "dos/wow", "two/wow"))

		actual := sub.root.children[0].named

		if !iter.SetEquals[*obj](
			actual,
			expected,
		) {
			t1.Errorf("\nexpected: %s\n  actual: %s", expected, actual)
		}
	}
}

func TestAssignmentLineReader1_1Heading2_2Zettels(t1 *testing.T) {
	t := test_logz.T{T: t1}

	input := `# wow

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
		expected := kennung.MakeEtikettSet(kennung.MustEtikett("wow"))
		actual := sub.root.children[0].etiketten

		if !kennung.EtikettSetEquals(
			actual,
			expected,
		) {
			t.Errorf("\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	{
		expected := kennung.MakeEtikettSet(kennung.MustEtikett("sub-wow"))

		if sub.root != sub.root.children[0].parent {
			t.Fatalf("%v, %v", sub.root, sub.root.children[0].parent)
		}

		l := len(sub.root.children[0].children)

		if l != 1 {
			t.Fatalf("\nexpected: %d\n  actual: %d", 1, l)
		}

		actual := sub.root.children[0].children[0].etiketten

		if !kennung.EtikettSetEquals(
			actual,
			expected,
		) {
			t.Errorf("\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	{
		expected := collections_value.MakeMutableValueSet[*obj](nil)
		expected.Add(makeObjWithHinAndBez(t.T, "one/wow", "uno"))
		expected.Add(makeObjWithHinAndBez(t.T, "dos/wow", "two/wow"))

		actual := sub.root.children[0].children[0].named

		if !iter.SetEquals[*obj](
			actual,
			expected,
		) {
			t1.Errorf("\nexpected: %q\n  actual: %q", expected, actual)
		}
	}
}

func TestAssignmentLineReader2_1Heading2_2_2Zettels(t *testing.T) {
	t1 := test_logz.T{T: t}

	input := `# wow

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
		expected := kennung.MakeEtikettSet(kennung.MustEtikett("wow"))
		actual := sub.root.children[0].etiketten

		if !kennung.EtikettSetEquals(
			actual,
			expected,
		) {
			t1.Errorf("\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	{
		expected := kennung.MakeEtikettSet(kennung.MustEtikett("sub-wow"))

		l := len(sub.root.children[0].children)
		if l != 1 {
			t1.Fatalf("\nexpected: %d\n  actual: %d", 1, l)
		}

		actual := sub.root.children[0].children[0].etiketten

		if !kennung.EtikettSetEquals(
			actual,
			expected,
		) {
			t1.Errorf("\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	{
		expected := kennung.MakeEtikettSet(kennung.MustEtikett("cow"))
		actual := sub.root.children[1].etiketten

		if !kennung.EtikettSetEquals(
			actual,
			expected,
		) {
			t1.Errorf("\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	{
		expected := collections_value.MakeMutableValueSet[*obj](nil)
		expected.Add(makeObjWithHinAndBez(t, "one/wow", "uno"))
		expected.Add(makeObjWithHinAndBez(t, "dos/wow", "two/wow"))

		actual := sub.root.children[0].named

		if !iter.SetEquals[*obj](actual, expected) {
			t1.Errorf("\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	{
		expected := collections_value.MakeMutableValueSet[*obj](nil)
		expected.Add(makeObjWithHinAndBez(t, "one/wow", "uno"))
		expected.Add(makeObjWithHinAndBez(t, "dos/wow", "two/wow"))

		actual := sub.root.children[1].named

		if !iter.SetEquals[*obj](actual, expected) {
			t1.Errorf("\nexpected: %s\n  actual: %s", expected, actual)
		}
	}
}

func TestAssignmentLineReader2_1Heading2_2_2ZettelsOffset(t *testing.T) {
	t1 := test_logz.T{T: t}

	input := `
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
		expected := kennung.MakeEtikettSet(kennung.MustEtikett("sub-wow"))
		actual := sub.root.children[0].etiketten

		if !kennung.EtikettSetEquals(actual, expected) {
			t1.Errorf("\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	{
		expected := kennung.MakeEtikettSet(kennung.MustEtikett("sub-cow"))

		l := len(sub.root.children)
		expLen := 2
		if l != expLen {
			t1.Fatalf("\nexpected: %d\n  actual: %d", expLen, l)
		}

		actual := sub.root.children[1].etiketten

		if !kennung.EtikettSetEquals(actual, expected) {
			t1.Errorf("\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	{
		expected := collections_value.MakeMutableValueSet[*obj](nil)
		expected.Add(makeObjWithHinAndBez(t, "four/wow", "quatro"))
		expected.Add(makeObjWithHinAndBez(t, "three/wow", "tres"))

		actual := sub.root.children[0].named

		if !iter.SetEquals[*obj](actual, expected) {
			t1.Errorf("\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	{
		expected := collections_value.MakeMutableValueSet[*obj](nil)
		expected.Add(makeObjWithHinAndBez(t, "one/wow", "uno"))
		expected.Add(makeObjWithHinAndBez(t, "dos/wow", "two/wow"))

		actual := sub.root.children[1].named

		if !iter.SetEquals[*obj](actual, expected) {
			t1.Errorf("\nexpected: %s\n  actual: %s", expected, actual)
		}
	}
}

func TestAssignmentLineReaderBigCheese(t *testing.T) {
	t1 := test_logz.T{T: t}

	input := `# task
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
		expected := kennung.MakeEtikettSet(kennung.MustEtikett("task"))
		actual := sub.root.children[0].etiketten

		if !kennung.EtikettSetEquals(actual, expected) {
			t1.Errorf("\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	// - [one/wow] uno
	// - [two/wow] dos/wow
	{
		expected := collections_value.MakeMutableValueSet[*obj](nil)
		expected.Add(makeObjWithHinAndBez(t, "one/wow", "uno"))
		expected.Add(makeObjWithHinAndBez(t, "two/wow", "dos/wow"))

		actual := sub.root.children[0].named

		if !iter.SetEquals[*obj](actual, expected) {
			t1.Errorf("\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	// ## priority-1
	{
		expected := kennung.MakeEtikettSet(kennung.MustEtikett("priority-1"))

		e := 2
		l := len(sub.root.children[0].children)
		if l != e {
			t1.Fatalf("\nexpected: %d\n  actual: %d", e, l)
		}

		actual := sub.root.children[0].children[0].etiketten

		if !kennung.EtikettSetEquals(actual, expected) {
			t1.Errorf("\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	// ### w-2022-07-09
	{
		expected := kennung.MakeEtikettSet(kennung.MustEtikett("w-2022-07-09"))
		actual := sub.root.children[0].children[0].children[0].etiketten

		if !kennung.EtikettSetEquals(actual, expected) {
			t1.Errorf("\nexpected: %s\n  actual: %s", expected, actual)
		}
	}

	// - [three/wow] tres
	{
		expected := collections_value.MakeMutableValueSet[*obj](nil)
		expected.Add(makeObjWithHinAndBez(t, "three/wow", "tres"))

		actual := sub.root.children[0].children[0].children[0].named

		if !iter.SetEquals[*obj](actual, expected) {
			t1.Errorf("\nexpected: %q\n  actual: %q", expected, actual)
		}
	}

	// ##
	// - [four/wow] quatro
	{
		expected := collections_value.MakeMutableValueSet[*obj](nil)
		expected.Add(makeObjWithHinAndBez(t, "four/wow", "quatro"))

		actual := sub.root.children[0].children[0].named

		if !iter.SetEquals[*obj](actual, expected) {
			t1.Errorf("\nexpected: %q\n  actual: %q", expected, actual)
		}
	}

	// ## priority-2
	// - [five/wow] cinco
	// - [six/wow] seis
	// `
	{
		expected := collections_value.MakeMutableValueSet[*obj](nil)
		expected.Add(makeObjWithHinAndBez(t, "five/wow", "cinco"))
		expected.Add(makeObjWithHinAndBez(t, "six/wow", "seis"))

		actual := sub.root.children[0].children[1].named

		if !iter.SetEquals[*obj](actual, expected) {
			t1.Errorf("\nexpected: %q\n  actual: %q", expected, actual)
		}
	}
}
