package commands

import (
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
	dep command.Dep,
) {
	args := dep.Args()
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

	if err := localWorkingCopy.GetStore().CreateOrUpdate(
		sk,
		sku.StoreOptions{},
	); err != nil {
		localWorkingCopy.CancelWithError(err)
	}

	localWorkingCopy.Must(localWorkingCopy.Unlock)
}

func (c EditConfig) editInVim(
	u *local_working_copy.Repo,
) (sk *sku.Transacted, err error) {
	var f *os.File

	if f, err = u.GetRepoLayout().TempLocal.FileTemp(); err != nil {
		err = errors.Wrap(err)
		return
	}

	p := f.Name()

	if err = f.Close(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.makeTempConfigFile(u, p); err != nil {
		err = errors.Wrap(err)
		return
	}

	openVimOp := user_ops.OpenEditor{
		VimOptions: vim_cli_options_builder.New().
			WithFileType("zit-konfig").
			Build(),
	}

	if err = openVimOp.Run(u, p); err != nil {
		err = errors.Wrap(err)
		return
	}

	if sk, err = c.readTempConfigFile(u, p); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c EditConfig) makeTempConfigFile(
	u *local_working_copy.Repo,
	p string,
) (err error) {
	var k *sku.Transacted

	if k, err = u.GetStore().ReadTransactedFromObjectId(&ids.Config{}); err != nil {
		err = errors.Wrap(err)
		return
	}

	var i sku.FSItem
	i.Reset()

	if err = i.Object.Set(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	i.MutableSetLike.Add(&i.Object)

	if err = u.GetFileEncoder().Encode(
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
