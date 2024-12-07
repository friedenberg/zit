package files

import (
	"reflect"
	"testing"

	"code.linenisgreat.com/zit/go/zit/src/bravo/test_logz"
)

func TestPathElements(t1 *testing.T) {
	t := test_logz.T{T: t1}

	path := "/wow/ok/great.ext"
	expected := []string{"ext", "great", "ok", "wow"}
	actual := PathElements(path)

	if reflect.DeepEqual(expected, actual) {
		t.AssertEqual(expected, actual)
	}
}
