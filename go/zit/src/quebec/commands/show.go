package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/golf/env_ui"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_local"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/box_format"
	pkg_query "code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

func init() {
	command.Register("show", &Show{})
}

type Show struct {
	command_components.EnvRepo
	command_components.LocalArchive
	command_components.Query

	complete command_components.Complete

	After  ids.Tai
	Before ids.Tai
	Format string
}

func (cmd *Show) SetFlagSet(flagSet *flag.FlagSet) {
	cmd.EnvRepo.SetFlagSet(flagSet)
	cmd.LocalArchive.SetFlagSet(flagSet)
	cmd.Query.SetFlagSet(flagSet)

	flagSet.StringVar(&cmd.Format, "format", "log", "format")
	flagSet.Var((*ids.TaiRFC3339Value)(&cmd.Before), "before", "")
	flagSet.Var((*ids.TaiRFC3339Value)(&cmd.After), "after", "")
}

func (cmd Show) CompletionGenres() ids.Genre {
	return ids.MakeGenre(
		genres.Tag,
	)
}

func (cmd Show) Complete(
	req command.Request,
	envLocal env_local.Env,
	commandLine command.CommandLine,
) {
	envRepo := cmd.MakeEnvRepo(req, false)
	archive := cmd.MakeLocalArchive(envRepo)

	if localWorkingCopy, ok := archive.(*local_working_copy.Repo); ok {
		args := commandLine.FlagsOrArgs[1:]

		if commandLine.InProgress != "" {
			args = args[:len(args)-1]
		}

		cmd.complete.CompleteObjects(
			req,
			localWorkingCopy,
			pkg_query.BuilderOptionDefaultGenres(genres.Tag),
			args...,
		)
	}
}

func (cmd Show) Run(req command.Request) {
	envRepo := cmd.MakeEnvRepo(req, false)
	archive := cmd.MakeLocalArchive(envRepo)

	args := req.PopArgs()

	if localWorkingCopy, ok := archive.(*local_working_copy.Repo); ok {
		query := cmd.MakeQuery(
			req,
			pkg_query.BuilderOptions(
				pkg_query.BuilderOptionWorkspace{Env: localWorkingCopy.GetEnvWorkspace()},
				pkg_query.BuilderOptionDefaultGenres(genres.Zettel),
			),
			localWorkingCopy,
			args,
		)

		cmd.runWithLocalWorkingCopyAndQuery(localWorkingCopy, query)
	} else {
		if len(args) != 0 {
			ui.Err().Print("ignoring arguments for archive repo")
		}

		cmd.runWithArchive(envRepo, archive)
	}
}

func (cmd Show) runWithLocalWorkingCopyAndQuery(
	repo *local_working_copy.Repo,
	query *pkg_query.Query,
) {
	var output interfaces.FuncIter[*sku.Transacted]

	if cmd.Format == "" && pkg_query.IsExactlyOneObjectId(query) {
		cmd.Format = "text"
	}

	{
		var err error

		if output, err = repo.MakeFormatFunc(cmd.Format, repo.GetUIFile()); err != nil {
			repo.CancelWithError(err)
		}
	}

	if !cmd.Before.IsEmpty() {
		old := output

		output = func(sk *sku.Transacted) (err error) {
			if !sk.GetTai().Before(cmd.Before) {
				return
			}

			if err = old(sk); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}
	}

	if !cmd.After.IsEmpty() {
		old := output

		output = func(sk *sku.Transacted) (err error) {
			if !sk.GetTai().After(cmd.After) {
				return
			}

			if err = old(sk); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}
	}

	if err := repo.GetStore().QueryTransacted(
		query,
		quiter.MakeSyncSerializer(output),
	); err != nil {
		repo.CancelWithError(err)
	}
}

// TODO add support for query group
func (cmd Show) runWithArchive(
	env env_ui.Env,
	archive repo.Repo,
) {
	boxFormat := box_format.MakeBoxTransactedArchive(
		env,
		env.GetCLIConfig().PrintOptions,
	)

	printer := string_format_writer.MakeDelim(
		"\n",
		env.GetUIFile(),
		string_format_writer.MakeFunc(
			func(
				writer interfaces.WriterAndStringWriter,
				object *sku.Transacted,
			) (n int64, err error) {
				env.ContinueOrPanicOnDone()
				return boxFormat.EncodeStringTo(object, writer)
			},
		),
	)

	inventoryListStore := archive.GetInventoryListStore()

	if err := inventoryListStore.ReadAllSkus(
		func(_, sk *sku.Transacted) (err error) {
			if err = printer(sk); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		env.CancelWithError(err)
	}
}
