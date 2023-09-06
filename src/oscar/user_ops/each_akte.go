package user_ops

import (
	"fmt"
	"os/exec"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/november/umwelt"
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
	args, err = shlex.Split(utility)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	args = append(args, akten...)

	cmd := exec.Command(args[0], args[1:]...)

	if err = cmd.Start(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
