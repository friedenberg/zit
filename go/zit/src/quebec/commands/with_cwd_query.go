package commands

import (
	"os"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/delta/gattungen"
	"code.linenisgreat.com/zit/src/india/matcher"
	"code.linenisgreat.com/zit/src/india/sku_fmt"
	"code.linenisgreat.com/zit/src/kilo/cwd"
	"code.linenisgreat.com/zit/src/oscar/umwelt"
)

type CommandWithCwdQuery interface {
	RunWithCwdQuery(
		store *umwelt.Umwelt,
		ms matcher.Group,
		cwdFiles *cwd.CwdFiles,
	) error
	DefaultGattungen() gattungen.Set
}

type commandWithCwdQuery struct {
	CommandWithCwdQuery
}

func (c commandWithCwdQuery) Complete(
	u *umwelt.Umwelt,
	args ...string,
) (err error) {
	var cgg CompletionGattungGetter
	ok := false

	if cgg, ok = c.CommandWithCwdQuery.(CompletionGattungGetter); !ok {
		return
	}

	w := sku_fmt.MakeWriterComplete(os.Stdout)
	defer errors.DeferredCloser(&err, w)

	if err = u.StoreObjekten().ReadAllSchwanzen(
		cgg.CompletionGattung(),
		w.WriteOne,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c commandWithCwdQuery) Run(u *umwelt.Umwelt, args ...string) (err error) {
	ids := u.MakeMetaIdSetWithoutExcludedHidden(
		c.DefaultGattungen(),
	)

	if err = ids.SetMany(args...); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.RunWithCwdQuery(u, ids, u.StoreUtil().GetCwdFiles()); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
