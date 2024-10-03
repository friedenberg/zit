package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/kilo/organize_text"
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
				"remove Objekten in working directory even if they have changes",
			)

			f.BoolVar(
				&c.includeParent,
				"include-mutter",
				false,
				"remove Objekten in working directory if they match their Mutter",
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
	u *env.Env,
	co sku.CheckedOutLike,
	eqwk *query.Group,
) bool {
	if c.force {
		return true
	}

	state := co.GetState()

	switch state {
	case checked_out_state.ExistsAndSame:
		return true

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
			&co.GetSkuExternalLike().GetSku().Metadata,
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
	u *env.Env,
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

	if err = u.GetStore().QueryCheckedOut(
		qg,
		func(co sku.CheckedOutLike) (err error) {
			if !c.shouldClean(u, co, qg) {
				return
			}

			if err = u.GetStore().DeleteCheckedOutLike(co); err != nil {
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

func (c Clean) runOrganize(u *env.Env, qg *query.Group) (err error) {
	opOrganize := user_ops.Organize{
		Env: u,
		Metadata: organize_text.Metadata{
			RepoId: qg.RepoId,
			OptionCommentSet: organize_text.MakeOptionCommentSet(
				nil,
				organize_text.OptionCommentUnknown(
					"instructions: to clean an object, delete it entirely",
				),
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
		func(el sku.ExternalLike) (err error) {
			if err = u.GetStore().DeleteExternalLike(
				qg.RepoId,
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
