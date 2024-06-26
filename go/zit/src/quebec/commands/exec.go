package commands

import (
	"flag"
	"io"
	"os"
	"os/exec"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/november/umwelt"
	"code.linenisgreat.com/zit/go/zit/src/papa/user_ops"
)

type Exec struct{}

func init() {
	registerCommand(
		"exec",
		func(f *flag.FlagSet) Command {
			c := &Exec{}
			return c
		},
	)
}

func (c Exec) Run(u *umwelt.Umwelt, args ...string) (err error) {
	if len(args) == 0 {
		err = errors.Normalf("needs at least Sku and possibly function name")
		return
	}

	k, args := args[0], args[1:]

	var sk *sku.Transacted

	if sk, err = u.GetSkuFromString(k); err != nil {
		err = errors.Wrap(err)
		return
	}

	switch {
	case strings.HasPrefix(sk.GetTyp().String(), "bash"):
		if err = c.runBash(u, sk, args...); err != nil {
			err = errors.Wrap(err)
			return
		}

	case strings.HasPrefix(sk.GetTyp().String(), "lua"):
		execLuaOp := user_ops.ExecLua{
			Umwelt: u,
		}

		if err = execLuaOp.Run(sk, args...); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (c Exec) runBash(
	u *umwelt.Umwelt,
	tz *sku.Transacted,
	args ...string,
) (err error) {
	var scriptPath string

	func() {
		var ar io.ReadCloser

		if ar, err = u.Standort().AkteReader(
			tz.GetAkteSha(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		var f *os.File

		if f, err = u.Standort().FileTempLocal(); err != nil {
			err = errors.Wrap(err)
			return
		}

		scriptPath = f.Name()

		defer errors.DeferredCloser(&err, f)

		if _, err = io.Copy(f, ar); err != nil {
			err = errors.Wrap(err)
			return
		}
	}()

	cmd := exec.Command(
		"bash",
		append([]string{scriptPath}, args...)...,
	)

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if err = cmd.Run(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
