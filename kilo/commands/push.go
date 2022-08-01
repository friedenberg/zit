package commands

import (
	"bytes"
	"flag"
	"os"
	"os/exec"

	"github.com/friedenberg/zit/alfa/errors"
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

func (c Push) Run(u _Umwelt, args ...string) (err error) {
	if len(args) == 0 {
		err = errors.Errorf("no remote specified")
		return
	}

	var remote _RemoteScript

	if remote, err = c.remoteScriptFromArg(u, args[0]); err != nil {
		err = errors.Error(err)
		return
	}

	if len(args) > 1 {
		args = args[1:]
	} else {
		args = []string{}
	}

	// var hins []_Hinweis

	// if _, hins, err = zs.Hinweisen().All(); err != nil {
	// 	err = errors.Error(err)
	// 	return
	// }

	// chains := make([]_ZettelsChain, len(hins))

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

func (c Push) remoteScriptFromArg(u _Umwelt, arg string) (remote _RemoteScript, err error) {
	ok := false

	if remote, ok = u.Konfig.RemoteScripts[arg]; !ok {
		p := u.DirZit("bin", arg)

		if !_FilesExist(p) {
			err = errors.Errorf("remote not defined: '%s'", arg)
			return
		}

		remote = _RemoteScriptFile{
			Path: p,
		}
	}

	return
}

func (c Push) runRemoteScript(u _Umwelt, remote _RemoteScript, args []string, b []byte) (err error) {
	var script *exec.Cmd

	if script, err = remote.Cmd(append([]string{"push"}, args...)); err != nil {
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
