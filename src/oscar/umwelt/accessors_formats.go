package umwelt

import (
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/echo/bezeichnung"
	"github.com/friedenberg/zit/src/echo/sha"
	"github.com/friedenberg/zit/src/etikett"
	"github.com/friedenberg/zit/src/foxtrot/hinweis"
	"github.com/friedenberg/zit/src/foxtrot/kennung"
	"github.com/friedenberg/zit/src/golf/fd"
	"github.com/friedenberg/zit/src/india/typ"
	"github.com/friedenberg/zit/src/india/zettel_external"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/mike/zettel_checked_out"
	"github.com/friedenberg/zit/src/november/store_fs"
)

func (u *Umwelt) FormatColorWriter() format.FuncColorWriter {
	if u.outIsTty {
		return format.MakeFormatWriterWithColor
	} else {
		return format.MakeFormatWriterNoopColor
	}
}

func (u *Umwelt) FormatSha() format.FormatWriterFunc[sha.Sha] {
	return sha.MakeCliFormat(
		u.FormatColorWriter(),
		u.StoreObjekten(),
	)
}

func (u *Umwelt) FormatHinweis() format.FormatWriterFunc[hinweis.Hinweis] {
	var a hinweis.Abbr

	if u.konfig.PrintAbbreviatedHinweisen {
		a = u.StoreObjekten()
	}

	return hinweis.MakeCliFormat(
		u.FormatColorWriter(),
		a,
		0,
		0,
	)
}

func (u *Umwelt) FormatBezeichnung() format.FormatWriterFunc[bezeichnung.Bezeichnung] {
	return bezeichnung.MakeCliFormat(u.FormatColorWriter())
}

func (u *Umwelt) FormatTyp() format.FormatWriterFunc[kennung.Typ] {
	return typ.MakeCliFormat(u.FormatColorWriter())
}

func (u *Umwelt) FormatEtikett() format.FormatWriterFunc[kennung.Etikett] {
	return etikett.MakeCliFormat(u.FormatColorWriter())
}

func (u *Umwelt) FormatTypTransacted(
	verb string,
) format.FormatWriterFunc[typ.Transacted] {
	return typ.MakeCliFormatTransacted(
		u.Standort(),
		u.FormatColorWriter(),
		u.FormatSha(),
		u.FormatTyp(),
		verb,
	)
}

func (u *Umwelt) FormatEtikettTransacted(
	verb string,
) format.FormatWriterFunc[etikett.Transacted] {
	return etikett.MakeCliFormatTransacted(
		u.Standort(),
		u.FormatColorWriter(),
		u.FormatSha(),
		u.FormatEtikett(),
		verb,
	)
}

func (u *Umwelt) FormatTypCheckedOut() format.FormatWriterFunc[typ.External] {
	return typ.MakeCliFormatExternal(
		u.Standort(),
		u.FormatColorWriter(),
		u.FormatSha(),
		u.FormatTyp(),
	)
}

func (u *Umwelt) FormatZettel() format.FormatWriterFunc[zettel.Objekte] {
	return zettel.MakeCliFormat(
		u.FormatBezeichnung(),
		format.MakeFormatStringer[kennung.EtikettSet](),
		u.FormatTyp(),
	)
}

func (u *Umwelt) FormatZettelExternal() format.FormatWriterFunc[zettel_external.Zettel] {
	return zettel_external.MakeCliFormatZettel(
		u.Standort(),
		u.FormatColorWriter(),
		u.FormatSha(),
		u.FormatZettel(),
	)
}

func (u *Umwelt) FormatZettelExternalAkte() format.FormatWriterFunc[zettel_external.Zettel] {
	return zettel_external.MakeCliFormatAkte(
		u.Standort(),
		u.FormatColorWriter(),
		u.FormatSha(),
	)
}

func (u *Umwelt) FormatZettelExternalFD() format.FormatWriterFunc[fd.FD] {
	return zettel_external.MakeCliFormatFD(
		u.Standort(),
		u.FormatColorWriter(),
	)
}

func (u *Umwelt) FormatZettelCheckedOut(
	mode zettel_checked_out.Mode,
) format.FormatWriterFunc[zettel_checked_out.Zettel] {
	return zettel_checked_out.MakeCliFormat(
		u.Standort(),
		u.FormatZettelExternal(),
		u.FormatZettelExternalAkte(),
		mode,
	)
}

func (u *Umwelt) FormatZettelCheckedOutFresh(
	mode zettel_checked_out.Mode,
) format.FormatWriterFunc[zettel_checked_out.Zettel] {
	return zettel_checked_out.MakeCliFormatFresh(
		u.Standort(),
		u.FormatZettelExternal(),
		u.FormatZettelExternalAkte(),
		mode,
	)
}

func (u *Umwelt) FormatZettelTransacted(verb string) format.FormatWriterFunc[zettel.Transacted] {
	return zettel.MakeCliFormatTransacted(
		u.FormatHinweis(),
		u.FormatSha(),
		u.FormatZettel(),
		verb,
	)
}

func (u *Umwelt) FormatFileNotRecognized() format.FormatWriterFunc[fd.FD] {
	return store_fs.MakeCliFormatNotRecognized(
		u.FormatColorWriter(),
		u.Standort(),
		u.FormatSha(),
	)
}

func (u *Umwelt) FormatFileRecognized() format.FormatWriterFunc[store_fs.FileRecognized] {
	return store_fs.MakeCliFormatRecognized(
		u.FormatColorWriter(),
		u.Standort(),
		u.FormatSha(),
		u.FormatZettel(),
	)
}

func (u *Umwelt) FormatDirDeleted() format.FormatWriterFunc[store_fs.Dir] {
	return store_fs.MakeCliFormatDirDeleted(
		u.FormatColorWriter(),
		u.Standort(),
	)
}

func (u *Umwelt) FormatFDDeleted() format.FormatWriterFunc[fd.FD] {
	return store_fs.MakeCliFormatFDDeleted(
		u.FormatColorWriter(),
		u.Standort(),
		u.FormatZettelExternalFD(),
	)
}
