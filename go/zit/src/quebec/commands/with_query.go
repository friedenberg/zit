package commands

import (
	"os"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/delta/gattungen"
	"code.linenisgreat.com/zit/src/india/matcher"
	"code.linenisgreat.com/zit/src/india/sku_fmt"
	"code.linenisgreat.com/zit/src/oscar/umwelt"
)

type CommandWithQuery interface {
	RunWithQuery(store *umwelt.Umwelt, ids matcher.Query) error
	DefaultGattungen() gattungen.Set
}

type commandWithQuery struct {
	CommandWithQuery
}

type CompletionGattungGetter interface {
	CompletionGattung() gattungen.Set
}

func (c commandWithQuery) Complete(
	u *umwelt.Umwelt,
	args ...string,
) (err error) {
	var cgg CompletionGattungGetter
	ok := false

	if cgg, ok = c.CommandWithQuery.(CompletionGattungGetter); !ok {
		return
	}

	cg := cgg.CompletionGattung()

	zw := sku_fmt.MakeWriterComplete(os.Stdout)
	defer errors.DeferredCloser(&err, &zw)

	w := zw.WriteZettelVerzeichnisse

	if err = u.StoreObjekten().ReadAllSchwanzen(
		cg,
		w,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c commandWithQuery) Run(u *umwelt.Umwelt, args ...string) (err error) {
	ids := u.MakeMetaIdSetWithExcludedHidden(c.DefaultGattungen())

	if err = ids.SetMany(args...); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.RunWithQuery(u, ids); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
