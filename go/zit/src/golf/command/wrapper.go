package command

import "flag"

type Wrapper struct {
	*flag.FlagSet
	Command2
}

func (wrapper Wrapper) GetFlagSet() *flag.FlagSet {
	return wrapper.FlagSet
}

func (wrapper Wrapper) SetFlagSet(f *flag.FlagSet) {
	wrapper.Command2.SetFlagSet(f)
}
