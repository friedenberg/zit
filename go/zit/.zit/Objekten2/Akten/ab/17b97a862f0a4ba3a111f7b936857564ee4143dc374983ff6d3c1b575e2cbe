package user_ops

import (
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/organize_text"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type Checkout struct {
	*env.Local
	Organize bool
	checkout_options.Options
	Open    bool
	Edit    bool
	Utility string
}

func (op Checkout) Run(
	skus sku.TransactedSet,
) (zsc sku.SkuTypeSetMutable, err error) {
	var k ids.RepoId

	if zsc, err = op.RunWithKasten(k, skus); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (op Checkout) RunWithKasten(
	kasten ids.RepoId,
	skus sku.TransactedSet,
) (zsc sku.SkuTypeSetMutable, err error) {
	b := op.Local.MakeQueryBuilder(
		ids.MakeGenre(genres.Zettel),
	).WithTransacted(
		skus,
		ids.SigilExternal,
	).WithRequireNonEmptyQuery()

	var qg *query.Group

	if qg, err = b.BuildQueryGroup(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if zsc, err = op.RunQuery(
		qg,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (op Checkout) RunQuery(
	qg *query.Group,
) (zsc sku.SkuTypeSetMutable, err error) {
	zsc = sku.MakeSkuTypeSetMutable()
	var l sync.Mutex

	onCheckedOut := func(col sku.SkuType) (err error) {
		l.Lock()
		defer l.Unlock()

		cl := col.Clone()

		if err = zsc.Add(cl); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	if op.Organize {
		if qg, err = op.runOrganize(qg, onCheckedOut); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = op.Local.GetStore().CheckoutQuery(
		op.Options,
		qg,
		onCheckedOut,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if op.Utility != "" {
		eachBlobOp := EachBlob{
			Utility: op.Utility,
			Local:     op.Local,
		}

		if err = eachBlobOp.Run(zsc); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if op.Open || op.Edit {
		if err = op.GetStore().Open(
			qg.RepoId,
			op.CheckoutMode,
			op.PrinterHeader(),
			zsc,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if op.Edit {
		if err = op.Reset(); err != nil {
			err = errors.Wrap(err)
			return
		}

		var qg *query.Group

		builder := op.MakeQueryBuilderExcludingHidden(ids.MakeGenre(genres.Zettel))

		if qg, err = builder.WithCheckedOut(
			zsc,
		).BuildQueryGroup(); err != nil {
			err = errors.Wrap(err)
			return
		}

		checkinOp := Checkin{}

		if err = checkinOp.Run(
			op.Local,
			qg,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (op Checkout) runOrganize(
	qgOriginal *query.Group,
	onCheckedOut interfaces.FuncIter[sku.SkuType],
) (qgModified *query.Group, err error) {
	opOrganize := Organize{
		Local: op.Local,
		Metadata: organize_text.Metadata{
			RepoId: qgOriginal.RepoId,
			OptionCommentSet: organize_text.MakeOptionCommentSet(
				// TODO add other OptionComments
				nil,
				&organize_text.OptionCommentUnknown{
					Value: "instructions: to prevent an object from being checked out, delete it entirely",
				},
			),
		},
		DontUseQueryGroupForOrganizeMetadata: true,
	}

	ui.Log().Print(qgOriginal)

	originalRepoId := qgOriginal.RepoId
	qgOriginal.RepoId.Reset()

	var organizeResults organize_text.OrganizeResults

	if organizeResults, err = opOrganize.RunWithQueryGroup(
		qgOriginal,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	var changeResults organize_text.Changes

	if changeResults, err = organize_text.ChangesFromResults(
		op.GetConfig().PrintOptions,
		organizeResults,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	b := op.MakeQueryBuilder(
		ids.MakeGenre(genres.TrueGenre()...),
	).WithTransacted(
		changeResults.After.AsTransactedSet(),
		ids.SigilExternal,
	).WithDoNotMatchEmpty().WithRequireNonEmptyQuery()

	if qgModified, err = b.BuildQueryGroup(); err != nil {
		err = errors.Wrap(err)
		return
	}

	qgModified.RepoId = originalRepoId

	return
}
