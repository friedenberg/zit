package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/id_set"
	"github.com/friedenberg/zit/src/golf/typ"
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
			MutableId: &kennung.Typ{},
		},
	)

	return
}

func (c EditTyp) RunWithIds(u *umwelt.Umwelt, ids id_set.Set) (err error) {
	typen := typ.MakeMutableSet(ids.Typen()...)

	printerType := u.PrinterTypCheckedOut("checked out")

	if err = typen.Each(
		func(tk kennung.Typ) (err error) {
			t := &typ.Named{
				Kennung: tk,
			}

			var tco *typ.External

			if tco, err = u.StoreWorkingDirectory().WriteTyp(t); err != nil {
				err = errors.Wrap(err)
				return
			}

			if printerType(tco); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
