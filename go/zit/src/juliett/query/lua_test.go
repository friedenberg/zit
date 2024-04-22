package query

import (
	"testing"

	"code.linenisgreat.com/zit/src/bravo/test_logz"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

func TestLuaFalse(t1 *testing.T) {
	t := test_logz.T{T: t1}

	m, err := MakeLua(
		`return { contains_sku = function (sku) return false end }`,
	)
	if err != nil {
		t.Errorf("expected no error but got %w", err)
	}

	sk := &sku.Transacted{}

	if err = sk.Kennung.SetWithKennung(&kennung.Etikett{}); err != nil {
		t.Fatal(err)
		return
	}

	if m.ContainsSku(sk) {
		t.Errorf("woops")
	}
}

func TestMatcherLuaTrue(t1 *testing.T) {
	t := test_logz.T{T: t1}

	m, err := MakeLua(
		`return { contains_sku = function (sku) return true end }`,
	)
	if err != nil {
		t.Errorf("expected no error but got %w", err)
	}

	sk := &sku.Transacted{}

	if err = sk.Kennung.SetWithKennung(&kennung.Etikett{}); err != nil {
		t.Fatal(err)
	}

	if !m.ContainsSku(sk) {
		t.Errorf("woops")
	}
}
