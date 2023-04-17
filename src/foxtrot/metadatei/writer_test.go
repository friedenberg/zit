package metadatei

import (
	"strings"
	"testing"

	"github.com/friedenberg/zit/src/bravo/test_logz"
)

func TestWriter1(t1 *testing.T) {
	t := test_logz.T{
		T: t1,
	}

	expectedOut := `---
metadatei
---

akte
`

	out := &strings.Builder{}

	sut := Writer{
		Metadatei: strings.NewReader("metadatei\n"),
		Akte:      strings.NewReader("akte\n"),
	}

	sut.WriteTo(out)

	if out.String() != expectedOut {
		t.Errorf("expected %q but got %q", expectedOut, out.String())
	}
}

func TestWriter2(t1 *testing.T) {
	t := test_logz.T{
		T: t1,
	}

	expectedOut := `---
metadatei
---
`

	out := &strings.Builder{}

	sut := Writer{
		Metadatei: strings.NewReader("metadatei\n"),
	}

	sut.WriteTo(out)

	if out.String() != expectedOut {
		t.Errorf("expected %q but got %q", expectedOut, out.String())
	}
}

func TestWriter3(t1 *testing.T) {
	t := test_logz.T{
		T: t1,
	}

	expectedOut := `akte
`

	out := &strings.Builder{}

	sut := Writer{
		Akte: strings.NewReader("akte\n"),
	}

	sut.WriteTo(out)

	if out.String() != expectedOut {
		t.Errorf("expected %q but got %q", expectedOut, out.String())
	}
}
