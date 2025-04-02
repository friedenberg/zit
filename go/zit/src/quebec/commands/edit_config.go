package commands

import (
	"fmt"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/vim_cli_options_builder"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
	"code.linenisgreat.com/zit/go/zit/src/papa/user_ops"
)

func init() {
	command.Register("edit-config", &EditConfig{})
}

type EditConfig struct {
	command_components.LocalWorkingCopy
}

func (cmd EditConfig) Run(
	dep command.Request,
) {
	args := dep.PopArgs()
	localWorkingCopy := cmd.MakeLocalWorkingCopy(dep)

	if len(args) > 0 {
		ui.Err().Print("Command edit-konfig ignores passed in arguments.")
	}

	var sk *sku.Transacted

	{
		var err error

		if sk, err = cmd.editInVim(localWorkingCopy); err != nil {
			localWorkingCopy.CancelWithError(err)
		}
	}

	localWorkingCopy.Must(localWorkingCopy.Reset)
	localWorkingCopy.Must(localWorkingCopy.Lock)

	if err := localWorkingCopy.GetStore().CreateOrUpdateDefaultProto(
		sk,
		sku.StoreOptions{},
	); err != nil {
		localWorkingCopy.CancelWithError(err)
	}

	localWorkingCopy.Must(localWorkingCopy.Unlock)
}

func (c EditConfig) editInVim(
	repo *local_working_copy.Repo,
) (sk *sku.Transacted, err error) {
	var file *os.File

	if file, err = repo.GetEnvRepo().GetTempLocal().FileTempWithTemplate(
		fmt.Sprintf("*.%s", repo.GetConfig().GetFileExtensions().GetFileExtensionConfig()),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	path := file.Name()

	if err = file.Close(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.makeTempConfigFile(repo, path); err != nil {
		err = errors.Wrap(err)
		return
	}

	openVimOp := user_ops.OpenEditor{
		VimOptions: vim_cli_options_builder.New().
			WithFileType("zit-object").
			Build(),
	}

	if err = openVimOp.Run(repo, path); err != nil {
		err = errors.Wrap(err)
		return
	}

	if sk, err = c.readTempConfigFile(repo, path); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c EditConfig) makeTempConfigFile(
	repo *local_working_copy.Repo,
	path string,
) (err error) {
	var k *sku.Transacted

	if k, err = repo.GetStore().ReadTransactedFromObjectId(&ids.Config{}); err != nil {
		err = errors.Wrap(err)
		return
	}

	var i sku.FSItem
	i.Reset()

	if err = i.Object.Set(path); err != nil {
		err = errors.Wrap(err)
		return
	}

	i.MutableSetLike.Add(&i.Object)

	if err = repo.GetFileEncoder().Encode(
		checkout_options.TextFormatterOptions{},
		k,
		&i,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c EditConfig) readTempConfigFile(
	u *local_working_copy.Repo,
	p string,
) (sk *sku.Transacted, err error) {
	sk = sku.GetTransactedPool().Get()

	if sk.ObjectId.Set("konfig"); err != nil {
		err = errors.Wrap(err)
		return
	}

	var f *os.File

	if f, err = files.Open(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, f)

	if err = u.GetStore().GetStoreFS().ReadOneExternalObjectReader(
		f,
		sk,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
