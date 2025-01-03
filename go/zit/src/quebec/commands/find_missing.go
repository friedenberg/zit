package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/november/repo_local"
)

type FindMissing struct{}

func init() {
	registerCommand(
		"find-missing",
		func(f *flag.FlagSet) CommandWithRepo {
			c := &FindMissing{}

			return c
		},
	)
}

func (c FindMissing) RunWithRepo(
	u *repo_local.Repo,
	args ...string,
) {
	var lookupStored map[sha.Bytes][]string

	{
		var err error

		if lookupStored, err = u.GetStore().MakeBlobShaBytesMap(); err != nil {
			u.CancelWithError(err)
		}
	}

	for _, shSt := range args {
		var sh sha.Sha

		if err := sh.Set(shSt); err != nil {
			u.CancelWithError(err)
		}

		oids, ok := lookupStored[sh.GetBytes()]

		if ok {
			u.GetUI().Printf("%s (checked in as %q)", &sh, oids)
		} else {
			u.GetUI().Printf("%s (missing)", &sh)
		}
	}
}
