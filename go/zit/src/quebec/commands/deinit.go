package commands

import (
	"flag"
	"fmt"
	"os"
	"path"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/bravo/files"
	"code.linenisgreat.com/zit/src/oscar/umwelt"
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
	if !c.Force && !c.getPermission(u) {
		return
	}

	base := path.Join(u.Standort().Dir(), ".zit")

	if err = files.SetAllowUserChangesRecursive(base); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = os.RemoveAll(base); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c Deinit) getPermission(u *umwelt.Umwelt) (success bool) {
	var err error
	errors.Err().Printf(
		"are you sure you want to deinit in %q? (y/*)",
		u.Standort().Dir(),
	)

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
