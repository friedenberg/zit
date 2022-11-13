package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/echo/typ"
	"github.com/friedenberg/zit/src/foxtrot/id_set"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type EditTyp struct {
}

func init() {
	registerCommand(
		"edit-typ",
		func(f *flag.FlagSet) Command {
			c := &EditTyp{}

			return commandWithIds{CommandWithIds: c}
		},
	)
}

func (c EditTyp) ProtoIdSet(u *umwelt.Umwelt) (is id_set.ProtoIdSet) {
	is = id_set.MakeProtoIdSet(
		id_set.ProtoId{
			MutableId: &typ.Typ{},
		},
	)

	return
}

func (c EditTyp) RunWithIds(u *umwelt.Umwelt, ids id_set.Set) (err error) {
	errors.PrintOut(ids.Typen())

	return
}
