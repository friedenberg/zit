package kennung

import "testing"

func TestAddNormalized(t *testing.T) {
	sut := MakeMutableSet(
		MustEtikett("project-2021-zit-test"),
		MustEtikett("project-2021-zit-ewwwwww"),
		MustEtikett("zz-archive-task-done"),
	)

	sutEx := sut.MutableCopy()

	toAdd := MustEtikett("project-2021-zit")

	AddNormalized(sut, toAdd)

	if !sut.Equals(sutEx) {
		t.Errorf("expected %v, but got %v", sutEx, sut)
	}
}

func TestAddNormalizedEmpty(t *testing.T) {
	sut := MakeMutableSet()
	toAdd := MustEtikett("project-2021-zit")

	sutEx := MakeMutableSet(toAdd)

	AddNormalized(sut, toAdd)

	if !sut.Equals(sutEx) {
		t.Errorf("expected %v, but got %v", sutEx, sut)
	}
}

func TestAddNormalizedFromEmptyBuild(t *testing.T) {
	toAdd := []Etikett{
		MustEtikett("priority"),
		MustEtikett("priority-1"),
	}

	sut := MakeMutableSet()

	sutEx := MakeMutableSet(
		MustEtikett("priority-1"),
	)

	for _, e := range toAdd {
		AddNormalized(sut, e)
	}

	if !sut.Equals(sutEx) {
		t.Errorf("expected %v, but got %v", sutEx, sut)
	}
}
