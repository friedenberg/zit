package matcher

import (
	"testing"

	"github.com/friedenberg/zit/src/bravo/test_logz"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/india/transacted"
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

	sk := &transacted.Etikett{}

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

	sk := &transacted.Etikett{}

	if !m.ContainsMatchable(sk) {
		t.Errorf("woops")
	}
}