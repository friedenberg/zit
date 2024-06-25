package commands

import (
	"flag"
	"io"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/alfa/vim_cli_options_builder"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/erworben"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/november/umwelt"
	"code.linenisgreat.com/zit/go/zit/src/papa/user_ops"
)

type EditSchlummernd struct{}

func init() {
	registerCommand(
		"schlummernd-edit",
		func(f *flag.FlagSet) Command {
			c := &EditSchlummernd{}

			return c
		},
	)
}

func (c EditSchlummernd) Run(u *umwelt.Umwelt, args ...string) (err error) {
	if len(args) > 0 {
		ui.Err().Print("Command edit-konfig ignores passed in arguments.")
	}

	var sh schnittstellen.ShaLike

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

func (c EditSchlummernd) editInVim(
	u *umwelt.Umwelt,
) (sh schnittstellen.ShaLike, err error) {
	var p string

	if p, err = c.makeTempKonfigFile(u); err != nil {
		err = errors.Wrap(err)
		return
	}

	openVimOp := user_ops.OpenVim{
		Options: vim_cli_options_builder.New().
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

func (c EditSchlummernd) makeTempKonfigFile(
	u *umwelt.Umwelt,
) (p string, err error) {
	var k *sku.Transacted

	if k, err = u.GetStore().ReadOne(&kennung.Konfig{}); err != nil {
		err = errors.Wrap(err)
		return
	}

	var f *os.File

	if f, err = u.Standort().FileTempLocal(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, f)

	p = f.Name()

	format := u.GetStore().GetKonfigAkteFormat()

	if _, err = format.FormatSavedAkte(f, k.GetAkteSha()); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c EditSchlummernd) readTempKonfigFile(
	u *umwelt.Umwelt,
	p string,
) (sh schnittstellen.ShaLike, err error) {
	var f *os.File

	if f, err = files.Open(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, f)

	format := u.GetStore().GetKonfigAkteFormat()

	var k erworben.Akte

	var aw schnittstellen.ShaWriteCloser

	if aw, err = u.Standort().AkteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, aw)

	// TODO-P3 offer option to edit again
	if _, err = format.ParseAkte(io.TeeReader(f, aw), &k); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh = aw.GetShaLike()

	return
}
