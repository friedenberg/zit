package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/organize_text"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
	"code.linenisgreat.com/zit/go/zit/src/papa/user_ops"
)

type Clean struct {
	force                    bool
	includeRecognizedBlobs   bool
	includeRecognizedZettels bool
	includeParent            bool
	organize                 bool
}

func init() {
	registerCommandWithQuery(
		"clean",
		func(f *flag.FlagSet) CommandWithQuery {
			c := &Clean{}

			f.BoolVar(
				&c.force,
				"force",
				false,
				"remove objects in working directory even if they have changes",
			)

			f.BoolVar(
				&c.includeParent,
				"include-mutter",
				false,
				"remove objects in working directory if they match their Mutter",
			)

			f.BoolVar(
				&c.includeRecognizedBlobs,
				"recognized-blobs",
				false,
				"remove blobs in working directory or args that are recognized",
			)

			f.BoolVar(
				&c.includeRecognizedZettels,
				"recognized-zettelen",
				false,
				"remove Zetteln in working directory or args that are recognized",
			)

			f.BoolVar(&c.organize, "organize", false, "")

			return c
		},
	)
}

func (c Clean) DefaultGenres() ids.Genre {
	return ids.MakeGenre(genres.TrueGenre()...)
}

func (c Clean) shouldClean(
	u *env.Local,
	co sku.SkuType,
	eqwk *query.Group,
) bool {
	if c.force {
		return true
	}

	state := co.GetState()

	switch state {
	case checked_out_state.CheckedOut:
		return sku.InternalAndExternalEqualsWithoutTai(co)

	case checked_out_state.Recognized:
		return !eqwk.ExcludeRecognized
	}

	if c.includeParent {
		mutter := sku.GetTransactedPool().Get()
		defer sku.GetTransactedPool().Put(mutter)

		err := u.GetStore().GetStreamIndex().ReadOneObjectId(
			co.GetSku().GetObjectId().String(),
			mutter,
		)

		errors.PanicIfError(err)

		if object_metadata.EqualerSansTai.Equals(
			&co.GetSkuExternal().GetSku().Metadata,
			&mutter.Metadata,
		) {
			return true
		}
	}

	return false
}

func (c Clean) ModifyBuilder(b *query.Builder) {
	b.WithHidden(nil)
}

func (c Clean) RunWithQuery(
	u *env.Local,
	qg *query.Group,
) (err error) {
	if c.organize {
		if err = c.runOrganize(u, qg); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	if err = u.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.GetStore().QuerySkuType(
		qg,
		func(co sku.SkuType) (err error) {
			if !c.shouldClean(u, co, qg) {
				return
			}

			if err = u.GetStore().DeleteCheckedOut(co); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.Unlock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c Clean) runOrganize(u *env.Local, qg *query.Group) (err error) {
	opOrganize := user_ops.Organize{
		Local: u,
		Metadata: organize_text.Metadata{
			RepoId: qg.RepoId,
			OptionCommentSet: organize_text.MakeOptionCommentSet(
				nil,
				&organize_text.OptionCommentUnknown{
					Value: "instructions: to clean an object, delete it entirely",
				},
			),
		},
		DontUseQueryGroupForOrganizeMetadata: true,
	}

	ui.Log().Print(qg)

	var organizeResults organize_text.OrganizeResults

	if organizeResults, err = opOrganize.RunWithQueryGroup(
		qg,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	var changes organize_text.Changes

	if changes, err = organize_text.ChangesFromResults(
		u.GetConfig().PrintOptions,
		organizeResults,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = changes.Removed.Each(
		func(el sku.SkuType) (err error) {
			if err = u.GetStore().DeleteCheckedOut(
				el,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.Unlock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
