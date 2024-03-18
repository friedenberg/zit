package commands

import (
	"os"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/matcher_proto"
	"code.linenisgreat.com/zit/src/india/sku_fmt"
	"code.linenisgreat.com/zit/src/kilo/cwd"
	"code.linenisgreat.com/zit/src/oscar/umwelt"
)

type CommandWithCwdQuery interface {
	RunWithCwdQuery(
		store *umwelt.Umwelt,
		ms matcher_proto.QueryGroup,
		cwdFiles *cwd.CwdFiles,
	) error
	DefaultGattungen() kennung.Gattung
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
	builder := u.MakeMetaIdSetWithoutExcludedHidden(
		c.DefaultGattungen(),
	)

	var ids matcher_proto.QueryGroup

	if ids, err = builder.BuildQueryGroup(args...); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.RunWithCwdQuery(u, ids, u.StoreUtil().GetCwdFiles()); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
