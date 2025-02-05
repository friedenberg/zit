package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
	"code.linenisgreat.com/zit/go/zit/src/papa/user_ops"
)

func init() {
	cmd := &Checkin{}
	command.Register("checkin", cmd)
	command.Register("add", cmd)
	command.Register("save", cmd)
}

type Checkin struct {
	command_components.LocalWorkingCopyWithQueryGroup

	IgnoreBlob bool
	Proto      sku.Proto

	command_components.Checkout

	CheckoutBlobAndRun string
	OpenBlob           bool
}

func (cmd *Checkin) SetFlagSet(f *flag.FlagSet) {
	cmd.LocalWorkingCopyWithQueryGroup.SetFlagSet(f)

	f.BoolVar(
		&cmd.IgnoreBlob,
		"ignore-blob",
		false,
		"do not change the blob",
	)

	f.StringVar(
		&cmd.CheckoutBlobAndRun,
		"each-blob",
		"",
		"checkout each Blob and run a utility",
	)

	cmd.Proto.SetFlagSet(f)
	cmd.Checkout.SetFlagSet(f)
}

func (c Checkin) DefaultGenres() ids.Genre {
	return ids.MakeGenre(genres.TrueGenre()...)
}

func (c *Checkin) ModifyBuilder(b *query.Builder) {
	b.
		WithDefaultSigil(ids.SigilExternal).
		WithRequireNonEmptyQuery()
}

func (cmd Checkin) Run(dep command.Request) {
	localWorkingCopy, queryGroup := cmd.MakeLocalWorkingCopyAndQueryGroup(
		dep,
		query.MakeBuilderOptions(cmd),
	)

	// envWorkspace := localWorkingCopy.GetEnvWorkspace()
	// envWorkspace.AssertInWorkspace(localWorkingCopy)

	op := user_ops.Checkin{
		Delete:             cmd.Delete,
		Organize:           cmd.Organize,
		Proto:              cmd.Proto,
		CheckoutBlobAndRun: cmd.CheckoutBlobAndRun,
		OpenBlob:           cmd.OpenBlob,
	}

	// TODO add auto dot operator
	if err := op.Run(localWorkingCopy, queryGroup); err != nil {
		dep.CancelWithError(err)
	}
}
