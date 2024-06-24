package user_ops

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/vim_cli_options_builder"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/november/umwelt"
)

type Edit struct {
	*umwelt.Umwelt
}

func (u Edit) Run(zsc sku.CheckedOutFSSet) (err error) {
	var filesZettelen []string

	if filesZettelen, err = sku.ToSliceFilesZettelen(zsc); err != nil {
		err = errors.Wrap(err)
		return
	}

	openVimOp := OpenVim{
		Options: vim_cli_options_builder.New().
			WithCursorLocation(2, 3).
			WithFileType("zit-zettel").
			WithInsertMode().
			Build(),
	}

	if _, err = openVimOp.Run(u.Umwelt, filesZettelen...); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.Reset(); err != nil {
		err = errors.Wrap(err)
		return
	}

	checkinOp := Checkin{}

	var ms *query.Group

	builder := u.MakeQueryBuilderExcludingHidden(kennung.MakeGattung(gattung.Zettel))

	col := sku.MakeCheckedOutLikeMutableSet()

	if err = zsc.Each(
		func(cofs *sku.CheckedOutFS) (err error) {
			if err = col.Add(cofs); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if ms, err = builder.WithCheckedOut(col).BuildQueryGroup(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = checkinOp.Run(u.Umwelt, ms); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
