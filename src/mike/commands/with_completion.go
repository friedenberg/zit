package commands

import (
	"github.com/friedenberg/zit/src/kilo/umwelt"
)

type WithCompletion interface {
	Complete(u *umwelt.Umwelt, args ...string) (err error)
}
