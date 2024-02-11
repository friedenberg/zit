package kennung

import (
	"bytes"
	"testing"

	"code.linenisgreat.com/zit/src/bravo/test_logz"
)

func kennung2WriteToReadFromData() []string {
	return []string{
		"one/uno",
		"konfig",
		"!md",
		"-etikett",
		"//kasten",
		"etikett",
	}
}

func TestKennung2WriteToReadFrom(t1 *testing.T) {
	t := test_logz.T{T: t1}
	for _, v := range kennung2WriteToReadFromData() {
		var k Kennung2
		t.AssertNoError(k.Set(v))

		var b bytes.Buffer

		_, err := k.WriteTo(&b)
		t.AssertNoError(err)

		var k2 Kennung2

		_, err = k2.ReadFrom(&b)
		t.AssertNoError(err)

		if k.String() != k2.String() {
			t.NotEqual(&k, &k2)
		}
	}
}
