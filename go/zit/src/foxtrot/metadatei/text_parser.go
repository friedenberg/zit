package metadatei

import (
	"io"
	"path"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/script_config"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/ohio"
	"github.com/friedenberg/zit/src/echo/fd"
	"github.com/friedenberg/zit/src/echo/format"
	"github.com/friedenberg/zit/src/echo/kennung"
)

type textParser struct {
	awf schnittstellen.AkteWriterFactory
	af  script_config.RemoteScript
}

func MakeTextParser(
	awf schnittstellen.AkteWriterFactory,
	akteFormatter script_config.RemoteScript,
) TextParser {
	if awf == nil {
		panic("nil AkteWriterFactory")
	}

	return textParser{
		awf: awf,
		af:  akteFormatter,
	}
}

func (f textParser) ParseMetadatei(
	r io.Reader,
	c TextParserContext,
) (n int64, err error) {
	m := c.GetMetadatei()
	Resetter.Reset(m)

	var n1 int64

	defer func() {
		c.SetAkteSha(&m.Akte)
	}()

	var akteFD fd.FD

	lr := format.MakeLineReaderConsumeEmpty(
		ohio.MakeLineReaderIterate(
			ohio.MakeLineReaderKeyValues(
				map[string]schnittstellen.FuncSetString{
					"#": m.Bezeichnung.Set,
					"%": func(v string) (err error) {
						m.Comments = append(m.Comments, v)
						return
					},
					"-": iter.MakeFuncSetString[
						kennung.Etikett,
						*kennung.Etikett,
					](m.GetEtikettenMutable()),
					"!": func(v string) (err error) {
						return f.readTyp(m, v, &akteFD)
					},
				},
			),
		),
	)

	var akteWriter sha.WriteCloser

	if akteWriter, err = f.awf.AkteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if akteWriter == nil {
		err = errors.Errorf("akte writer is nil")
		return
	}

	defer errors.DeferredCloser(&err, akteWriter)

	mr := Reader{
		// RequireMetadatei: true,
		Metadatei: lr,
		Akte:      akteWriter,
	}

	// if cmg, ok := c.(checkout_mode.Getter); ok {
	// 	var checkoutMode checkout_mode.Mode

	// 	if checkoutMode, err = cmg.GetCheckoutMode(); err != nil {
	// 		err = errors.Wrap(err)
	// 		return
	// 	}

	// 	mr.RequireMetadatei = checkoutMode.IncludesObjekte()
	// }

	if n, err = mr.ReadFrom(r); err != nil {
		n += n1
		err = errors.Wrap(err)
		return
	}

	n += n1

	inlineAkteSha := sha.Make(akteWriter.GetShaLike())

	if !m.Akte.IsNull() && !akteFD.GetShaLike().IsNull() {
		err = errors.Wrap(
			MakeErrHasInlineAkteAndFilePath(
				&akteFD,
				inlineAkteSha,
			),
		)

		return
	} else if !akteFD.GetShaLike().IsNull() {
		if afs, ok := c.(fd.AkteFDSetter); ok {
			afs.SetAkteFD(&akteFD)
		}

		m.Akte.SetShaLike(akteFD.GetShaLike())
	}

	switch {
	case m.Akte.IsNull() && !inlineAkteSha.IsNull():
		m.Akte.SetShaLike(inlineAkteSha)

	case !m.Akte.IsNull() && inlineAkteSha.IsNull():
		// noop

	case !m.Akte.IsNull() && !inlineAkteSha.IsNull() &&
		!m.Akte.Equals(inlineAkteSha):
		err = errors.Wrap(
			MakeErrHasInlineAkteAndMetadateiSha(
				inlineAkteSha,
				&m.Akte,
			),
		)

		return
	}

	return
}

func (f textParser) readTyp(
	m *Metadatei,
	desc string,
	akteFD *fd.FD,
) (err error) {
	if desc == "" {
		return
	}

	tail := path.Ext(desc)
	head := desc[:len(desc)-len(tail)]

	//! <path>.<typ ext>
	switch {
	case files.Exists(desc):
		if err = m.Typ.Set(tail); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = akteFD.SetWithAkteWriterFactory(
			desc,
			f.awf,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

	//! <sha>.<typ ext>
	case tail != "":
		if err = f.setAkteSha(m, head); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = m.Typ.Set(tail); err != nil {
			err = errors.Wrap(err)
			return
		}

	//! <sha>
	case tail == "":
		if err = f.setAkteSha(m, head); err == nil {
			return
		}

		err = nil

		fallthrough

	//! <typ ext>
	default:
		if err = m.Typ.Set(head); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (f textParser) setAkteSha(
	m *Metadatei,
	maybeSha string,
) (err error) {
	if err = m.Akte.Set(maybeSha); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}