package commands

import (
	"bytes"
	"flag"
	"os"
	"os/exec"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/hotel/erworben"
	"github.com/friedenberg/zit/src/oscar/umwelt"
)

type Import struct {
}

func init() {
	registerCommand(
		"import",
		func(f *flag.FlagSet) Command {
			c := &Import{}

			return c
		},
	)
}

func (c Import) Run(u *umwelt.Umwelt, args ...string) (err error) {
	if len(args) == 0 {
		err = errors.Errorf("no remote specified")
		return
	}

	var remote erworben.RemoteScript

	if remote, err = c.remoteScriptFromArg(u, args[0]); err != nil {
		err = errors.Wrap(err)
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
		err = errors.Wrap(err)
		return
	}

	return
}

func (c Import) remoteScriptFromArg(
	u *umwelt.Umwelt,
	arg string,
) (remote erworben.RemoteScript, err error) {
	p := u.Standort().DirZit("bin", arg)

	if !files.Exists(p) {
		err = errors.Errorf("remote not defined: '%s'", arg)
		return
	}

	remote = erworben.RemoteScriptFile{
		Path: p,
	}

	return
}

func (c Import) runRemoteScript(
	u *umwelt.Umwelt,
	remote erworben.RemoteScript,
	args []string,
	b []byte,
) (err error) {
	var script *exec.Cmd

	if script, err = remote.Cmd(append([]string{"pull"}, args...)...); err != nil {
		err = errors.Wrap(err)
		return
	}

	script.Stdout = os.Stdout
	script.Stderr = os.Stderr

	r := bytes.NewBuffer(b)
	script.Stdin = r

	if err = script.Run(); err != nil {
		err = errors.Normal(err)
		return
	}

	return
}
