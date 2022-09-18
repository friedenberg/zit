package etikett

import (
	"testing"

	"github.com/friedenberg/zit/src/bravo/test_logz"
)

func assertSetRemovesPrefixes(t *testing.T, ac Set, ex Set, prefix string) {
	t.Helper()

	ac.RemovePrefixes(Etikett{Value: prefix})

	if !ac.Equals(ex) {
		test_logz.Errorf(
      test_logz.T{T: t, Skip: 1},
			"removing prefixes doesn't match:\nexpected: %q\n  actual: %q",
			ex,
			ac,
		)
	}
}
