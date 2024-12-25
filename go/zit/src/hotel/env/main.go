package env

import (
	"flag"
	"io"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/golf/mutable_config_blobs"
)

type Env struct {
	errors.Context

	in  fd.Std
	out fd.Std
	err fd.Std

	flags *flag.FlagSet

	dir_layout.Layout

	cliConfig mutable_config_blobs.Cli
}

func MakeDefault(
	layout dir_layout.Layout,
) *Env {
	return &Env{
		Context: errors.MakeContextDefault(),
		in:      fd.MakeStd(os.Stdin),
		out:     fd.MakeStd(os.Stdout),
		err:     fd.MakeStd(os.Stderr),
		Layout:  layout,
	}
}

func Make(
	context errors.Context,
	flags *flag.FlagSet,
	kCli mutable_config_blobs.Cli,
	primitiveFSHome dir_layout.Layout,
) *Env {
	return &Env{
		Context:   context,
		in:        fd.MakeStd(os.Stdin),
		out:       fd.MakeStd(os.Stdout),
		err:       fd.MakeStd(os.Stderr),
		flags:     flags,
		cliConfig: kCli,
		Layout:    primitiveFSHome,
	}
}

func (u *Env) GetIn() fd.Std {
	return u.in
}

func (u *Env) GetInFile() io.Reader {
	return u.in.File
}

func (u *Env) GetOut() fd.Std {
	return u.out
}

func (u *Env) GetOutFile() interfaces.WriterAndStringWriter {
	return u.out.File
}

func (u *Env) GetErr() fd.Std {
	return u.err
}

func (u *Env) GetErrFile() interfaces.WriterAndStringWriter {
	return u.err.File
}

func (u *Env) GetCLIConfig() mutable_config_blobs.Cli {
	return u.cliConfig
}

func (u *Env) GetDirLayoutPrimitive() dir_layout.Layout {
	return u.Layout
}
