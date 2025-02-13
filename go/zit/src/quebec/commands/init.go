package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

func init() {
	command.Register("init", &Init{})
}

type Init struct {
	command_components.Genesis
}

func (cmd *Init) SetFlagSet(f *flag.FlagSet) {
	cmd.Genesis.SetFlagSet(f)
}

func (cmd *Init) Run(req command.Request) {
	switch len(req.Args()) {
	case 0:
		req.CancelWithBadRequestf("expected a repo id, but got nothing")

	default:
		req.CancelWithBadRequestf("only acceptable argument is a repo id, but got %q", req.Args())

	case 1:
		break
	}

	cmd.OnTheFirstDay(req, req.Args()[0])
}
