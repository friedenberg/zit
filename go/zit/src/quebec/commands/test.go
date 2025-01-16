package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/golf/command"
)

type Test struct{}

func init() {
	command.Register("test", &Test{})
}

func (*Test) SetFlagSet(*flag.FlagSet) {}

func (c Test) Run(dep command.Request) {}
