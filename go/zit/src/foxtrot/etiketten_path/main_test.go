package etiketten_path

import (
	"bytes"
	"testing"

	"code.linenisgreat.com/zit/src/bravo/test_logz"
	"code.linenisgreat.com/zit/src/delta/catgut"
)

func TestReadWrite(t1 *testing.T) {
	t := test_logz.T{T: t1}

	b := new(bytes.Buffer)
	var sut Path

	sut.Add(catgut.MakeFromString("one"))
	sut.Add(catgut.MakeFromString("two"))
	sut.Add(catgut.MakeFromString("three"))

	{
		n, err := sut.WriteTo(b)
		t.AssertNoError(err)
		if int(n) != b.Len() {
			t.NotEqual(b.Len(), n)
		}
	}

	b.Reset()

	{
		n, err := sut.ReadFrom(b)
		t.AssertEOF(err)

		if int(n) != b.Len() {
			t.NotEqual(b.Len(), n)
		}

		if sut.Len() != 3 {
			t.NotEqual(3, sut.Len())
		}

		if !sut[0].EqualsString("one") {
			t.NotEqual("one", sut[0])
		}

		if !sut[1].EqualsString("two") {
			t.NotEqual("two", sut[1])
		}

		if !sut[2].EqualsString("three") {
			t.NotEqual("three", sut[2])
		}
	}
}
