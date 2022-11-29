package commands

import (
	"flag"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/vim_cli_options_builder"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/echo/konfig"
	"github.com/friedenberg/zit/src/november/umwelt"
	"github.com/friedenberg/zit/src/oscar/user_ops"
)

type EditKonfig struct {
}

func init() {
	registerCommand(
		"edit-konfig",
		func(f *flag.FlagSet) Command {
			c := &EditKonfig{}

			return c
		},
	)
}

func (c EditKonfig) Run(u *umwelt.Umwelt, args ...string) (err error) {
	if len(args) > 0 {
		errors.Err().Print("Command edit-konfig ignores passed in arguments.")
	}

	var p string

	if p, err = c.makeTempKonfigFile(u); err != nil {
		err = errors.Wrap(err)
		return
	}

	openVimOp := user_ops.OpenVim{
		Options: vim_cli_options_builder.New().
			WithCursorLocation(2, 3).
			WithFileType("zit-konfig").
			WithInsertMode().
			Build(),
	}

	if _, err = openVimOp.Run(u, p); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.Reset(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var k *konfig.Objekte

	if k, err = c.readTempKonfigFile(u, p); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, u.Unlock)

	if _, err = u.StoreObjekten().Konfig().Update(k); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c EditKonfig) makeTempKonfigFile(
	u *umwelt.Umwelt,
) (p string, err error) {
	var k konfig.Transacted

	if k, err = u.StoreObjekten().Konfig().Read(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var f *os.File

	if f, err = files.TempFile(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, f.Close)

	p = f.Name()

	format := konfig.MakeFormatText(u.StoreObjekten())

	if _, err = format.WriteFormat(f, &k.Objekte); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c EditKonfig) readTempKonfigFile(
	u *umwelt.Umwelt,
	p string,
) (k *konfig.Objekte, err error) {
	var f *os.File

	if f, err = files.Open(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, f.Close)

	format := konfig.MakeFormatText(u.StoreObjekten())

	k = &konfig.Objekte{}

	//TODO offer option to edit again
	if _, err = format.ReadFormat(f, k); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
