package user_ops

import (
	"fmt"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/lima/organize_text"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
)

// TODO support using query results for organize population
type CreateOrganizeFile struct {
	*local_working_copy.Repo
	organize_text.Options
}

func (cmd CreateOrganizeFile) RunAndWrite(
	writer io.Writer,
) (results *organize_text.Text, err error) {
	if results, err = cmd.Run(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if _, err = results.WriteTo(writer); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (cmd CreateOrganizeFile) Run() (results *organize_text.Text, err error) {
	count := cmd.Options.Skus.Len()

	if cmd.Options.Limit == 0 && count > 30 && !cmd.GetCLIConfig().IsDryRun() {
		if !cmd.Confirm(
			fmt.Sprintf(
				"a large number (%d) of objects would be edited in organize. continue to organize?",
				count,
			),
		) {
			err = errors.BadRequestf("aborting")
			return
		}
	}

	if results, err = organize_text.New(cmd.Options); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
