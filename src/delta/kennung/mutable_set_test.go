package kennung

import (
	"testing"

	"github.com/friedenberg/zit/src/bravo/test_logz"
)

func TestAddNormalized(t1 *testing.T) {
	t := test_logz.T{T: t1}

	sut := MakeEtikettMutableSet(
		MustEtikett("project-2021-zit-test"),
		MustEtikett("project-2021-zit-ewwwwww"),
		MustEtikett("zz-archive-task-done"),
	)

	sutEx := sut.ImmutableClone()

	toAdd := MustEtikett("project-2021-zit")

	AddNormalized(sut, toAdd)

	if !sut.EqualsSetLike(sutEx) {
		t.NotEqual(sutEx, sut)
	}
}

func TestAddNormalizedEmpty(t *testing.T) {
	sut := MakeEtikettMutableSet()
	toAdd := MustEtikett("project-2021-zit")

	sutEx := MakeEtikettMutableSet(toAdd)

	AddNormalized(sut, toAdd)

	if !sut.EqualsSetLike(sutEx) {
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

	for _, e := range toAdd {
		AddNormalized(sut, e)
	}

	if !sut.EqualsSetLike(sutEx) {
		t.Errorf("expected %v, but got %v", sutEx, sut)
	}
}
