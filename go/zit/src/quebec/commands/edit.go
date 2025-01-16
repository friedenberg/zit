package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
	"code.linenisgreat.com/zit/go/zit/src/papa/user_ops"
)

func init() {
	command.Register(
		"edit",
		&Edit{
			Workspace:    true,
			CheckoutMode: checkout_mode.MetadataOnly,
		},
	)
}

type Edit struct {
	command_components.LocalWorkingCopyWithQueryGroup

	// TODO-P3 add force
	Workspace bool
	command_components.Checkout
	CheckoutMode checkout_mode.Mode
}

func (cmd *Edit) SetFlagSet(f *flag.FlagSet) {
	cmd.LocalWorkingCopyWithQueryGroup.SetFlagSet(f)

	cmd.Checkout.SetFlagSet(f)

	f.Var(&cmd.CheckoutMode, "mode", "mode for checking out the object")
	f.BoolVar(&cmd.Workspace, "use-workspace", true, "checkout the object into the current workspace (CWD)")
}

func (c Edit) CompletionGenres() ids.Genre {
	return ids.MakeGenre(
		genres.Tag,
		genres.Zettel,
		genres.Type,
		genres.Repo,
	)
}

func (c Edit) DefaultGenres() ids.Genre {
	return ids.MakeGenre(
		genres.Tag,
		genres.Zettel,
		genres.Type,
		genres.Repo,
	)
}

func (cmd Edit) Run(dep command.Dep) {
	localWorkingCopy, queryGroup := cmd.MakeLocalWorkingCopyAndQueryGroup(
		dep,
		query.MakeBuilderOptions(cmd),
	)

	options := checkout_options.Options{
		CheckoutMode: cmd.CheckoutMode,
	}

	opEdit := user_ops.Checkout{
		Repo:    localWorkingCopy,
		Options: options,
		Edit:    true,
	}

	opEdit.Options.Workspace = cmd.Workspace

	if _, err := opEdit.RunQuery(queryGroup); err != nil {
		localWorkingCopy.CancelWithError(err)
	}
}
