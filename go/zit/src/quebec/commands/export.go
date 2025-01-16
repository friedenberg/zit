package commands

import (
	"bufio"
	"flag"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/delta/age"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

func init() {
	registerCommand(
		"export",
		&Export{
			CompressionType: immutable_config.CompressionTypeEmpty,
		},
	)
}

type Export struct {
	command_components.LocalWorkingCopyWithQueryGroup

	AgeIdentity     age.Identity
	CompressionType immutable_config.CompressionType
}

func (cmd *Export) SetFlagSet(f *flag.FlagSet) {
	cmd.LocalWorkingCopyWithQueryGroup.SetFlagSet(f)

	f.Var(&cmd.AgeIdentity, "age-identity", "")
	cmd.CompressionType.SetFlagSet(f)
}

func (c Export) DefaultSigil() ids.Sigil {
	return ids.MakeSigil(ids.SigilHistory, ids.SigilHidden)
}

func (c Export) DefaultGenres() ids.Genre {
	return ids.MakeGenre(genres.InventoryList)
}

func (cmd Export) Run(dep command.Dep) {
	localWorkingCopy, queryGroup := cmd.MakeLocalWorkingCopyAndQueryGroup(
		dep,
		query.MakeBuilderOptions(cmd),
	)

	var list *sku.List

	{
		var err error

		if list, err = localWorkingCopy.MakeInventoryList(queryGroup); err != nil {
			localWorkingCopy.CancelWithError(err)
		}
	}

	var ag age.Age

	if err := ag.AddIdentity(cmd.AgeIdentity); err != nil {
		localWorkingCopy.CancelWithErrorAndFormat(err, "age-identity: %q", &cmd.AgeIdentity)
	}

	var wc io.WriteCloser

	o := dir_layout.WriteOptions{
		Config: dir_layout.MakeConfig(
			&ag,
			cmd.CompressionType,
			false,
		),
		Writer: localWorkingCopy.GetUIFile(),
	}

	{
		var err error

		if wc, err = dir_layout.NewWriter(o); err != nil {
			localWorkingCopy.CancelWithError(err)
		}
	}

	defer localWorkingCopy.MustClose(wc)

	bw := bufio.NewWriter(wc)
	defer localWorkingCopy.MustFlush(bw)

	printer := localWorkingCopy.MakePrinterBoxArchive(bw, localWorkingCopy.GetConfig().PrintOptions.PrintTime)

	var sk *sku.Transacted
	var hasMore bool

	for {
		localWorkingCopy.ContinueOrPanicOnDone()

		sk, hasMore = list.Pop()

		if !hasMore {
			break
		}

		if err := printer(sk); err != nil {
			localWorkingCopy.CancelWithError(err)
		}
	}
}
