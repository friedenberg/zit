package metadatei

import (
	"io"
	"path"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/bravo/ohio"
	"github.com/friedenberg/zit/src/bravo/script_config"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/delta/kennung"
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
	m.Reset()

	etiketten := kennung.MakeEtikettMutableSet()

	var n1 int64

	defer func() {
		m.Etiketten = etiketten.CloneSetPtrLike()
		c.SetMetadatei(m)
		c.SetAkteSha(m.AkteSha)
	}()

	var akteFD kennung.FD

	lr := format.MakeLineReaderConsumeEmpty(
		ohio.MakeLineReaderIterate(
			ohio.MakeLineReaderKeyValues(
				map[string]schnittstellen.FuncSetString{
					"#": m.Bezeichnung.Set,
					"%": ohio.MakeLineReaderNop(),
					"-": iter.MakeFuncSetString[
						kennung.Etikett,
						*kennung.Etikett,
					](etiketten),
					"!": func(v string) (err error) {
						return f.readTyp(&m, v, &akteFD)
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

	if !m.AkteSha.IsNull() && !akteFD.Sha.IsNull() {
		err = errors.Wrap(ErrHasInlineAkteAndFilePath{
			AkteFD:    akteFD,
			InlineSha: inlineAkteSha,
		})

		return
	} else if !akteFD.Sha.IsNull() {
		if afs, ok := c.(kennung.AkteFDSetter); ok {
			afs.SetAkteFD(akteFD)
		}

		m.AkteSha = akteFD.Sha
	}

	switch {
	case m.AkteSha.IsNull() && !inlineAkteSha.IsNull():
		m.AkteSha = inlineAkteSha

	case !m.AkteSha.IsNull() && inlineAkteSha.IsNull():
		// noop

	case !m.AkteSha.IsNull() && !inlineAkteSha.IsNull() &&
		!m.AkteSha.Equals(inlineAkteSha):
		err = errors.Wrap(ErrHasInlineAkteAndMetadateiSha{
			InlineSha:    inlineAkteSha,
			MetadateiSha: m.AkteSha,
		})

		return
	}

	return
}

func (f textParser) readTyp(
	m *Metadatei,
	desc string,
	akteFD *kennung.FD,
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

		if *akteFD, err = kennung.FDFromPathWithAkteWriterFactory(
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
	var sh sha.Sha

	if err = sh.Set(maybeSha); err != nil {
		err = errors.Wrap(err)
		return
	}

	m.AkteSha = sh

	return
}
