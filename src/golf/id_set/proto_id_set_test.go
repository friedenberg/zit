package id_set

import (
	"os"
	"testing"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/test_logz"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/hinweis"
	"github.com/friedenberg/zit/src/foxtrot/ts"
)

func TestMain(m *testing.M) {
	errors.SetTesting()
	code := m.Run()
	os.Exit(code)
}

func TestMakeProtoIdSet(t1 *testing.T) {
	t := test_logz.T{T: t1}

	sut := MakeProtoIdSet(
		ProtoId{
			Setter: &hinweis.Hinweis{},
		},
		ProtoId{
			Setter: &kennung.Etikett{},
		},
		ProtoId{
			Setter: &kennung.Typ{},
		},
		ProtoId{
			Setter: &ts.Time{},
		},
	)

	eLen := 4

	if sut.Len() != eLen {
		t.Errorf("expected %d but got %d", eLen, sut.Len())
	}

	if !sut.Contains(&hinweis.Hinweis{}) {
		t.Errorf("expected sut to contain hinweis, but did not")
	}

	eString := "test/wow"
	// var set Set
	var err error

	if _, err = sut.Make(eString); err != nil {
		t.Errorf("expected sut create hinweis, but failed: %s", err)
	}
}
