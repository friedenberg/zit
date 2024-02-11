package sku_fmt

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/erworben_cli_print_options"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/charlie/string_format_writer"
	"code.linenisgreat.com/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/src/echo/bezeichnung"
	"code.linenisgreat.com/zit/src/echo/fd"
	"code.linenisgreat.com/zit/src/echo/kennung"
)

type cliCheckedOut struct {
	options erworben_cli_print_options.PrintOptions

	writeTyp         bool
	writeBezeichnung bool
	writeEtiketten   bool

	rightAlignedWriter            schnittstellen.StringFormatWriter[string]
	shaStringFormatWriter         schnittstellen.StringFormatWriter[schnittstellen.ShaLike]
	kennungStringFormatWriter     schnittstellen.StringFormatWriter[*kennung.Kennung2]
	fdStringFormatWriter          schnittstellen.StringFormatWriter[*fd.FD]
	typStringFormatWriter         schnittstellen.StringFormatWriter[*kennung.Typ]
	bezeichnungStringFormatWriter schnittstellen.StringFormatWriter[*bezeichnung.Bezeichnung]
	etikettenStringFormatWriter   schnittstellen.StringFormatWriter[*kennung.Etikett]
}

func MakeCliCheckedOutFormat(
	options erworben_cli_print_options.PrintOptions,
	shaStringFormatWriter schnittstellen.StringFormatWriter[schnittstellen.ShaLike],
	fdStringFormatWriter schnittstellen.StringFormatWriter[*fd.FD],
	kennungStringFormatWriter schnittstellen.StringFormatWriter[*kennung.Kennung2],
	typStringFormatWriter schnittstellen.StringFormatWriter[*kennung.Typ],
	bezeichnungStringFormatWriter schnittstellen.StringFormatWriter[*bezeichnung.Bezeichnung],
	etikettenStringFormatWriter schnittstellen.StringFormatWriter[*kennung.Etikett],
) *cliCheckedOut {
	return &cliCheckedOut{
		options:                       options,
		writeTyp:                      true,
		writeBezeichnung:              true,
		writeEtiketten:                true,
		rightAlignedWriter:            string_format_writer.MakeRightAligned(),
		shaStringFormatWriter:         shaStringFormatWriter,
		kennungStringFormatWriter:     kennungStringFormatWriter,
		fdStringFormatWriter:          fdStringFormatWriter,
		typStringFormatWriter:         typStringFormatWriter,
		bezeichnungStringFormatWriter: bezeichnungStringFormatWriter,
		etikettenStringFormatWriter:   etikettenStringFormatWriter,
	}
}

func (f *cliCheckedOut) WriteStringFormat(
	sw schnittstellen.WriterAndStringWriter,
	co *CheckedOut,
) (n int64, err error) {
	var (
		n1 int
		n2 int64
	)

	{
		var stateString string

		if co.State == checked_out_state.StateError {
			stateString = co.Error.Error()
		} else {
			stateString = co.State.String()
		}

		n2, err = f.rightAlignedWriter.WriteStringFormat(sw, stateString)
		n += n2

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	o := &co.External
	fds := o.GetFDsPtr()
	n1, err = sw.WriteString("[")
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	m := fds.GetCheckoutMode()

	switch {
	case co.IsImport:
		fallthrough

	case m == checkout_mode.ModeAkteOnly:
		n2, err = f.kennungStringFormatWriter.WriteStringFormat(sw, &o.Kennung)
		n += n2

		if err != nil {
			err = errors.Wrap(err)
			return
		}

	default:
		n2, err = f.fdStringFormatWriter.WriteStringFormat(sw, &fds.Objekte)
		n += n2

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if co.State == checked_out_state.StateConflicted {
		n1, err = sw.WriteString("]")
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	n1, err = sw.WriteString("@")
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n2, err = f.shaStringFormatWriter.WriteStringFormat(sw, o.GetAkteSha())
	n += n2

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	if f.writeTyp {
		t := o.GetMetadatei().GetTypPtr()

		if len(t.String()) > 0 {
			n1, err = sw.WriteString(" !")
			n += int64(n1)

			if err != nil {
				err = errors.Wrap(err)
				return
			}

			n2, err = f.typStringFormatWriter.WriteStringFormat(sw, t)
			n += n2

			if err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	didWriteBezeichnung := false
	if f.writeBezeichnung {
		b := o.GetMetadatei().GetBezeichnungPtr()

		if !b.IsEmpty() {
			didWriteBezeichnung = true

			n1, err = sw.WriteString(" \"")
			n += int64(n1)

			if err != nil {
				err = errors.Wrap(err)
				return
			}

			n2, err = f.bezeichnungStringFormatWriter.WriteStringFormat(sw, b)
			n += n2

			if err != nil {
				err = errors.Wrap(err)
				return
			}

			n1, err = sw.WriteString("\"")
			n += int64(n1)

			if err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	if f.writeEtiketten && !didWriteBezeichnung {
		b := o.GetMetadatei().GetEtiketten()

		for _, v := range iter.SortedValues[kennung.Etikett](b) {
			n1, err = sw.WriteString(" ")
			n += int64(n1)

			if err != nil {
				err = errors.Wrap(err)
				return
			}

			n2, err = f.etikettenStringFormatWriter.WriteStringFormat(sw, &v)
			n += n2

			if err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	if m != checkout_mode.ModeObjekteOnly {
		n1, err = sw.WriteString("\n")
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		n2, err = f.rightAlignedWriter.WriteStringFormat(
			sw,
			"",
		)
		n += n2

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		n1, err = sw.WriteString(" ")
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		n2, err = f.fdStringFormatWriter.WriteStringFormat(
			sw,
			&fds.Akte,
		)
		n += n2

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	n1, err = sw.WriteString("]")
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
