package matcher

import (
	"testing"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/bravo/test_logz"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

func TestMatcherLuaFalse(t1 *testing.T) {
	t := test_logz.T{T: t1}

	m, err := MakeMatcherLua(
		kennung.Index{},
		`function contains_matchable(sku) return false end`,
	)
	if err != nil {
		t.Errorf("expected no error but got %w", err)
	}

	sk := &sku.Transacted{}

	if err = sk.Kennung.SetWithKennung(&kennung.Etikett{}); err != nil {
		err = errors.Wrap(err)
		return
	}

	if m.ContainsMatchable(sk) {
		t.Errorf("woops")
	}
}

func TestMatcherLuaTrue(t1 *testing.T) {
	t := test_logz.T{T: t1}

	m, err := MakeMatcherLua(
		kennung.Index{},
		`function contains_matchable(sku) return true end`,
	)
	if err != nil {
		t.Errorf("expected no error but got %w", err)
	}

	sk := &sku.Transacted{}

	if err = sk.Kennung.SetWithKennung(&kennung.Etikett{}); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !m.ContainsMatchable(sk) {
		t.Errorf("woops")
	}
}
