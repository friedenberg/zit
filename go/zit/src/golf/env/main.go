package env

import (
	"io"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/debug"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/config_mutable_cli"
)

type Env struct {
	*errors.Context

	options Options

	in  fd.Std
	ui  fd.Std
	out fd.Std
	err fd.Std

	dir_layout.Layout

	debug *debug.Context

	cliConfig config_mutable_cli.Config
}

func MakeDefault(
	layout dir_layout.Layout,
) *Env {
	return Make(
		errors.MakeContextDefault(),
		config_mutable_cli.Config{},
		layout,
		Options{},
	)
}

func Make(
	context *errors.Context,
	kCli config_mutable_cli.Config,
	dirLayout dir_layout.Layout,
	options Options,
) *Env {
	e := &Env{
		Context:   context,
		options:   options,
		in:        fd.MakeStd(os.Stdin),
		out:       fd.MakeStd(os.Stdout),
		err:       fd.MakeStd(os.Stderr),
		cliConfig: kCli,
		Layout:    dirLayout,
	}

	if options.UIFileIsStderr {
		e.ui = e.err
	} else {
		e.ui = e.out
	}

	{
		var err error

		if e.debug, err = debug.MakeContext(context, kCli.Debug); err != nil {
			context.CancelWithError(err)
		}
	}

	if kCli.Verbose && !kCli.Quiet {
		ui.SetVerbose(true)
	} else {
		ui.SetOutput(io.Discard)
	}

	if kCli.Todo {
		ui.SetTodoOn()
	}

	return e
}

func (u Env) GetOptions() Options {
	return u.options
}

func (u *Env) GetIn() fd.Std {
	return u.in
}

func (u *Env) GetInFile() io.Reader {
	return u.in.GetFile()
}

func (u *Env) GetUI() fd.Std {
	return u.ui
}

func (u *Env) GetUIFile() interfaces.WriterAndStringWriter {
	return u.ui.GetFile()
}

func (u *Env) GetOut() fd.Std {
	return u.out
}

func (u *Env) GetOutFile() interfaces.WriterAndStringWriter {
	return u.out.GetFile()
}

func (u *Env) GetErr() fd.Std {
	return u.err
}

func (u *Env) GetErrFile() interfaces.WriterAndStringWriter {
	return u.err.GetFile()
}

func (u *Env) GetCLIConfig() config_mutable_cli.Config {
	return u.cliConfig
}

func (u *Env) GetDirLayout() dir_layout.Layout {
	return u.Layout
}
