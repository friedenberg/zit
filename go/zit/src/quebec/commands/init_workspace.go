package commands

import (
	"flag"
	"os"
	"path/filepath"

	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/golf/config_mutable_blobs"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_local"
	"code.linenisgreat.com/zit/go/zit/src/hotel/workspace_config_blobs"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

func init() {
	command.Register(
		"init-workspace",
		&InitWorkspace{},
	)
}

type InitWorkspace struct {
	command_components.Env
	command_components.LocalWorkingCopy

	DefaultQueryGroup string
	Proto             sku.Proto
}

func (cmd *InitWorkspace) SetFlagSet(f *flag.FlagSet) {
	cmd.LocalWorkingCopy.SetFlagSet(f)
	// TODO add command.Completer variants of tags, type, and query flags

	cmd.Proto.SetFlagSetTags(
		f,
		"tags added for new objects in `checkin`, `new`, `organize`",
	)

	cmd.Proto.SetFlagSetType(
		f,
		"type used for new objects in `new` and `organize`",
	)

	f.StringVar(&cmd.DefaultQueryGroup, "query", "", "default query for `show`")
}

func (cmd InitWorkspace) Complete(
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

		if !dirEntry.IsDir() {
			continue
		}

		if files.WalkDirIgnoreFuncHidden(dirEntry) {
			continue
		}

		envLocal.GetUI().Printf("%s/\tdirectory", dirEntry.RelPath)
	}
}

func (cmd InitWorkspace) Run(req command.Request) {
	argc := len(req.Args())
	envLocal := cmd.MakeEnv(req)

	switch argc {
	case 0:
		break

	default:
		if !envLocal.Confirm(
			"more than one argument passed in. ignore everything after first and continue?",
		) {
			req.CancelWithBadRequestf("aborted")
			return
		}

		fallthrough

	case 1:
		dir := req.Args()[0]

		if err := envLocal.MakeDir(dir); err != nil {
			req.CancelWithError(err)
			return
		}

		if err := os.Chdir(dir); err != nil {
			req.CancelWithError(err)
			return
		}
	}

	localWorkingCopy := cmd.MakeLocalWorkingCopy(req)

	blob := &workspace_config_blobs.V0{
		Query: cmd.DefaultQueryGroup,
		Defaults: config_mutable_blobs.DefaultsV1OmitEmpty{
			Type: cmd.Proto.Type,
			Tags: quiter.Elements(cmd.Proto.Tags),
		},
	}

	if err := localWorkingCopy.GetEnvWorkspace().CreateWorkspace(
		blob,
	); err != nil {
		req.CancelWithError(err)
	}
}
