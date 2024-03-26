package kennung

import (
	"bytes"
	"testing"

	"code.linenisgreat.com/zit/src/bravo/test_logz"
)

func TestSigilContains(t1 *testing.T) {
	t := test_logz.T{T: t1}

	sut := SigilAll

	if !sut.ContainsOneOf(SigilSchwanzen) {
		t.Errorf("expected SigilAll to contain SigilSchwanzen")
	}

	sut = SigilSchwanzen
	sut.Add(SigilHidden)

	if !sut.ContainsOneOf(SigilSchwanzen) {
		t.Errorf("expected sut to contain SigilSchwanzen")
	}

	if !sut.ContainsOneOf(SigilHidden) {
		t.Errorf("expected sut to contain SigilHidden")
	}

	if !sut.ContainsOneOf(sut) {
		t.Errorf("expected sut to contain sut")
	}

	if !sut.Contains(sut) {
		t.Errorf("expected sut to contain sut")
	}

	other := SigilHistory
	other.Add(SigilHidden)

	if !sut.ContainsOneOf(other) {
		t.Errorf("expected sut to contain one ofother")
	}

	if sut.Contains(other) {
		t.Errorf("expected sut not to contain other")
	}

	if sut.ContainsOneOf(SigilCwd) {
		t.Errorf("expected sut not to contain SigilCwd")
	}
}

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
