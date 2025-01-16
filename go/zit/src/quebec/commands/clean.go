package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/organize_text"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
	"code.linenisgreat.com/zit/go/zit/src/papa/user_ops"
)

func init() {
	command.Register("clean", &Clean{})
}

type Clean struct {
	command_components.LocalWorkingCopyWithQueryGroup

	force                    bool
	includeRecognizedBlobs   bool
	includeRecognizedZettels bool
	includeParent            bool
	organize                 bool
}

func (c *Clean) SetFlagSet(f *flag.FlagSet) {
	c.LocalWorkingCopyWithQueryGroup.SetFlagSet(f)

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
}

func (c Clean) DefaultGenres() ids.Genre {
	return ids.MakeGenre(genres.TrueGenre()...)
}

func (c Clean) ModifyBuilder(b *query.Builder) {
	b.WithHidden(nil)
}

func (cmd Clean) Run(dep command.Dep) {
	localWorkingCopy, queryGroup := cmd.MakeLocalWorkingCopyAndQueryGroup(
		dep,
		query.MakeBuilderOptions(cmd),
	)

	if cmd.organize {
		if err := cmd.runOrganize(localWorkingCopy, queryGroup); err != nil {
			localWorkingCopy.CancelWithError(err)
		}

		return
	}

	localWorkingCopy.Must(localWorkingCopy.Lock)

	if err := localWorkingCopy.GetStore().QuerySkuType(
		queryGroup,
		func(co sku.SkuType) (err error) {
			if !cmd.shouldClean(localWorkingCopy, co, queryGroup) {
				return
			}

			if err = localWorkingCopy.GetStore().DeleteCheckedOut(co); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		localWorkingCopy.CancelWithError(err)
	}

	localWorkingCopy.Must(localWorkingCopy.Unlock)
}

func (c Clean) runOrganize(u *local_working_copy.Repo, qg *query.Group) (err error) {
	opOrganize := user_ops.Organize{
		Repo: u,
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

	u.Must(u.Lock)

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

	u.Must(u.Unlock)

	return
}

func (c Clean) shouldClean(
	u *local_working_copy.Repo,
	co sku.SkuType,
	qg *query.Group,
) bool {
	if c.force {
		return true
	}

	state := co.GetState()

	switch state {
	case checked_out_state.CheckedOut:
		return sku.InternalAndExternalEqualsWithoutTai(co)

	case checked_out_state.Recognized:
		return !qg.ExcludeRecognized
	}

	if c.includeParent {
		mutter := sku.GetTransactedPool().Get()
		defer sku.GetTransactedPool().Put(mutter)

		err := u.GetStore().GetStreamIndex().ReadOneObjectId(
			co.GetSku().GetObjectId(),
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
