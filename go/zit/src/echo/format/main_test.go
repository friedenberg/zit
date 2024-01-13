package format

import (
	"os"
	"testing"

	"github.com/friedenberg/zit/src/alfa/errors"
)

func TestMain(m *testing.M) {
	errors.SetTesting()
	code := m.Run()
	os.Exit(code)
}
