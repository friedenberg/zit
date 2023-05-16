package commands

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type Deinit struct {
	Force bool
}

func init() {
	registerCommand(
		"deinit",
		func(f *flag.FlagSet) Command {
			c := &Deinit{}

			f.BoolVar(
				&c.Force,
				"force",
				false,
				"force deinit",
			)

			return c
		},
	)
}

func (c Deinit) Run(u *umwelt.Umwelt, args ...string) (err error) {
	if !c.Force && !c.getPermission() {
		return
	}

	base := path.Join(u.Standort().Dir(), ".zit")
	err = os.RemoveAll(base)

	if err != nil {
		return
	}

	return
}

func (c Deinit) getPermission() (success bool) {
	var err error
	errors.Err().Printf("are you sure you want to deinit? (y/*)")

	var answer rune
	var n int

	if n, err = fmt.Scanf("%c", &answer); err != nil {
		errors.Err().Printf("failed to read answer: %s", err)
		return
	}

	if n != 1 {
		errors.Err().Printf("failed to read at exactly 1 answer: %s", err)
		return
	}

	if answer == 'y' || answer == 'Y' {
		success = true
	}

	return
}
