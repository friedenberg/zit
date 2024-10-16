package commands

import (
	"flag"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
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
	u *env.Env,
	args ...string,
) (err error) {
	lookupStored := make(map[sha.Bytes]struct{}, len(args))
	var l sync.Mutex

	if err = u.GetStore().QueryPrimitive(
		sku.MakePrimitiveQueryGroup(),
		func(sk *sku.Transacted) (err error) {
			l.Lock()
			defer l.Unlock()

			lookupStored[sk.Metadata.Blob.GetBytes()] = struct{}{}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, shSt := range args {
		var sh sha.Sha

		if err = sh.Set(shSt); err != nil {
			err = errors.Wrap(err)
			return
		}

		_, ok := lookupStored[sh.GetBytes()]

		if ok {
			ui.Out().Printf("%s (checked in)", &sh)
		} else {
			ui.Out().Printf("%s (missing)", &sh)
		}
	}

	return
}
