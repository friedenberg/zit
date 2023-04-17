package metadatei

import (
	"io"
	"os"
	"path"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/script_config"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type ParserContext interface {
	GetMetadateiPtr() *Metadatei
	SetAkteFD(kennung.FD) error
}

type TextParser interface {
	Parse(io.Reader, ParserContext) (int64, error)
}

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

func (f textParser) Parse(
	r io.Reader,
	c ParserContext,
) (n int64, err error) {
	m := c.GetMetadateiPtr()
	etiketten := kennung.MakeEtikettMutableSet()

	var n1 int64

	defer func() {
		m.Etiketten = etiketten.ImmutableClone()
	}()

	lr := format.MakeLineReaderConsumeEmpty(
		format.MakeLineReaderIterate(
			format.MakeLineReaderKeyValues(
				map[string]schnittstellen.FuncSetString{
					"#": m.Bezeichnung.Set,
					"%": format.MakeLineReaderNop(),
					"-": collections.MakeFuncSetString[
						kennung.Etikett,
						*kennung.Etikett,
					](etiketten),
					"!": func(v string) (err error) {
						return f.readTyp(c, v)
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
		RequireMetadatei: true,
		Metadatei:        lr,
		Akte:             akteWriter,
	}

	if n, err = mr.ReadFrom(r); err != nil {
		n += n1
		err = errors.Wrap(err)
		return
	}

	n += n1

	inlineAkteSha := sha.Make(akteWriter.Sha())

	switch {
	case m.AkteSha.IsNull() && !inlineAkteSha.IsNull():
		m.AkteSha = inlineAkteSha

	case !m.AkteSha.IsNull() && inlineAkteSha.IsNull():
		// noop

	case !m.AkteSha.IsNull() && !inlineAkteSha.IsNull():
		err = ErrHasInlineAkteAndFilePath{
			Metadatei: *m,
		}

		return
	}

	return
}

func (tp textParser) readExternalAkte(
	p string,
) (sh sha.Sha, err error) {
	var f *os.File

	if f, err = files.OpenExclusiveReadOnly(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, f)

	var akteWriter sha.WriteCloser

	if akteWriter, err = tp.awf.AkteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if akteWriter == nil {
		err = errors.Errorf("akte writer is nil")
		return
	}

	defer errors.DeferredCloser(&err, akteWriter)

	if _, err = io.Copy(akteWriter, f); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh = sha.Make(akteWriter.Sha())

	return
}

func (f textParser) readTyp(
	c ParserContext,
	desc string,
) (err error) {
	m := c.GetMetadateiPtr()

	if desc == "" {
		return
	}

	tail := path.Ext(desc)
	head := desc[:len(desc)-len(tail)]

	errors.TodoP3("handle akte descs that are invalid files")
	//! <path>.<typ ext>
	switch {
	case files.Exists(desc):
		if err = m.Typ.Set(tail); err != nil {
			err = errors.Wrap(err)
			return
		}

		var externalAkteSha sha.Sha

		if externalAkteSha, err = f.readExternalAkte(desc); err != nil {
			err = errors.Wrap(err)
			return
		}

		c.GetMetadateiPtr().AkteSha = externalAkteSha

	//! <sha>.<typ ext>
	case tail != "":
		if err = f.setAkteSha(c, head); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = m.Typ.Set(tail); err != nil {
			err = errors.Wrap(err)
			return
		}

	//! <sha>
	case tail == "":
		if err = f.setAkteSha(c, head); err == nil {
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
	c ParserContext,
	maybeSha string,
) (err error) {
	var sh sha.Sha

	if err = sh.Set(maybeSha); err != nil {
		err = errors.Wrap(err)
		return
	}

	c.GetMetadateiPtr().AkteSha = sh

	return
}
