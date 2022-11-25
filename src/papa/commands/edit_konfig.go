package commands

import (
	"flag"
	"io"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/vim_cli_options_builder"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/charlie/sha"
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
	var p string

	if p, err = c.partOne(u); err != nil {
		err = errors.Wrap(err)
		return
	}

	openVimOp := user_ops.OpenVim{
		Options: vim_cli_options_builder.New().
			WithCursorLocation(2, 3).
			WithFileType("toml").
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

	if k, err = c.partTwo(u, p); err != nil {
		err = errors.Wrap(err)
		return
	}

	errors.Err().Print(k)

	//TODO read, validate, and compile konfig file
	//TODO checkin new konfig

	return
}

func (c EditKonfig) partOne(
	u *umwelt.Umwelt,
) (p string, err error) {
	var k *konfig.Transacted

	if k, err = u.StoreObjekten().Konfig().Read(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var f *os.File

	if f, err = files.TempFile(); err != nil {
		err = errors.Wrap(err)
		return
	}

	p = f.Name()

	defer errors.Deferred(&err, f.Close)

	var ar sha.ReadCloser

	if ar, err = u.StoreObjekten().AkteReader(
		k.Named.Stored.Objekte.Sha,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, ar.Close)

	if _, err = io.Copy(f, ar); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c EditKonfig) partTwo(
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

	if _, err = format.ReadFormat(f, k); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
