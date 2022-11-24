package kennung

import (
	"testing"

	"github.com/friedenberg/zit/src/bravo/test_logz"
)

func TestNormalize(t *testing.T) {
	type testEntry struct {
		ac EtikettSet
		ex EtikettSet
	}

	testEntries := map[string]testEntry{
		"removes all": testEntry{
			ac: MakeEtikettSet(
				MustEtikett("project-2021-zit"),
				MustEtikett("project-2021-zit-test"),
				MustEtikett("project-2021-zit-ewwwwww"),
				MustEtikett("project-archive-task-done"),
			),
			ex: MakeEtikettSet(
				MustEtikett("project-2021-zit-test"),
				MustEtikett("project-2021-zit-ewwwwww"),
				MustEtikett("project-archive-task-done"),
			),
		},
		"removes non": testEntry{
			ac: MakeEtikettSet(
				MustEtikett("project-2021-zit"),
				MustEtikett("project-2021-zit-test"),
				MustEtikett("project-2021-zit-ewwwwww"),
				MustEtikett("zz-archive-task-done"),
			),
			ex: MakeEtikettSet(
				MustEtikett("project-2021-zit-test"),
				MustEtikett("project-2021-zit-ewwwwww"),
				MustEtikett("zz-archive-task-done"),
			),
		},
		"removes right order": testEntry{
			ac: MakeEtikettSet(
				MustEtikett("priority"),
				MustEtikett("priority-1"),
			),
			ex: MakeEtikettSet(
				MustEtikett("priority-1"),
			),
		},
	}

	for d, te := range testEntries {
		t.Run(
			d,
			func(t1 *testing.T) {
				t := test_logz.T{T: t1}
				ac := WithRemovedCommonPrefixes(te.ac)

				if !ac.Equals(te.ex) {
					t.Errorf(
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
		ac     EtikettSet
		ex     EtikettSet
		prefix string
	}

	testEntries := map[string]testEntry{
		"removes all": testEntry{
			ac: MakeEtikettSet(
				MustEtikett("project-2021-zit"),
				MustEtikett("project-2021-zit-test"),
				MustEtikett("project-2021-zit-ewwwwww"),
				MustEtikett("project-archive-task-done"),
			),
			ex:     MakeEtikettSet(),
			prefix: "project",
		},
		"removes non": testEntry{
			ac: MakeEtikettSet(
				MustEtikett("project-2021-zit"),
				MustEtikett("project-2021-zit-test"),
				MustEtikett("project-2021-zit-ewwwwww"),
				MustEtikett("zz-archive-task-done"),
			),
			ex: MakeEtikettSet(
				MustEtikett("project-2021-zit"),
				MustEtikett("project-2021-zit-test"),
				MustEtikett("project-2021-zit-ewwwwww"),
				MustEtikett("zz-archive-task-done"),
			),
			prefix: "xx",
		},
		"removes one": testEntry{
			ac: MakeEtikettSet(
				MustEtikett("project-2021-zit"),
				MustEtikett("project-2021-zit-test"),
				MustEtikett("project-2021-zit-ewwwwww"),
				MustEtikett("zz-archive-task-done"),
			),
			ex: MakeEtikettSet(
				MustEtikett("project-2021-zit"),
				MustEtikett("project-2021-zit-test"),
				MustEtikett("project-2021-zit-ewwwwww"),
			),
			prefix: "zz",
		},
		"removes most": testEntry{
			ac: MakeEtikettSet(
				MustEtikett("project-2021-zit"),
				MustEtikett("project-2021-zit-test"),
				MustEtikett("project-2021-zit-ewwwwww"),
				MustEtikett("zz-archive-task-done"),
			),
			ex: MakeEtikettSet(
				MustEtikett("zz-archive-task-done"),
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
	s := MakeEtikettSet(
		MustEtikett("project-2021-zit"),
		MustEtikett("zz-archive-task-done"),
	)

	ex := Expanded(s, ExpanderRight)

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
	s := MakeEtikettSet(
		MustEtikett("project-2021-zit"),
		MustEtikett("zz-archive-task-done"),
	)

	ex := s.IntersectPrefixes(MakeEtikettSet(MustEtikett("project")))

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
// 	e := MustEtikett("this-is-a-tag")
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
	a := MakeEtikettSet(
		MustEtikett("project-2021-zit"),
		MustEtikett("task-todo"),
	)

	b := MakeEtikettSet(
		MustEtikett("project-2021-zit"),
		MustEtikett("zz-archive-task-done"),
	)

	d := MakeSetDelta(a, b)

	c_expected := MakeEtikettSet(
		MustEtikett("zz-archive-task-done"),
	)

	if !c_expected.Equals(d.Added) {
		t.Errorf("expected\n%s\nactual:\n%s", c_expected, d.Added)
	}

	d_expected := MakeEtikettSet(
		MustEtikett("task-todo"),
	)

	if !d_expected.Equals(d.Removed) {
		t.Errorf("expected\n%s\nactual:\n%s", d_expected, d.Removed)
	}
}
