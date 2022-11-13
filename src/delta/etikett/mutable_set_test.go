package etikett

import "testing"

func TestAddNormalized(t *testing.T) {
	sut := MakeMutableSet(
		Etikett{Value: "project-2021-zit-test"},
		Etikett{Value: "project-2021-zit-ewwwwww"},
		Etikett{Value: "zz-archive-task-done"},
	)

	sutEx := sut.MutableCopy()

	toAdd := Etikett{Value: "project-2021-zit"}

	AddNormalized(sut, toAdd)

	if !sut.Equals(sutEx) {
		t.Errorf("expected %v, but got %v", sutEx, sut)
	}
}

func TestAddNormalizedEmpty(t *testing.T) {
	sut := MakeMutableSet()
	toAdd := Etikett{Value: "project-2021-zit"}

	sutEx := MakeMutableSet(toAdd)

	AddNormalized(sut, toAdd)

	if !sut.Equals(sutEx) {
		t.Errorf("expected %v, but got %v", sutEx, sut)
	}
}

func TestAddNormalizedFromEmptyBuild(t *testing.T) {
	toAdd := []Etikett{
		Etikett{Value: "priority"},
		Etikett{Value: "priority-1"},
	}

	sut := MakeMutableSet()

	sutEx := MakeMutableSet(
		Etikett{Value: "priority-1"},
	)

	for _, e := range toAdd {
		AddNormalized(sut, e)
	}

	if !sut.Equals(sutEx) {
		t.Errorf("expected %v, but got %v", sutEx, sut)
	}
}
