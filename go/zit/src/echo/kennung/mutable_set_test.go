package kennung

import (
	"testing"

	"code.linenisgreat.com/zit/go/zit/src/bravo/test_logz"
)

func TestAddNormalized(t1 *testing.T) {
	t := test_logz.T{T: t1}

	sut := MakeEtikettMutableSet(
		MustEtikett("project-2021-zit-test"),
		MustEtikett("project-2021-zit-ewwwwww"),
		MustEtikett("zz-archive-task-done"),
	)

	sutEx := sut.CloneSetPtrLike()

	toAdd := MustEtikett("project-2021-zit")

	AddNormalizedEtikett(sut, &toAdd)

	if !EtikettSetEquals(sut, sutEx) {
		t.NotEqual(sutEx, sut)
	}
}

func TestAddNormalizedEmpty(t *testing.T) {
	sut := MakeEtikettMutableSet()
	toAdd := MustEtikett("project-2021-zit")

	sutEx := MakeEtikettMutableSet(toAdd)

	AddNormalizedEtikett(sut, &toAdd)

	if !EtikettSetEquals(sut, sutEx) {
		t.Errorf("expected %v, but got %v", sutEx, sut)
	}
}

func TestAddNormalizedFromEmptyBuild(t *testing.T) {
	toAdd := []Etikett{
		MustEtikett("priority"),
		MustEtikett("priority-1"),
	}

	sut := MakeEtikettMutableSet()

	sutEx := MakeEtikettMutableSet(
		MustEtikett("priority-1"),
	)

	for i := range toAdd {
		e := &toAdd[i]
		AddNormalizedEtikett(sut, e)
	}

	if !EtikettSetEquals(sut, sutEx) {
		t.Errorf("expected %v, but got %v", sutEx, sut)
	}
}
