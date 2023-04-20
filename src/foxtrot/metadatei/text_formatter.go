package metadatei

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/ohio"
	"github.com/friedenberg/zit/src/bravo/script_config"
)

type textFormatter struct {
	textFormatterCommon
	sequence []schnittstellen.FuncWriterElementInterface[TextFormatterContext]
	// standort      standort.Standort
	// InlineChecker kennung.InlineTypChecker
	// akteFactory schnittstellen.AkteReaderFactory
	// AkteFormatter              script_config.RemoteScript
	// ExcludeMetadatei           bool
}

func MakeTextFormatterMetadateiOnly(
	akteFactory schnittstellen.AkteReaderFactory,
	akteFormatter script_config.RemoteScript,
) textFormatter {
	if akteFactory == nil {
		panic("akte reader factory is nil")
	}

	common := textFormatterCommon{
		akteFactory:   akteFactory,
		akteFormatter: akteFormatter,
	}

	return textFormatter{
		textFormatterCommon: common,
		sequence: []schnittstellen.FuncWriterElementInterface[TextFormatterContext]{
			common.writeBoundary,
			common.writeCommonMetadateiFormat,
			common.writeShaTyp,
			common.writeBoundary,
		},
	}
}

func MakeTextFormatterMetadateiInlineAkte(
	akteFactory schnittstellen.AkteReaderFactory,
	akteFormatter script_config.RemoteScript,
) textFormatter {
	if akteFactory == nil {
		panic("akte reader factory is nil")
	}

	common := textFormatterCommon{
		akteFactory:   akteFactory,
		akteFormatter: akteFormatter,
	}

	return textFormatter{
		textFormatterCommon: common,
		sequence: []schnittstellen.FuncWriterElementInterface[TextFormatterContext]{
			common.writeBoundary,
			common.writeCommonMetadateiFormat,
			common.writeTyp,
			common.writeBoundary,
			common.writeNewLine,
			common.writeAkte,
		},
	}
}

func MakeTextFormatterExcludeMetadatei(
	akteFactory schnittstellen.AkteReaderFactory,
	akteFormatter script_config.RemoteScript,
) textFormatter {
	if akteFactory == nil {
		panic("akte reader factory is nil")
	}

	common := textFormatterCommon{
		akteFactory:   akteFactory,
		akteFormatter: akteFormatter,
	}

	return textFormatter{
		textFormatterCommon: common,
		sequence: []schnittstellen.FuncWriterElementInterface[TextFormatterContext]{
			common.writeAkte,
		},
	}
}

// func MakeTextFormatterExcludeMetadatei(
// 	standort standort.Standort,
// 	inlineChecker kennung.InlineTypChecker,
// 	akteFactory schnittstellen.AkteIOFactory,
// 	akteFormatter script_config.RemoteScript,
// ) textFormatter {
// 	return textFormatter{
// 		standort:         standort,
// 		InlineChecker:    inlineChecker,
// 		AkteFactory:      akteFactory,
// 		AkteFormatter:    akteFormatter,
// 		ExcludeMetadatei: true,
// 	}
// }

// func MakeTextFormatterIncludeAkte(
// 	standort standort.Standort,
// 	inlineChecker kennung.InlineTypChecker,
// 	akteFactory schnittstellen.AkteIOFactory,
// 	akteFormatter script_config.RemoteScript,
// ) textFormatter {
// 	return textFormatter{
// 		standort:      standort,
// 		InlineChecker: inlineChecker,
// 		AkteFactory:   akteFactory,
// 		AkteFormatter: akteFormatter,
// 	}
// }

// func MakeTextFormatterAkteShaOnly(
// 	standort standort.Standort,
// 	akteFactory schnittstellen.AkteIOFactory,
// 	akteFormatter script_config.RemoteScript,
// ) textFormatter {
// 	return textFormatter{
// 		standort:      standort,
// 		AkteFactory:   akteFactory,
// 		AkteFormatter: akteFormatter,
// 	}
// }

func (f textFormatter) Format(
	w io.Writer,
	c TextFormatterContext,
) (n int64, err error) {
	return ohio.WriteSeq(w, c, f.sequence...)
}

// func (f textFormatter) writeMetadatei(
// 	w1 io.Writer,
// 	c TextFormatterContext,
// ) (n int64, err error) {
// 	w := format.NewLineWriter()
// 	m := c.GetMetadatei()

// 	if m.Bezeichnung.String() != "" || !f.doNotWriteEmptyBezeichnung {
// 		w.WriteLines(
// 			fmt.Sprintf("# %s", m.Bezeichnung),
// 		)
// 	}

// 	for _, e := range collections.SortedValues(m.Etiketten) {
// 		if e.IsEmpty() {
// 			continue
// 		}

// 		w.WriteFormat("- %s", e)
// 	}

// 	// ap := c.GetAktePath()

// 	switch {
// 	// case ap != "":
// 	// 	w.WriteLines(
// 	// 		fmt.Sprintf("! %s", ap),
// 	// 	)

// 	// case includeAkteSha:
// 	// 	sh := c.GetAkteSha()

// 	// 	w.WriteLines(
// 	// 		fmt.Sprintf("! %s.%s", sh, m.Typ),
// 	// 	)

// 	default:
// 		w.WriteLines(
// 			fmt.Sprintf("! %s", m.Typ),
// 		)
// 	}

// 	if n, err = w.WriteTo(w1); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	return
// }

// func (f textFormatter) writeToExternalAkte(
// 	w1 io.Writer,
// 	c TextFormatterContext,
// ) (n int64, err error) {
// 	m := c.GetMetadatei()
// 	w := format.NewLineWriter()

// 	w.WriteLines(
// 		Boundary,
// 		fmt.Sprintf("# %s", m.Bezeichnung),
// 	)

// 	for _, e := range collections.SortedValues(m.Etiketten) {
// 		w.WriteFormat("- %s", e)
// 	}

// 	ap := c.GetAktePath()

// 	if strings.Index(ap, "\n") != -1 {
// 		panic(errors.Errorf("ExternalAktePath contains newline: %q", ap))
// 	}

// 	w.WriteLines(
// 		fmt.Sprintf("! %s", ap),
// 	)

// 	w.WriteLines(
// 		Boundary,
// 	)

// 	n, err = w.WriteTo(w1)

// 	if err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	var ar io.ReadCloser

// 	if ar, err = f.akteFactory.AkteReader(c.GetAkteSha()); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	if ar == nil {
// 		err = errors.Errorf("akte reader is nil")
// 		return
// 	}

// 	defer errors.Deferred(&err, ar.Close)

// 	var file *os.File

// 	if file, err = files.Create(c.GetAktePath()); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	defer errors.Deferred(&err, file.Close)

// 	var n1 int64

// 	n1, err = io.Copy(file, ar)
// 	n += n1

// 	if err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	return
// }
