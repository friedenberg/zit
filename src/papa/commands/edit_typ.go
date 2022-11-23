package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/echo/typ"
	"github.com/friedenberg/zit/src/foxtrot/id_set"
	"github.com/friedenberg/zit/src/foxtrot/typ_checked_out"
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
			MutableId: &typ.Kennung{},
		},
	)

	return
}

func (c EditTyp) RunWithIds(u *umwelt.Umwelt, ids id_set.Set) (err error) {
	typen := typ.MakeMutableSet(ids.Typen()...)

	printerType := u.PrinterTypCheckedOut("checked out")

	if err = typen.Each(
		func(tk typ.Kennung) (err error) {
			t := &typ.Typ{
				Kennung: tk,
			}

			var tco *typ_checked_out.Typ

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
