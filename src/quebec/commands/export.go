package commands

import (
	"flag"
	"os"
	"os/exec"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/india/konfig"
	"github.com/friedenberg/zit/src/oscar/umwelt"
)

type Export struct {
}

func init() {
	registerCommand(
		"export",
		func(f *flag.FlagSet) Command {
			c := &Export{}

			return c
		},
	)
}

func (c Export) Run(u *umwelt.Umwelt, args ...string) (err error) {
	if len(args) == 0 {
		err = errors.Normalf("no remote specified")
		return
	}

	var remote konfig.RemoteScript

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
	if err = c.runRemoteScript(u, remote, args); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c Export) remoteScriptFromArg(
	u *umwelt.Umwelt,
	arg string,
) (remote konfig.RemoteScript, err error) {
	p := u.Standort().DirZit("bin", arg)

	if !files.Exists(p) {
		err = errors.Errorf("remote not defined: '%s'", arg)
		return
	}

	remote = konfig.RemoteScriptFile{
		Path: p,
	}

	return
}

func (c Export) runRemoteScript(
	u *umwelt.Umwelt,
	remote konfig.RemoteScript,
	args []string,
) (err error) {
	var script *exec.Cmd

	if script, err = remote.Cmd(append([]string{"push"}, args...)...); err != nil {
		err = errors.Wrap(err)
		return
	}

	script.Stdin = os.Stdin
	script.Stdout = os.Stdout
	script.Stderr = os.Stderr

	if err = script.Run(); err != nil {
		err = errors.Normal(err)
		return
	}

	return
}