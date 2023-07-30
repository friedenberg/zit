package objekte

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/checkout_mode"
	"github.com/friedenberg/zit/src/bravo/string_writer_format"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/echo/bezeichnung"
)

type CliOptions struct{}

type cli struct {
	options CliOptions

	writeTyp         bool
	writeBezeichnung bool
	writeEtiketten   bool

	rightAlignedWriter            schnittstellen.StringFormatWriter[string]
	shaStringFormatWriter         schnittstellen.StringFormatWriter[schnittstellen.ShaLike]
	kennungStringFormatWriter     schnittstellen.StringFormatWriter[kennung.KennungPtr]
	fdStringFormatWriter          schnittstellen.StringFormatWriter[*kennung.FD]
	typStringFormatWriter         schnittstellen.StringFormatWriter[*kennung.Typ]
	bezeichnungStringFormatWriter schnittstellen.StringFormatWriter[*bezeichnung.Bezeichnung]
	etikettenStringFormatWriter   schnittstellen.StringFormatWriter[kennung.EtikettSet]
}

func MakeCliFormat(
	options CliOptions,
	shaStringFormatWriter schnittstellen.StringFormatWriter[schnittstellen.ShaLike],
	fdStringFormatWriter schnittstellen.StringFormatWriter[*kennung.FD],
	kennungStringFormatWriter schnittstellen.StringFormatWriter[kennung.KennungPtr],
	typStringFormatWriter schnittstellen.StringFormatWriter[*kennung.Typ],
	bezeichnungStringFormatWriter schnittstellen.StringFormatWriter[*bezeichnung.Bezeichnung],
	etikettenStringFormatWriter schnittstellen.StringFormatWriter[kennung.EtikettSet],
) *cli {
	return &cli{
		options:                       options,
		writeTyp:                      true,
		writeBezeichnung:              true,
		writeEtiketten:                true,
		rightAlignedWriter:            string_writer_format.MakeRightAligned(),
		shaStringFormatWriter:         shaStringFormatWriter,
		kennungStringFormatWriter:     kennungStringFormatWriter,
		fdStringFormatWriter:          fdStringFormatWriter,
		typStringFormatWriter:         typStringFormatWriter,
		bezeichnungStringFormatWriter: bezeichnungStringFormatWriter,
		etikettenStringFormatWriter:   etikettenStringFormatWriter,
	}
}

func (f *cli) WriteStringFormat(
	sw io.StringWriter,
	colp CheckedOutLikePtr,
) (n int64, err error) {
	var (
		n1 int
		n2 int64
	)

	state := colp.GetState()

	n2, err = f.rightAlignedWriter.WriteStringFormat(
		sw,
		state.String(),
	)
	n += n2

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	o := colp.GetExternalLikePtr()
	fds := o.GetFDsPtr()
	n1, err = sw.WriteString("[")
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	var m checkout_mode.Mode

	if m, err = fds.GetCheckoutMode(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if m == checkout_mode.ModeAkteOnly {
		n2, err = f.kennungStringFormatWriter.WriteStringFormat(
			sw,
			o.GetKennungPtr(),
		)
		n += n2

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		n2, err = f.fdStringFormatWriter.WriteStringFormat(
			sw,
			&fds.Objekte,
		)
		n += n2

		if err != nil {
			err = errors.Wrap(err)
			return
		}
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
		t := o.GetMetadateiPtr().GetTypPtr()

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
		b := o.GetMetadateiPtr().GetBezeichnungPtr()

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
		b := o.GetMetadateiPtr().GetEtiketten()

		if b.Len() > 0 {
			n1, err = sw.WriteString(" ")
			n += int64(n1)

			if err != nil {
				err = errors.Wrap(err)
				return
			}

			n2, err = f.etikettenStringFormatWriter.WriteStringFormat(sw, b)
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
