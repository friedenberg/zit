package collections

import (
	"os"
	"testing"

	"code.linenisgreat.com/zit/src/alfa/errors"
)

func TestMain(m *testing.M) {
	errors.SetTesting()
	code := m.Run()
	os.Exit(code)
}
