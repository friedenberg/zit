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

func (c Checkout) ModifyBuilder(b *query.Builder) {
	b.
		WithPermittedSigil(ids.SigilLatest).
		WithPermittedSigil(ids.SigilHidden).
		WithRequireNonEmptyQuery()
}

func (cmd Checkout) Run(req command.Request) {
	localWorkingCopy := cmd.MakeLocalWorkingCopy(req)
	envWorkspace := localWorkingCopy.GetEnvWorkspace()

	queryGroup := cmd.MakeQueryIncludingWorkspace(
		req,
		query.BuilderOptions(
			query.BuilderOptionsOld(cmd),
			query.BuilderOptionWorkspace{Env: envWorkspace},
			query.BuilderOptionDefaultGenres(genres.Zettel),
		),
		localWorkingCopy,
		req.PopArgs(),
	)

	opCheckout := user_ops.Checkout{
		Repo:     localWorkingCopy,
		Organize: cmd.Organize,
		Options:  cmd.CheckoutOptions,
	}

	envWorkspace.AssertNotTemporaryOrOfferToCreate(localWorkingCopy)

	if _, err := opCheckout.RunQuery(queryGroup); err != nil {
		localWorkingCopy.CancelWithError(err)
	}
}
