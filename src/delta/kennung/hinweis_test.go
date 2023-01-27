package kennung

import (
	"testing"

	"github.com/friedenberg/zit/src/bravo/test_logz"
)

func TestMake(t1 *testing.T) {
	t := test_logz.T{T: t1}
	in := "ceroplastes/midtown"
	var sut Hinweis

	if err := sut.Set(in); err != nil {
		t.Errorf("expected no error but got: '%s'", err)
	}

	ex := in
	ac := sut.String()

	if ex != ac {
		t.Errorf("expected %q but got %q", ex, ac)
	}
}

func TestMakeKopfUndScwhanz(t1 *testing.T) {
	t := test_logz.T{T: t1}
	k := "ceroplastes"
	s := "midtown"

	var sut Hinweis
	var err error

	if sut, err = MakeHinweisKopfUndSchwanz(k, s); err != nil {
		t.Errorf("expected no error but got: '%s'", err)
	}

	ex := k + "/" + s
	ac := sut.String()

	if ex != ac {
		t.Errorf("expected %q but got %q", ex, ac)
	}
}
