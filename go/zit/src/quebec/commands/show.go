package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/options_print"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/box_format"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

func init() {
	registerCommand("show", &Show{})
}

type Show struct {
	command_components.RepoLayout
	command_components.LocalArchive
	command_components.QueryGroup

	After  ids.Tai
	Before ids.Tai
	Format string
}

func (cmd *Show) SetFlagSet(f *flag.FlagSet) {
	cmd.RepoLayout.SetFlagSet(f)
	cmd.LocalArchive.SetFlagSet(f)
	cmd.QueryGroup.SetFlagSet(f)

	f.StringVar(&cmd.Format, "format", "log", "format")
	f.Var((*ids.TaiRFC3339Value)(&cmd.Before), "before", "")
	f.Var((*ids.TaiRFC3339Value)(&cmd.After), "after", "")
}

func (cmd Show) CompletionGenres() ids.Genre {
	return ids.MakeGenre(
		genres.Zettel,
		genres.Tag,
		genres.Type,
		genres.InventoryList,
		genres.Repo,
	)
}

func (cmd Show) DefaultGenres() ids.Genre {
	return ids.MakeGenre(
		genres.Zettel,
	)
}

func (cmd Show) Run(dep command.Dep) {
	repoLayout := cmd.MakeRepoLayout(dep, false)

	archive := cmd.MakeLocalArchive(repoLayout)

	args := dep.Args()

	if localWorkingCopy, ok := archive.(*local_working_copy.Repo); ok {
		switch {
		case archive.GetRepoLayout().Env.GetCLIConfig().Complete:
			cmd.CompleteWithRepo(cmd, localWorkingCopy, args...)

		default:
			qg := cmd.MakeQueryGroup(
				query.MakeBuilderOptions(cmd),
				localWorkingCopy,
				args,
			)

			cmd.runWithLocalWorkingCopyAndQuery(localWorkingCopy, qg)
		}
	} else {
		if len(args) != 0 {
			ui.Err().Print("ignoring arguments for archive repo")
		}

		cmd.runWithArchive(archive)
	}
}

func (cmd Show) runWithLocalWorkingCopyAndQuery(
	repo *local_working_copy.Repo,
	qg *query.Group,
) {
	var f interfaces.FuncIter[*sku.Transacted]

	if cmd.Format == "" && qg.IsExactlyOneObjectId() {
		cmd.Format = "text"
	}

	{
		var err error

		if f, err = repo.MakeFormatFunc(cmd.Format, repo.GetUIFile()); err != nil {
			repo.CancelWithError(err)
		}
	}

	if !cmd.Before.IsEmpty() {
		old := f

		f = func(sk *sku.Transacted) (err error) {
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
		old := f

		f = func(sk *sku.Transacted) (err error) {
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
		qg,
		quiter.MakeSyncSerializer(f),
	); err != nil {
		repo.CancelWithError(err)
	}
}

func (cmd Show) runWithArchive(
	archive repo.Archive,
	// qg *query.Group,
) {
	boxFormat := box_format.MakeBoxTransactedArchive(
		archive.GetRepoLayout().Env,
		options_print.V0{}.WithPrintTai(true),
	)

	f := string_format_writer.MakeDelim(
		"\n",
		archive.GetRepoLayout().GetUIFile(),
		string_format_writer.MakeFunc(
			func(w interfaces.WriterAndStringWriter, o *sku.Transacted) (n int64, err error) {
				return boxFormat.WriteStringFormat(w, o)
			},
		),
	)

	inventoryListStore := archive.GetInventoryListStore()

	if err := inventoryListStore.ReadAllSkus(
		func(_, sk *sku.Transacted) (err error) {
			if err = f(sk); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		archive.GetRepoLayout().CancelWithError(err)
	}
}
