package kennung

import (
	"bytes"
	"encoding/gob"
	"testing"

	"code.linenisgreat.com/zit/go/zit/src/bravo/test_logz"
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

	var sut *Hinweis
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

func TestGob(t1 *testing.T) {
	t := test_logz.T{T: t1}
	k := "ceroplastes"
	s := "midtown"

	var sut *Hinweis
	var err error

	if sut, err = MakeHinweisKopfUndSchwanz(k, s); err != nil {
		t.Errorf("expected no error but got: '%s'", err)
	}

	b := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(b)

	err = enc.Encode(sut)

	if err != nil {
		t.Errorf("expected no error but got: '%s'", err)
	}

	var sut2 Hinweis

	dec := gob.NewDecoder(b)

	if err = dec.Decode(&sut2); err != nil {
		t.Errorf("expected no error but got: '%s'", err)
	}

	ex := k + "/" + s
	ac := sut2.String()

	if ac != ex {
		t.Errorf("expected %q but got %q", ex, ac)
	}
}
