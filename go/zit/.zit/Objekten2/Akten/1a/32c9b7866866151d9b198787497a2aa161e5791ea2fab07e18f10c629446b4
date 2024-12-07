package catgut

import (
	"io"
	"strings"
	"testing"

	"code.linenisgreat.com/zit/go/zit/src/bravo/test_logz"
)

func TestSliceReader(t1 *testing.T) {
	t := test_logz.T{T: t1}
	input := Slice{
		data: [2][]byte{
			[]byte("test"),
			[]byte("string"),
		},
	}

	sut := MakeSliceReader(input)

	var actual strings.Builder

	n1, err := io.Copy(&actual, sut)
	t.AssertNoError(err)
	n := int(n1)

	if n != input.Len() {
		t.NotEqual(input.Len(), n)
	}
}
