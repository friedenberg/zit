package format

import (
	"os"
	"testing"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
)

func TestMain(m *testing.M) {
	errors.SetTesting()
	code := m.Run()
	os.Exit(code)
}
