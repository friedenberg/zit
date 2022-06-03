package etikett

import "testing"

func TestExpandedRight(t *testing.T) {
	s := NewSet(
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
