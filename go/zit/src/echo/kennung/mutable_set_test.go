package kennung

import (
	"testing"

	"code.linenisgreat.com/zit/go/zit/src/bravo/test_logz"
)

func TestAddNormalized(t1 *testing.T) {
	t := test_logz.T{T: t1}

	sut := MakeTagMutableSet(
		MustTag("project-2021-zit-test"),
		MustTag("project-2021-zit-ewwwwww"),
		MustTag("zz-archive-task-done"),
	)

	sutEx := sut.CloneSetPtrLike()

	toAdd := MustTag("project-2021-zit")

	AddNormalizedEtikett(sut, &toAdd)

	if !TagSetEquals(sut, sutEx) {
		t.NotEqual(sutEx, sut)
	}
}

func TestAddNormalizedEmpty(t *testing.T) {
	sut := MakeTagMutableSet()
	toAdd := MustTag("project-2021-zit")

	sutEx := MakeTagMutableSet(toAdd)

	AddNormalizedEtikett(sut, &toAdd)

	if !TagSetEquals(sut, sutEx) {
		t.Errorf("expected %v, but got %v", sutEx, sut)
	}
}

func TestAddNormalizedFromEmptyBuild(t *testing.T) {
	toAdd := []Tag{
		MustTag("priority"),
		MustTag("priority-1"),
	}

	sut := MakeTagMutableSet()

	sutEx := MakeTagMutableSet(
		MustTag("priority-1"),
	)

	for i := range toAdd {
		e := &toAdd[i]
		AddNormalizedEtikett(sut, e)
	}

	if !TagSetEquals(sut, sutEx) {
		t.Errorf("expected %v, but got %v", sutEx, sut)
	}
}
