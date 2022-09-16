package etikett

import "testing"

func assertSetRemovesPrefixes(t *testing.T, ac Set, ex Set, prefix string) {
	t.Helper()

	ac.RemovePrefixes(Etikett{Value: prefix})

	if !ac.Equals(ex) {
		t.Errorf(
			"removing prefixes doesn't match:\nexpected: %q\n  actual: %q",
			ex,
			ac,
		)
	}
}
