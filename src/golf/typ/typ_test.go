package typ

import (
	"testing"

	"github.com/friedenberg/zit/src/delta/kennung"
)

func TestAkteExt(t *testing.T) {
	var e kennung.Typ
	var err error

	if err = e.Set(".md"); err != nil {
		t.Fatalf("%s", err)
	}

	actual := e.String()
	expected := "md"

	if expected != actual {
		t.Fatalf("expected %q, but got %q", expected, actual)
	}
}

func TestAkteExt1(t *testing.T) {
	var e kennung.Typ
	var err error

	if err = e.Set("md"); err != nil {
		t.Fatalf("%s", err)
	}

	actual := e.String()
	expected := "md"

	if expected != actual {
		t.Fatalf("expected %q, but got %q", expected, actual)
	}
}