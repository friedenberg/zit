package kennung

import (
	"testing"

	"github.com/friedenberg/zit/src/bravo/test_logz"
)

func assertSetRemovesPrefixes(t1 *testing.T, ac1 EtikettSet, ex EtikettSet, prefix string) {
	t := test_logz.T{
		T:    t1,
		Skip: 1,
	}

	ac := ac1.MutableClone()
	RemovePrefixes(ac, MustEtikett(prefix))

	if !ac.Equals(ex) {
		t.Errorf(
			"removing prefixes doesn't match:\nexpected: %q\n  actual: %q",
			ex,
			ac,
		)
	}
}
