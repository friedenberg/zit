package matcher

import (
	"testing"

	"github.com/friedenberg/zit/src/bravo/test_logz"
	"github.com/friedenberg/zit/src/echo/kennung"
)

func TestMatcherLua(t1 *testing.T) {
	t1.Skip()
	t := test_logz.T{T: t1}

	m := MakeMatcherWithLua(
		kennung.MustEtikett("base"),
		kennung.MustEtikett("conditional"),
		"",
	)

	if m.ContainsMatchable(nil) {
		t.Errorf("woops")
	}
}
