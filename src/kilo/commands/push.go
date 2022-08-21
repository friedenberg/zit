package commands

import (
	"bytes"
	"flag"
	"os"
	"os/exec"

	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/charlie/files"
	"github.com/friedenberg/zit/src/delta/konfig"
	"github.com/friedenberg/zit/src/echo/umwelt"
)

type Push struct {
}

func init() {
	registerCommand(
		"push",
		func(f *flag.FlagSet) Command {
			c := &Push{}

			return c
		},
	)
}

func (c Push) Run(u *umwelt.Umwelt, args ...string) (err error) {
	if len(args) == 0 {
		err = errors.Normalf("no remote specified")
		return
	}

	var remote konfig.RemoteScript

	if remote, err = c.remoteScriptFromArg(u, args[0]); err != nil {
		err = errors.Error(err)
		return
	}

	if len(args) > 1 {
		args = args[1:]
	} else {
		args = []string{}
	}

	// var hins []hinweis.Hinweis

	// if _, hins, err = zs.Hinweisen().All(); err != nil {
	// 	err = errors.Error(err)
	// 	return
	// }

	// chains := make([]zettels.Chain, len(hins))

	// for i, h := range hins {
	// 	if chains[i], err = zs.AllInChain(h); err != nil {
	// 		err = errors.Error(err)
	// 		return
	// 	}
	// }

	// b, err := json.Marshal(chains)

	// if err != nil {
	// 	logz.Print(err)
	// 	return
	// }
	b := []byte{}

	if err = c.runRemoteScript(u, remote, args, b); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (c Push) remoteScriptFromArg(u *umwelt.Umwelt, arg string) (remote konfig.RemoteScript, err error) {
	ok := false

	if remote, ok = u.Konfig.RemoteScripts[arg]; !ok {
		p := u.DirZit("bin", arg)

		if !files.Exists(p) {
			err = errors.Errorf("remote not defined: '%s'", arg)
			return
		}

		remote = konfig.RemoteScriptFile{
			Path: p,
		}
	}

	return
}

func (c Push) runRemoteScript(u *umwelt.Umwelt, remote konfig.RemoteScript, args []string, b []byte) (err error) {
	var script *exec.Cmd

	if script, err = remote.Cmd(append([]string{"push"}, args...)...); err != nil {
		err = errors.Error(err)
		return
	}

	script.Stdin = os.Stdin
	script.Stdout = os.Stdout
	script.Stderr = os.Stderr

	r := bytes.NewBuffer(b)
	script.Stdin = r

	if err = script.Run(); err != nil {
		err = errors.Error(err)
		return
	}

	return
}
