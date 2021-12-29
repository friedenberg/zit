package commands

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/friedenberg/zit/alfa/stdprinter"
)

type New struct {
	Filter       _ScriptValue
	ValidateOnly bool
}

func init() {
	registerCommand(
		"new",
		func(f *flag.FlagSet) Command {
			c := &New{}

			f.Var(&c.Filter, "filter", "a script to run for each file to transform it the standard zettel format")
			f.BoolVar(&c.ValidateOnly, "validate-only", false, "do not actually add the zettels, just validate the format")

			return commandWithZettels{c}
		},
	)
}

func (c New) RunWithZettels(u _Umwelt, zs _Zettels, args ...string) (err error) {
	f := _ZettelFormatsText{}

	if c.ValidateOnly {
		if len(args) == 0 {
			_Errf("when -validate-only is set, paths to existing zettels must be provided")
			return
		}

		for _, arg := range args {
			_, err := c.zettelsFromPath(u, zs, f, arg)

			if err != nil {
				_Errf("%s: err: %s\n", arg, err)
			} else {
				_Outf("%s: ok\n", arg)
			}

			continue
		}

		return
	}

	if len(args) == 0 {
		var toCreate string

		if toCreate, err = c.writeEmptyAndOpen(u, zs, f); err != nil {
			err = _Error(err)
			return
		}

		args = append(args, toCreate)
	}

	u.Lock.Lock()
	defer _PanicIfError(u.Lock.Unlock())

	for _, arg := range args {
		var toCreate []_Zettel

		if toCreate, err = c.zettelsFromPath(u, zs, f, arg); err != nil {
			err = _Errorf("zettel text format error for path: %s: %w", arg, err)
			return
		}

		for _, z := range toCreate {
			var named _NamedZettel

			if named, err = zs.Create(z); err != nil {
				c.handleStoreError(u, named, arg, err)
				err = nil
				return
			}

			_Outf("[%s %s]\n", named.Hinweis, named.Sha)
		}
	}

	return
}

func (c New) writeEmptyAndOpen(u _Umwelt, zs _Zettels, format _ZettelFormat) (out string, err error) {
	var f *os.File

	if f, err = _TempFileWithPattern("*.md"); err != nil {
		err = _Error(err)
		return
	}

	defer _Close(f)

	out = f.Name()

	etiketten := _EtikettNewSet()

	for e, t := range u.Konfig.Tags {
		if !t.AddToNewZettels {
			continue
		}

		if err = etiketten.AddString(e); err != nil {
			err = _Error(err)
			return
		}
	}

	ctx := _ZettelFormatContextWrite{
		Out:               f,
		AkteReaderFactory: zs,
		Zettel: _Zettel{
			Etiketten: etiketten,
			AkteExt:   _AkteExt{Value: "md"},
		},
	}

	if _, err = format.WriteTo(ctx); err != nil {
		err = _Error(err)
		return
	}

	vimArgs := []string{
		"-c",
		`call cursor(2, 3)`,
		"-c",
		`startinsert!`,
		"-c",
		"set ft=zit.zettel",
		"-c",
		"source ~/.vim/syntax/zit.zettel.vim",
	}

	if err = _OpenVimWithArgs(vimArgs, out); err != nil {
		err = _Error(err)
		return
	}

	return
}

func (c New) zettelsFromPath(u _Umwelt, zs _Zettels, format _ZettelFormat, p string) (out []_Zettel, err error) {
	var r io.Reader

	log.Print("running")

	if r, err = c.Filter.Run(p); err != nil {
		err = _Error(err)
		return
	}

	defer c.Filter.Close()

	ctx := _ZettelFormatContextRead{
		In:                r,
		AkteWriterFactory: zs,
	}

	if _, err = format.ReadFrom(&ctx); err != nil {
		err = _Error(err)
		return
	}

	if ctx.RecoverableError != nil {
		var errAkteInlineAndFilePath _ZettelFormatTextErrAkteInlineAndFilePath

		if _ErrorAs(ctx.RecoverableError, &errAkteInlineAndFilePath) {
			var z1 _Zettel

			if z1, err = errAkteInlineAndFilePath.Recover(); err != nil {
				err = _Error(err)
				return
			}

			out = append(out, z1)
		} else {
			err = _Errorf("unsupported recoverable error: %w", ctx.RecoverableError)
			return
		}
	}

	out = append(out, ctx.Zettel)

	return
}

func (c New) handleStoreError(u _Umwelt, z _NamedZettel, f string, in error) {
	var err error

	var lostError _VerlorenAndGefundenError
	var normalError _ErrorsStackTracer

	if _ErrorAs(in, &lostError) {
		var p string

		if p, err = lostError.AddToLostAndFound(u.DirZit("Verloren+Gefunden")); err != nil {
			stdprinter.Error(err)
			return
		}

		_Outf("lost+found: %s: %s\n", lostError.Error(), p)

	} else if _ErrorAs(in, &normalError) {
		_Errf("%s\n", normalError.Error())
	} else {
		err = _Errorf("writing zettel failed: %s: %w", f, in)
		stdprinter.Error(err)
	}
}
