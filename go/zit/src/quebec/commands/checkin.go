package commands

import (
	"flag"
	"path/filepath"

	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_local"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
	"code.linenisgreat.com/zit/go/zit/src/papa/user_ops"
)

func init() {
	cmd := &Checkin{
		Proto: sku.MakeProto(nil),
	}

	command.Register("checkin", cmd)
	command.Register("add", cmd)
	command.Register("save", cmd)
}

type Checkin struct {
	command_components.LocalWorkingCopyWithQueryGroup

	complete command_components.Complete

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

	cmd.complete.SetFlagsProto(
		&cmd.Proto,
		f,
		"description to use for new zettels",
		"tags added for new zettels",
		"type used for new zettels",
	)

	cmd.Checkout.SetFlagSet(f)
}

// TODO refactor into common
func (cmd *Checkin) Complete(
	_ command.Request,
	envLocal env_local.Env,
	commandLine command.CommandLine,
) {
	searchDir := envLocal.GetCwd()

	if commandLine.InProgress != "" && files.Exists(commandLine.InProgress) {
		var err error

		if commandLine.InProgress, err = filepath.Abs(commandLine.InProgress); err != nil {
			envLocal.CancelWithError(err)
			return
		}

		if searchDir, err = filepath.Rel(searchDir, commandLine.InProgress); err != nil {
			envLocal.CancelWithError(err)
			return
		}
	}

	for dirEntry, err := range files.WalkDir(searchDir) {
		if err != nil {
			envLocal.CancelWithError(err)
			return
		}

		if files.WalkDirIgnoreFuncHidden(dirEntry) {
			continue
		}

		if !dirEntry.IsDir() {
			envLocal.GetUI().Printf("%s\tfile", dirEntry.RelPath)
		} else {
			envLocal.GetUI().Printf("%s/\tdirectory", dirEntry.RelPath)
		}
	}
}

func (c *Checkin) ModifyBuilder(b *query.Builder) {
	b.
		WithRequireNonEmptyQuery()
}

func (cmd Checkin) Run(dep command.Request) {
	localWorkingCopy, queryGroup := cmd.MakeLocalWorkingCopyAndQueryGroup(
		dep,
		query.BuilderOptionsOld(
			cmd,
			query.BuilderOptionDefaultSigil(ids.SigilExternal),
			query.BuilderOptionDefaultGenres(genres.All()...),
		),
	)

	// envWorkspace := localWorkingCopy.GetEnvWorkspace()
	// envWorkspace.AssertInWorkspace(localWorkingCopy)
	workspace := localWorkingCopy.GetEnvWorkspace()
	workspaceTags := workspace.GetDefaults().GetTags()

	for t := range workspaceTags.All() {
		cmd.Proto.Tags.Add(t)
	}

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
