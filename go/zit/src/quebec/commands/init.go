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
	repoId := req.PopArg("repo-id")
	req.AssertNoMoreArgs()
	cmd.OnTheFirstDay(req, repoId)
}
