package sku_fmt

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/erworben_cli_print_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type cliCheckedOutLike struct {
	options erworben_cli_print_options.PrintOptions

	rightAlignedWriter          schnittstellen.StringFormatWriter[string]
	shaStringFormatWriter       schnittstellen.StringFormatWriter[schnittstellen.ShaLike]
	kennungStringFormatWriter   schnittstellen.StringFormatWriter[*kennung.Kennung2]
	fdStringFormatWriter        schnittstellen.StringFormatWriter[*fd.FD]
	metadateiStringFormatWriter schnittstellen.StringFormatWriter[*metadatei.Metadatei]
}

func MakeCliCheckedOutLikeFormat(
	options erworben_cli_print_options.PrintOptions,
	shaStringFormatWriter schnittstellen.StringFormatWriter[schnittstellen.ShaLike],
	fdStringFormatWriter schnittstellen.StringFormatWriter[*fd.FD],
	kennungStringFormatWriter schnittstellen.StringFormatWriter[*kennung.Kennung2],
	metadateiStringFormatWriter schnittstellen.StringFormatWriter[*metadatei.Metadatei],
) *cliCheckedOutLike {
	return &cliCheckedOutLike{
		options:                     options,
		rightAlignedWriter:          string_format_writer.MakeRightAligned(),
		shaStringFormatWriter:       shaStringFormatWriter,
		kennungStringFormatWriter:   kennungStringFormatWriter,
		fdStringFormatWriter:        fdStringFormatWriter,
		metadateiStringFormatWriter: metadateiStringFormatWriter,
	}
}

func (f *cliCheckedOutLike) WriteStringFormat(
	sw schnittstellen.WriterAndStringWriter,
	co sku.CheckedOutLike,
) (n int64, err error) {
	var (
		n1    int
		n2    int64
		state = co.GetState()
	)

	{
		var stateString string

		if state == checked_out_state.StateError {
			stateString = co.GetError().Error()
		} else {
			stateString = state.String()
		}

		n2, err = f.rightAlignedWriter.WriteStringFormat(sw, stateString)
		n += n2

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	n1, err = sw.WriteString("[")
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	cofs, ok := co.(*CheckedOut)

	if !ok {
		n1, err = fmt.Fprintf(sw, "unsupported check out type: %T", co)
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}
		n1, err = sw.WriteString("]")
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	o := &cofs.External
	fds := o.GetFDsPtr()
	m := fds.GetCheckoutMode()

	switch {
	case cofs.IsImport:
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

	if state == checked_out_state.StateConflicted {
		n1, err = sw.WriteString("]")
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	n2, err = f.metadateiStringFormatWriter.WriteStringFormat(sw, o.GetMetadatei())
	n += n2

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	if m != checkout_mode.ModeObjekteOnly && m != checkout_mode.ModeNone {
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
