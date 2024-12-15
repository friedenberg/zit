package commands

import (
	"flag"
	"io"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/alfa/vim_cli_options_builder"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/mutable_config_blobs"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
	"code.linenisgreat.com/zit/go/zit/src/papa/user_ops"
)

type DormantEdit struct{}

func init() {
	registerCommand(
		"schlummernd-edit",
		func(f *flag.FlagSet) Command {
			c := &DormantEdit{}

			return c
		},
	)
}

func (c DormantEdit) Run(u *env.Local, args ...string) (err error) {
	if len(args) > 0 {
		ui.Err().Print("Command edit-konfig ignores passed in arguments.")
	}

	var sh interfaces.Sha

	if sh, err = c.editInVim(u); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.Reset(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, u.Unlock)

	if _, err = u.GetStore().UpdateKonfig(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c DormantEdit) editInVim(
	u *env.Local,
) (sh interfaces.Sha, err error) {
	var p string

	if p, err = c.makeTempKonfigFile(u); err != nil {
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

	if sh, err = c.readTempKonfigFile(u, p); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c DormantEdit) makeTempKonfigFile(
	u *env.Local,
) (p string, err error) {
	var k *sku.Transacted

	if k, err = u.GetStore().ReadTransactedFromObjectId(&ids.Config{}); err != nil {
		err = errors.Wrap(err)
		return
	}

	var f *os.File

	if f, err = u.GetDirectoryLayout().TempLocal.FileTemp(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, f)

	p = f.Name()

	format := u.GetStore().GetConfigBlobFormat()

	if _, err = format.FormatSavedBlob(f, k.GetBlobSha()); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c DormantEdit) readTempKonfigFile(
	u *env.Local,
	p string,
) (sh interfaces.Sha, err error) {
	var f *os.File

	if f, err = files.Open(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, f)

	format := u.GetStore().GetConfigBlobFormat()

	var k mutable_config_blobs.V0

	var aw interfaces.ShaWriteCloser

	if aw, err = u.GetDirectoryLayout().BlobWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, aw)

	// TODO-P3 offer option to edit again
	if _, err = format.ParseBlob(io.TeeReader(f, aw), &k); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh = aw.GetShaLike()

	return
}
