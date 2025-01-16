package command

import "flag"

type commandWrapper struct {
	*flag.FlagSet
	Command2
}

func (wrapper commandWrapper) GetFlagSet() *flag.FlagSet {
	return wrapper.FlagSet
}

func (wrapper commandWrapper) SetFlagSet(f *flag.FlagSet) {
	wrapper.Command2.SetFlagSet(f)
}
