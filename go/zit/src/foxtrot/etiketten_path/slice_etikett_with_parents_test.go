package etiketten_path

import (
	"testing"

	"code.linenisgreat.com/zit/src/bravo/test_logz"
	"code.linenisgreat.com/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/src/echo/kennung"
)

func TestContains(t1 *testing.T) {
	t := test_logz.T{T: t1}

	var s SliceEtikettWithParents

	s.Add(catgut.MakeFromString("%chrome-tab_id-1"), nil)
	s.Add(catgut.MakeFromString("%chrome-tab_id-2"), nil)

	var e kennung.Kennung2
	e.Set("%chrome-tab_id")

	if _, ok := s.ContainsKennungEtikett(&e); !ok {
		t.Errorf("expected %q to be in %s", &e, &s)
	}
}
