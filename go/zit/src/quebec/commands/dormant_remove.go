package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
)

type DormantRemove struct{}

func init() {
	registerCommandOld(
		"dormant-remove",
		func(f *flag.FlagSet) WithLocalWorkingCopy {
			c := &DormantRemove{}

			return c
		},
	)
}

func (c DormantRemove) Run(u *local_working_copy.Repo, args ...string) {
	u.Must(u.Lock)

	for _, v := range args {
		cs := catgut.MakeFromString(v)

		if err := u.GetDormantIndex().RemoveDormantTag(cs); err != nil {
			u.CancelWithError(err)
		}
	}

	u.Must(u.Unlock)
}
