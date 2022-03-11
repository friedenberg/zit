package commands

import (
	"flag"
)

type WriteObjekte struct {
	Type _Type
}

func init() {
	registerCommand(
		"write-objekte",
		func(f *flag.FlagSet) Command {
			c := &WriteObjekte{
				Type: _TypeUnknown,
			}

			f.Var(&c.Type, "type", "ObjekteType")

			return commandWithAge{c}
		},
	)
}

func (c WriteObjekte) RunWithAge(u _Umwelt, age _Age, args ...string) (err error) {
	objektePath, err := _ObjekteWriteAndMove(u.In, age, u.DirZit(), c.Type)

	_Out(objektePath)

	return
}
