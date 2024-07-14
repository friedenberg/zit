package ids

import (
	"testing"
)

func TestBlobExt(t *testing.T) {
	var e Type
	var err error

	if err = e.Set("!md"); err != nil {
		t.Fatalf("%s", err)
	}

	actual := e.String()
	expected := "md"

	if expected != actual {
		t.Fatalf("expected %q, but got %q", expected, actual)
	}
}

func TestBlobExt1(t *testing.T) {
	var e Type
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
