package commands

import (
	"flag"

	"github.com/friedenberg/zit/bravo/stdprinter"
	"github.com/friedenberg/zit/bravo/node_type"
	"github.com/friedenberg/zit/charlie/age"
	"github.com/friedenberg/zit/delta/objekte"
	"github.com/friedenberg/zit/delta/umwelt"
)

type WriteObjekte struct {
	Type node_type.Type
}

func init() {
	registerCommand(
		"write-objekte",
		func(f *flag.FlagSet) Command {
			c := &WriteObjekte{
				Type: node_type.TypeUnknown,
			}

			f.Var(&c.Type, "type", "ObjekteType")

			return commandWithAge{c}
		},
	)
}

func (c WriteObjekte) RunWithAge(u *umwelt.Umwelt, age age.Age, args ...string) (err error) {
	objektePath, err := objekte.WriteAndMove(u.In, age, u.DirZit(), c.Type)

	stdprinter.Out(objektePath)

	return
}
