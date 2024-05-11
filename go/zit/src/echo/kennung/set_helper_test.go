package kennung

import (
	"testing"

	"code.linenisgreat.com/zit/src/bravo/test_logz"
)

func assertSetRemovesPrefixes(
	t1 *testing.T,
	ac1 EtikettSet,
	ex EtikettSet,
	prefix string,
) {
	t := &test_logz.T{T: t1}
	t = t.Skip(1)

	ac := ac1.CloneMutableSetPtrLike()
	RemovePrefixes(ac, MustEtikett(prefix))

	if !EtikettSetEquals(ac, ex) {
		t.Errorf(
			"removing prefixes doesn't match:\nexpected: %q\n  actual: %q",
			ex,
			ac,
		)
	}
}
