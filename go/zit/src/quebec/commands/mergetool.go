package commands

import (
	"bufio"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/builtin_types"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

func init() {
	command.Register("merge-tool", &Mergetool{})
}

type Mergetool struct {
	command_components.LocalWorkingCopyWithQueryGroup
}

func (cmd Mergetool) Run(req command.Request) {
	localWorkingCopy, queryGroup := cmd.MakeLocalWorkingCopyAndQueryGroup(
		req,
		query.BuilderOptionsOld(
			cmd,
			query.BuilderOptionDefaultGenres(genres.All()...),
		),
	)

	envWorkspace := localWorkingCopy.GetEnvWorkspace()
	envWorkspace.AssertNotTemporary(req)

	conflicted := sku.MakeSkuTypeSetMutable()

	if err := localWorkingCopy.GetStore().QuerySkuType(
		queryGroup,
		func(co sku.SkuType) (err error) {
			if co.GetState() != checked_out_state.Conflicted {
				return
			}

			if err = conflicted.Add(co.Clone()); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		localWorkingCopy.CancelWithError(err)
	}

	localWorkingCopy.Must(localWorkingCopy.Lock)

	if conflicted.Len() == 0 {
		// TODO-P2 return status 1 and use Err
		localWorkingCopy.GetUI().Printf("nothing to merge")
		return
	}

	for co := range conflicted.All() {
		cmd.doOne(localWorkingCopy, co)
	}

	localWorkingCopy.Must(localWorkingCopy.Unlock)
}

func (c Mergetool) doOne(u *local_working_copy.Repo, co *sku.CheckedOut) {
	tm := sku.Conflicted{
		CheckedOut: co.Clone(),
	}

	var conflict *fd.FD

	{
		var err error

		if conflict, err = u.GetStore().GetStoreFS().GetConflictOrError(
			co.GetSkuExternal(),
		); err != nil {
			u.CancelWithError(err)
		}
	}

	var f *os.File

	{
		var err error

		if f, err = files.Open(conflict.GetPath()); err != nil {
			u.CancelWithError(err)
		}

		defer u.MustClose(f)
	}

	br := bufio.NewReader(f)

	bs := u.GetStore().GetTypedBlobStore().InventoryList

	if err := tm.ReadConflictMarker(
		bs.IterInventoryListBlobSkusFromReader(
			builtin_types.DefaultOrPanic(genres.InventoryList),
			br,
		),
	); err != nil {
		u.CancelWithError(err)
	}

	if err := u.GetStore().RunMergeTool(
		tm,
	); err != nil {
		u.CancelWithError(err)
	}
}
