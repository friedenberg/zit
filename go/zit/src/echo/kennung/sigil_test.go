package kennung

import (
	"bytes"
	"testing"

	"code.linenisgreat.com/zit/src/bravo/test_logz"
)

func TestSigilReadWrite(t1 *testing.T) {
	t := test_logz.T{T: t1}

	sut := SigilAll
	b := bytes.NewBuffer(nil)

	{
		n, err := sut.WriteTo(b)
		t.AssertNoError(err)
		if n != 1 {
			t.NotEqual(1, n)
		}
	}

	var actual Sigil

	{
		n, err := actual.ReadFrom(b)
		t.AssertNoError(err)
		if n != 1 {
			t.NotEqual(1, n)
		}
	}

	if actual != sut {
		t.NotEqual(sut, actual)
	}
}
