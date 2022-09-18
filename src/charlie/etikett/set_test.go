package etikett

import (
	"testing"

	"github.com/friedenberg/zit/src/bravo/test_logz"
)

func TestNormalize(t *testing.T) {
	type testEntry struct {
		ac Set
		ex Set
	}

	testEntries := map[string]testEntry{
		"removes all": testEntry{
			ac: MakeSet(
				Etikett{Value: "project-2021-zit"},
				Etikett{Value: "project-2021-zit-test"},
				Etikett{Value: "project-2021-zit-ewwwwww"},
				Etikett{Value: "project-archive-task-done"},
			),
			ex: MakeSet(
				Etikett{Value: "project-2021-zit-test"},
				Etikett{Value: "project-2021-zit-ewwwwww"},
				Etikett{Value: "project-archive-task-done"},
			),
		},
		"removes non": testEntry{
			ac: MakeSet(
				Etikett{Value: "project-2021-zit"},
				Etikett{Value: "project-2021-zit-test"},
				Etikett{Value: "project-2021-zit-ewwwwww"},
				Etikett{Value: "zz-archive-task-done"},
			),
			ex: MakeSet(
				Etikett{Value: "project-2021-zit-test"},
				Etikett{Value: "project-2021-zit-ewwwwww"},
				Etikett{Value: "zz-archive-task-done"},
			),
		},
	}

	for d, te := range testEntries {
		t.Run(
			d,
			func(t *testing.T) {
				ac := te.ac.WithRemovedCommonPrefixes()

				if !ac.Equals(te.ex) {
					test_logz.Errorf(
						test_logz.T{T: t},
						"removing prefixes doesn't match:\nexpected: %q\n  actual: %q",
						te.ex,
						ac,
					)
				}
			},
		)
	}
}

func TestRemovePrefixes(t *testing.T) {
	type testEntry struct {
		ac     Set
		ex     Set
		prefix string
	}

	testEntries := map[string]testEntry{
		"removes all": testEntry{
			ac: MakeSet(
				Etikett{Value: "project-2021-zit"},
				Etikett{Value: "project-2021-zit-test"},
				Etikett{Value: "project-2021-zit-ewwwwww"},
				Etikett{Value: "project-archive-task-done"},
			),
			ex:     MakeSet(),
			prefix: "project",
		},
		"removes non": testEntry{
			ac: MakeSet(
				Etikett{Value: "project-2021-zit"},
				Etikett{Value: "project-2021-zit-test"},
				Etikett{Value: "project-2021-zit-ewwwwww"},
				Etikett{Value: "zz-archive-task-done"},
			),
			ex: MakeSet(
				Etikett{Value: "project-2021-zit"},
				Etikett{Value: "project-2021-zit-test"},
				Etikett{Value: "project-2021-zit-ewwwwww"},
				Etikett{Value: "zz-archive-task-done"},
			),
			prefix: "xx",
		},
		"removes one": testEntry{
			ac: MakeSet(
				Etikett{Value: "project-2021-zit"},
				Etikett{Value: "project-2021-zit-test"},
				Etikett{Value: "project-2021-zit-ewwwwww"},
				Etikett{Value: "zz-archive-task-done"},
			),
			ex: MakeSet(
				Etikett{Value: "project-2021-zit"},
				Etikett{Value: "project-2021-zit-test"},
				Etikett{Value: "project-2021-zit-ewwwwww"},
			),
			prefix: "zz",
		},
		"removes most": testEntry{
			ac: MakeSet(
				Etikett{Value: "project-2021-zit"},
				Etikett{Value: "project-2021-zit-test"},
				Etikett{Value: "project-2021-zit-ewwwwww"},
				Etikett{Value: "zz-archive-task-done"},
			),
			ex: MakeSet(
				Etikett{Value: "zz-archive-task-done"},
			),
			prefix: "project",
		},
	}

	for d, te := range testEntries {
		t.Run(
			d,
			func(t *testing.T) {
				assertSetRemovesPrefixes(t, te.ac, te.ex, te.prefix)
			},
		)
	}
}

func TestExpandedRight(t *testing.T) {
	s := MakeSet(
		Etikett{Value: "project-2021-zit"},
		Etikett{Value: "zz-archive-task-done"},
	)

	ex := s.Expanded(ExpanderRight{})

	expected := []string{
		"project",
		"project-2021",
		"project-2021-zit",
		"zz",
		"zz-archive",
		"zz-archive-task",
		"zz-archive-task-done",
	}

	actual := ex.SortedString()

	if !stringSliceEquals(actual, expected) {
		t.Errorf(
			"expanded tags don't match:\nexpected: %q\n  actual: %q",
			expected,
			actual,
		)
	}
}

func TestPrefixIntersection(t *testing.T) {
	s := MakeSet(
		Etikett{Value: "project-2021-zit"},
		Etikett{Value: "zz-archive-task-done"},
	)

	ex := s.IntersectPrefixes(MakeSet(Etikett{Value: "project"}))

	expected := []string{
		"project-2021-zit",
	}

	actual := ex.SortedString()

	if !stringSliceEquals(actual, expected) {
		t.Errorf(
			"expanded tags don't match:\nexpected: %q\n  actual: %q",
			expected,
			actual,
		)
	}
}

// func TestExpansionRight(t *testing.T) {
// 	e := Etikett{Value: "this-is-a-tag"}
// 	ex := e.Expanded(ExpanderRight{})
// 	expected := []string{
// 		"this",
// 		"this-is",
// 		"this-is-a",
// 		"this-is-a-tag",
// 	}

// 	actual := ex.SortedString()

// 	if !stringSliceEquals(actual, expected) {
// 		t.Errorf(
// 			"expanded tags don't match:\nexpected: %q\n  actual: %q",
// 			expected,
// 			actual,
// 		)
// 	}
// }

func TestDelta1(t *testing.T) {
	a := MakeSet(
		Etikett{Value: "project-2021-zit"},
		Etikett{Value: "task-todo"},
	)

	b := MakeSet(
		Etikett{Value: "project-2021-zit"},
		Etikett{Value: "zz-archive-task-done"},
	)

	d := MakeSetDelta(a, b)

	c_expected := MakeSet(
		Etikett{Value: "zz-archive-task-done"},
	)

	if !c_expected.Equals(d.Added) {
		t.Errorf("expected\n%s\nactual:\n%s", c_expected, d.Added)
	}

	d_expected := MakeSet(
		Etikett{Value: "task-todo"},
	)

	if !d_expected.Equals(d.Removed) {
		t.Errorf("expected\n%s\nactual:\n%s", d_expected, d.Removed)
	}
}
