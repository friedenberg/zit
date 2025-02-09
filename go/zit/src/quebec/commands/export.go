package commands

import (
	"bufio"
	"flag"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/delta/age"
	"code.linenisgreat.com/zit/go/zit/src/delta/config_immutable"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/env_dir"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

func init() {
	command.Register(
		"export",
		&Export{
			CompressionType: config_immutable.CompressionTypeEmpty,
		},
	)
}

type Export struct {
	command_components.LocalWorkingCopyWithQueryGroup

	AgeIdentity     age.Identity
	CompressionType config_immutable.CompressionType
}

func (cmd *Export) SetFlagSet(f *flag.FlagSet) {
	cmd.LocalWorkingCopyWithQueryGroup.SetFlagSet(f)

	f.Var(&cmd.AgeIdentity, "age-identity", "")
	cmd.CompressionType.SetFlagSet(f)
}

func (cmd Export) Run(dep command.Request) {
	localWorkingCopy, queryGroup := cmd.MakeLocalWorkingCopyAndQueryGroup(
		dep,
		query.BuilderOptionsOld(
			cmd,
			query.BuilderOptionDefaultSigil(
				ids.SigilHistory,
				ids.SigilHidden,
			),
			query.BuilderOptionDefaultGenres(
				genres.InventoryList,
			),
		),
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

	o := env_dir.WriteOptions{
		Config: env_dir.MakeConfig(
			&cmd.CompressionType,
			&ag,
			false,
		),
		Writer: localWorkingCopy.GetUIFile(),
	}

	{
		var err error

		if wc, err = env_dir.NewWriter(o); err != nil {
			localWorkingCopy.CancelWithError(err)
		}
	}

	defer localWorkingCopy.MustClose(wc)

	bw := bufio.NewWriter(wc)
	defer localWorkingCopy.MustFlush(bw)

	printer := localWorkingCopy.MakePrinterBoxArchive(bw, localWorkingCopy.GetConfig().GetCLIConfig().PrintOptions.PrintTime)

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
