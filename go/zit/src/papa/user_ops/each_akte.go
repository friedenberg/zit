package user_ops

import (
	"fmt"
	"os/exec"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/november/umwelt"
	"github.com/google/shlex"
)

type EachAkte struct{}

func (c EachAkte) Run(
	u *umwelt.Umwelt,
	utility string,
	akten ...string,
) (err error) {
	if len(akten) == 0 {
		return
	}

	v := fmt.Sprintf("running utility: %q", utility)

	if err = u.PrinterHeader()(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	var args []string

	if args, err = shlex.Split(utility); err != nil {
		err = errors.Wrap(err)
		return
	}

	args = append(args, akten...)

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = u.Out()
	cmd.Stdin = u.In()
	cmd.Stderr = u.Err()

	if err = cmd.Start(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
