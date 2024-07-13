package query

import (
	"testing"

	"code.linenisgreat.com/zit/go/zit/src/bravo/test_logz"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func TestLuaFalse(t1 *testing.T) {
	t := test_logz.T{T: t1}

	m, err := MakeLua(
		nil,
		`return { contains_sku = function (sku) return false end }`,
		nil,
	)
	if err != nil {
		t.Errorf("expected no error but got %w", err)
	}

	sk := &sku.Transacted{}

	if err = sk.Kennung.SetWithIdLike(&ids.Tag{}); err != nil {
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
		nil,
		`return { contains_sku = function (sku) return true end }`,
		nil,
	)
	if err != nil {
		t.Errorf("expected no error but got %w", err)
	}

	sk := &sku.Transacted{}

	if err = sk.Kennung.SetWithIdLike(&ids.Tag{}); err != nil {
		t.Fatal(err)
	}

	if !m.ContainsSku(sk) {
		t.Errorf("woops")
	}
}
