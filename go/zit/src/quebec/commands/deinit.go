package commands

import (
	"flag"
	"fmt"
	"path"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
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

func (c Deinit) Run(u *env.Local, args ...string) (err error) {
	if !c.Force && !c.getPermission(u) {
		return
	}

	base := path.Join(u.GetDirectoryLayout().Dir())

	if err = files.SetAllowUserChangesRecursive(base); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.GetDirectoryLayout().DeleteAll(base); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c Deinit) getPermission(u *env.Local) (success bool) {
	var err error
	ui.Err().Printf(
		"are you sure you want to deinit in %q? (y/*)",
		u.GetDirectoryLayout().Dir(),
	)

	var answer rune
	var n int

	if n, err = fmt.Scanf("%c", &answer); err != nil {
		ui.Err().Printf("failed to read answer: %s", err)
		return
	}

	if n != 1 {
		ui.Err().Printf("failed to read at exactly 1 answer: %s", err)
		return
	}

	if answer == 'y' || answer == 'Y' {
		success = true
	}

	return
}
