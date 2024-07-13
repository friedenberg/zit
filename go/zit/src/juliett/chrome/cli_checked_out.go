package chrome

import (
	"net/url"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/erworben_cli_print_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/bezeichnung"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type cliCheckedOut struct {
	options erworben_cli_print_options.PrintOptions

	rightAlignedWriter          interfaces.StringFormatWriter[string]
	shaStringFormatWriter       interfaces.StringFormatWriter[interfaces.ShaLike]
	kennungStringFormatWriter   interfaces.StringFormatWriter[*kennung.ObjectId]
	metadateiStringFormatWriter interfaces.StringFormatWriter[*metadatei.Metadatei]

	typStringFormatWriter         interfaces.StringFormatWriter[*kennung.Type]
	bezeichnungStringFormatWriter interfaces.StringFormatWriter[*bezeichnung.Bezeichnung]
	etikettenStringFormatWriter   interfaces.StringFormatWriter[*kennung.Tag]
}

func MakeCliCheckedOutFormat(
	options erworben_cli_print_options.PrintOptions,
	shaStringFormatWriter interfaces.StringFormatWriter[interfaces.ShaLike],
	kennungStringFormatWriter interfaces.StringFormatWriter[*kennung.ObjectId],
	metadateiStringFormatWriter interfaces.StringFormatWriter[*metadatei.Metadatei],
	typStringFormatWriter interfaces.StringFormatWriter[*kennung.Type],
	bezeichnungStringFormatWriter interfaces.StringFormatWriter[*bezeichnung.Bezeichnung],
	etikettenStringFormatWriter interfaces.StringFormatWriter[*kennung.Tag],
) *cliCheckedOut {
	return &cliCheckedOut{
		options:                       options,
		rightAlignedWriter:            string_format_writer.MakeRightAligned(),
		shaStringFormatWriter:         shaStringFormatWriter,
		kennungStringFormatWriter:     kennungStringFormatWriter,
		metadateiStringFormatWriter:   metadateiStringFormatWriter,
		typStringFormatWriter:         typStringFormatWriter,
		bezeichnungStringFormatWriter: bezeichnungStringFormatWriter,
		etikettenStringFormatWriter:   etikettenStringFormatWriter,
	}
}

func (f *cliCheckedOut) WriteStringFormat(
	sw interfaces.WriterAndStringWriter,
	col sku.CheckedOutLike,
) (n int64, err error) {
	co := col.(*CheckedOut)
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

	n1, err = sw.WriteString("[")

	if co.State != checked_out_state.StateUntracked {
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		n2, err = f.kennungStringFormatWriter.WriteStringFormat(
			sw,
			&co.Internal.Kennung,
		)
		n += n2

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		n2, err = f.metadateiStringFormatWriter.WriteStringFormat(
			sw,
			&co.Internal.Metadatei,
		)
		n += n2

		if err != nil {
			err = errors.Wrap(err)
			return
		}

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
	}

	item := co.External.item
	browser := &co.External.browser

	{
		n1, err = sw.WriteString("!")
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		n2, err = f.typStringFormatWriter.WriteStringFormat(
			sw,
			&browser.Metadatei.Typ,
		)
		n += n2

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		if !browser.Metadatei.Bezeichnung.IsEmpty() {
			n1, err = sw.WriteString(" ")
			n += int64(n1)

			if err != nil {
				err = errors.Wrap(err)
				return
			}

			n2, err = f.bezeichnungStringFormatWriter.WriteStringFormat(
				sw,
				&browser.Metadatei.Bezeichnung,
			)
			n += n2

			if err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	{
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

		var u *url.URL

		if u, err = item.GetUrl(); err != nil {
			err = errors.Wrap(err)
			return
		}

		n1, err = sw.WriteString(u.String())
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}

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

		for _, v := range iter.SortedValues(browser.Metadatei.GetEtiketten()) {
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

	n1, err = sw.WriteString("]")
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
