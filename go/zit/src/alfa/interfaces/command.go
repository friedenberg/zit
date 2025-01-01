package interfaces

import "flag"

type CommandComponent interface {
	SetFlagSet(*flag.FlagSet)
}
