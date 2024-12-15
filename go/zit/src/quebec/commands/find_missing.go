package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type FindMissing struct{}

func init() {
	registerCommand(
		"find-missing",
		func(f *flag.FlagSet) Command {
			c := &FindMissing{}

			return c
		},
	)
}

func (c FindMissing) Run(
	u *env.Local,
	args ...string,
) (err error) {
	var lookupStored map[sha.Bytes][]string

	if lookupStored, err = u.GetStore().MakeBlobShaBytesMap(); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, shSt := range args {
		var sh sha.Sha

		if err = sh.Set(shSt); err != nil {
			err = errors.Wrap(err)
			return
		}

		oids, ok := lookupStored[sh.GetBytes()]

		if ok {
			ui.Out().Printf("%s (checked in as %q)", &sh, oids)
		} else {
			ui.Out().Printf("%s (missing)", &sh)
		}
	}

	return
}
