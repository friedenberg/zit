package kennung

import (
	"testing"

	"code.linenisgreat.com/zit/go/zit/src/bravo/expansion"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/test_logz"
	"code.linenisgreat.com/zit/go/zit/src/delta/collections_delta"
)

func TestNormalize(t *testing.T) {
	type testEntry struct {
		ac TagSet
		ex TagSet
	}

	testEntries := map[string]testEntry{
		"removes all": {
			ac: MakeTagSet(
				MustTag("project-2021-zit"),
				MustTag("project-2021-zit-test"),
				MustTag("project-2021-zit-ewwwwww"),
				MustTag("project-archive-task-done"),
			),
			ex: MakeTagSet(
				MustTag("project-2021-zit-test"),
				MustTag("project-2021-zit-ewwwwww"),
				MustTag("project-archive-task-done"),
			),
		},
		"removes non": {
			ac: MakeTagSet(
				MustTag("project-2021-zit"),
				MustTag("project-2021-zit-test"),
				MustTag("project-2021-zit-ewwwwww"),
				MustTag("zz-archive-task-done"),
			),
			ex: MakeTagSet(
				MustTag("project-2021-zit-test"),
				MustTag("project-2021-zit-ewwwwww"),
				MustTag("zz-archive-task-done"),
			),
		},
		"removes right order": {
			ac: MakeTagSet(
				MustTag("priority"),
				MustTag("priority-1"),
			),
			ex: MakeTagSet(
				MustTag("priority-1"),
			),
		},
	}

	for d, te := range testEntries {
		t.Run(
			d,
			func(t1 *testing.T) {
				t := test_logz.T{T: t1}
				ac := WithRemovedCommonPrefixes(te.ac)

				if !TagSetEquals(ac, te.ex) {
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
		ac     TagSet
		ex     TagSet
		prefix string
	}

	testEntries := map[string]testEntry{
		"removes all": {
			ac: MakeTagSet(
				MustTag("project-2021-zit"),
				MustTag("project-2021-zit-test"),
				MustTag("project-2021-zit-ewwwwww"),
				MustTag("project-archive-task-done"),
			),
			ex:     MakeTagSet(),
			prefix: "project",
		},
		"removes non": {
			ac: MakeTagSet(
				MustTag("project-2021-zit"),
				MustTag("project-2021-zit-test"),
				MustTag("project-2021-zit-ewwwwww"),
				MustTag("zz-archive-task-done"),
			),
			ex: MakeTagSet(
				MustTag("project-2021-zit"),
				MustTag("project-2021-zit-test"),
				MustTag("project-2021-zit-ewwwwww"),
				MustTag("zz-archive-task-done"),
			),
			prefix: "xx",
		},
		"removes one": {
			ac: MakeTagSet(
				MustTag("project-2021-zit"),
				MustTag("project-2021-zit-test"),
				MustTag("project-2021-zit-ewwwwww"),
				MustTag("zz-archive-task-done"),
			),
			ex: MakeTagSet(
				MustTag("project-2021-zit"),
				MustTag("project-2021-zit-test"),
				MustTag("project-2021-zit-ewwwwww"),
			),
			prefix: "zz",
		},
		"removes most": {
			ac: MakeTagSet(
				MustTag("project-2021-zit"),
				MustTag("project-2021-zit-test"),
				MustTag("project-2021-zit-ewwwwww"),
				MustTag("zz-archive-task-done"),
			),
			ex: MakeTagSet(
				MustTag("zz-archive-task-done"),
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
	s := MakeTagSet(
		MustTag("project-2021-zit"),
		MustTag("zz-archive-task-done"),
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

	actual := iter.SortedStrings[Tag](ex)

	if !stringSliceEquals(actual, expected) {
		t.Errorf(
			"expanded tags don't match:\nexpected: %q\n  actual: %q",
			expected,
			actual,
		)
	}
}

func TestPrefixIntersection(t *testing.T) {
	s := MakeTagSet(
		MustTag("project-2021-zit"),
		MustTag("zz-archive-task-done"),
	)

	ex := IntersectPrefixes(s, MustTag("project"))

	expected := []string{
		"project-2021-zit",
	}

	actual := iter.SortedStrings(ex)

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
	from := MakeTagSet(
		MustTag("project-2021-zit"),
		MustTag("task-todo"),
	)

	to := MakeTagSet(
		MustTag("project-2021-zit"),
		MustTag("zz-archive-task-done"),
	)

	d := collections_delta.MakeSetDelta[Tag](from, to)

	c_expected := MakeTagSet(
		MustTag("zz-archive-task-done"),
	)

	if !iter.SetEquals[Tag](c_expected, d.GetAdded()) {
		t.Errorf("expected\n%s\nactual:\n%s", c_expected, d.GetAdded())
	}

	d_expected := MakeTagSet(
		MustTag("task-todo"),
	)

	if !iter.SetEquals[Tag](d_expected, d.GetRemoved()) {
		t.Errorf("expected\n%s\nactual:\n%s", d_expected, d.GetRemoved())
	}
}
