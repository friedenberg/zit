package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/typ"
	"github.com/friedenberg/zit/src/kilo/umwelt"
)

type EditTyp struct {
}

func init() {
	registerCommand(
		"edit-typ",
		func(f *flag.FlagSet) Command {
			c := &EditTyp{}

			return commandWithTypen{c}
		},
	)
}

func (c EditTyp) RunWithTypen(u *umwelt.Umwelt, typen ...typ.Typ) (err error) {
	errors.PrintOut(typen)

	return
}
