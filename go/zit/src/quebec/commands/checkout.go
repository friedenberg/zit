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
		"checkout",
		&Checkout{
			CheckoutOptions: checkout_options.Options{
				CheckoutMode: checkout_mode.MetadataOnly,
			},
		},
	)
}

type Checkout struct {
	command_components.LocalWorkingCopyWithQueryGroup

	CheckoutOptions checkout_options.Options
	Organize        bool
}

func (c *Checkout) SetFlagSet(f *flag.FlagSet) {
	c.LocalWorkingCopyWithQueryGroup.SetFlagSet(f)
	f.BoolVar(&c.Organize, "organize", false, "")
	c.CheckoutOptions.SetFlagSet(f)
}

func (c Checkout) DefaultGenres() ids.Genre {
	return ids.MakeGenre(
		genres.Zettel,
	)
}

func (c Checkout) ModifyBuilder(b *query.Builder) {
	b.
		WithPermittedSigil(ids.SigilLatest).
		WithPermittedSigil(ids.SigilHidden).
		WithDefaultGenres(ids.MakeGenre(genres.Zettel)).
		WithRequireNonEmptyQuery()
}

func (cmd Checkout) Run(req command.Request) {
	localWorkingCopy, queryGroup := cmd.MakeLocalWorkingCopyAndQueryGroup(
		req,
		query.MakeBuilderOptions(cmd),
	)

	localWorkingCopy.AssertCLINotComplete()

	opCheckout := user_ops.Checkout{
		Repo:     localWorkingCopy,
		Organize: cmd.Organize,
		Options:  cmd.CheckoutOptions,
	}

	envWorkspace := localWorkingCopy.GetEnvWorkspace()
	envWorkspace.AssertInWorkspaceOrOfferToCreate(localWorkingCopy)

	if _, err := opCheckout.RunQuery(queryGroup); err != nil {
		localWorkingCopy.CancelWithError(err)
	}
}
