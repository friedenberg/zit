package command

import "flag"

type Wrapper struct {
	*flag.FlagSet
	Command
}

func (wrapper Wrapper) GetFlagSet() *flag.FlagSet {
	return wrapper.FlagSet
}

func (wrapper Wrapper) SetFlagSet(f *flag.FlagSet) {
	wrapper.Command.SetFlagSet(f)
}
