package ids

import (
	"testing"

	"code.linenisgreat.com/zit/go/zit/src/bravo/test_logz"
)

func assertSetRemovesPrefixes(
	t1 *testing.T,
	ac1 TagSet,
	ex TagSet,
	prefix string,
) {
	t := &test_logz.T{T: t1}
	t = t.Skip(1)

	ac := ac1.CloneMutableSetPtrLike()
	RemovePrefixes(ac, MustTag(prefix))

	if !TagSetEquals(ac, ex) {
		t.Errorf(
			"removing prefixes doesn't match:\nexpected: %q\n  actual: %q",
			ex,
			ac,
		)
	}
}
