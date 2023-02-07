package kennung

import (
	"testing"

	"github.com/friedenberg/zit/src/bravo/test_logz"
	"github.com/friedenberg/zit/src/echo/ts"
)

func TestMakeProtoIdSet(t1 *testing.T) {
	t := test_logz.T{T: t1}

	sut := MakeProtoIdSet(
		ProtoId{
			Setter: &Hinweis{},
		},
		ProtoId{
			Setter: &Etikett{},
		},
		ProtoId{
			Setter: &Typ{},
		},
		ProtoId{
			Setter: &ts.Time{},
		},
	)

	eLen := 4

	if sut.Len() != eLen {
		t.Errorf("expected %d but got %d", eLen, sut.Len())
	}

	if !sut.Contains(&Hinweis{}) {
		t.Errorf("expected sut to contain hinweis, but did not")
	}

	eString := "test/wow"
	// var set Set
	var err error

	if _, err = sut.Make(eString); err != nil {
		t.Errorf("expected sut create hinweis, but failed: %s", err)
	}
}
