package kennung

import (
	"testing"

	"code.linenisgreat.com/zit/src/bravo/expansion"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/bravo/test_logz"
	"code.linenisgreat.com/zit/src/delta/collections_delta"
)

func TestNormalize(t *testing.T) {
	type testEntry struct {
		ac EtikettSet
		ex EtikettSet
	}

	testEntries := map[string]testEntry{
		"removes all": {
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
		"removes non": {
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
		"removes right order": {
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

				if !EtikettSetEquals(ac, te.ex) {
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
		"removes all": {
			ac: MakeEtikettSet(
				MustEtikett("project-2021-zit"),
				MustEtikett("project-2021-zit-test"),
				MustEtikett("project-2021-zit-ewwwwww"),
				MustEtikett("project-archive-task-done"),
			),
			ex:     MakeEtikettSet(),
			prefix: "project",
		},
		"removes non": {
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
		"removes one": {
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
		"removes most": {
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

	ex := Expanded(s, expansion.ExpanderRight)

	expected := []string{
		"project",
		"project-2021",
		"project-2021-zit",
		"zz",
		"zz-archive",
		"zz-archive-task",
		"zz-archive-task-done",
	}

	actual := iter.SortedStrings[Etikett](ex)

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

	ex := IntersectPrefixes(s, MakeEtikettSet(MustEtikett("project")))

	expected := []string{
		"project-2021-zit",
	}

	actual := iter.SortedStrings[Etikett](ex)

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
	from := MakeEtikettSet(
		MustEtikett("project-2021-zit"),
		MustEtikett("task-todo"),
	)

	to := MakeEtikettSet(
		MustEtikett("project-2021-zit"),
		MustEtikett("zz-archive-task-done"),
	)

	d := collections_delta.MakeSetDelta[Etikett](from, to)

	c_expected := MakeEtikettSet(
		MustEtikett("zz-archive-task-done"),
	)

	if !iter.SetEquals[Etikett](c_expected, d.GetAdded()) {
		t.Errorf("expected\n%s\nactual:\n%s", c_expected, d.GetAdded())
	}

	d_expected := MakeEtikettSet(
		MustEtikett("task-todo"),
	)

	if !iter.SetEquals[Etikett](d_expected, d.GetRemoved()) {
		t.Errorf("expected\n%s\nactual:\n%s", d_expected, d.GetRemoved())
	}
}
