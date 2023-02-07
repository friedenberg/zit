package umwelt

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/echo/bezeichnung"
	"github.com/friedenberg/zit/src/echo/sha_cli_format"
	"github.com/friedenberg/zit/src/foxtrot/fd"
	"github.com/friedenberg/zit/src/hotel/etikett"
	"github.com/friedenberg/zit/src/hotel/typ"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/kasten"
	"github.com/friedenberg/zit/src/kilo/zettel_external"
	"github.com/friedenberg/zit/src/lima/zettel_checked_out"
	"github.com/friedenberg/zit/src/mike/store_fs"
)

func (u *Umwelt) FormatColorWriter() format.FuncColorWriter {
	if u.outIsTty {
		return format.MakeFormatWriterWithColor
	} else {
		return format.MakeFormatWriterNoopColor
	}
}

func (u *Umwelt) FormatSha(
	a schnittstellen.FuncAbbreviateValue,
) schnittstellen.FuncWriterFormat[schnittstellen.Sha] {
	return sha_cli_format.MakeCliFormat(u.FormatColorWriter(), a)
}

func (u *Umwelt) FormatHinweis() schnittstellen.FuncWriterFormat[kennung.Hinweis] {
	var a schnittstellen.FuncAbbreviateKorper

	if u.konfig.PrintAbbreviatedHinweisen {
		a = u.StoreObjekten().GetAbbrStore().AbbreviateHinweis
	}

	return kennung.MakeHinweisCliFormat(
		u.FormatColorWriter(),
		a,
		0,
		0,
	)
}

func (u *Umwelt) FormatBezeichnung() schnittstellen.FuncWriterFormat[bezeichnung.Bezeichnung] {
	return bezeichnung.MakeCliFormat(u.FormatColorWriter())
}

func (u *Umwelt) FormatTyp() schnittstellen.FuncWriterFormat[kennung.Typ] {
	return typ.MakeCliFormat(u.FormatColorWriter())
}

func (u *Umwelt) FormatEtikett() schnittstellen.FuncWriterFormat[kennung.Etikett] {
	return etikett.MakeCliFormat(u.FormatColorWriter())
}

func (u *Umwelt) FormatKasten() schnittstellen.FuncWriterFormat[kennung.Kasten] {
	return kasten.MakeCliFormat(u.FormatColorWriter())
}

func (u *Umwelt) FormatTypTransacted(
	verb string,
) schnittstellen.FuncWriterFormat[typ.Transacted] {
	return typ.MakeCliFormatTransacted(
		u.Standort(),
		u.FormatColorWriter(),
		u.FormatSha(u.StoreObjekten().GetAbbrStore().AbbreviateSha),
		u.FormatTyp(),
		verb,
	)
}

func (u *Umwelt) FormatEtikettTransacted(
	verb string,
) schnittstellen.FuncWriterFormat[etikett.Transacted] {
	return etikett.MakeCliFormatTransacted(
		u.Standort(),
		u.FormatColorWriter(),
		u.FormatSha(u.StoreObjekten().GetAbbrStore().AbbreviateSha),
		u.FormatEtikett(),
		verb,
	)
}

func (u *Umwelt) FormatKastenTransacted(
	verb string,
) schnittstellen.FuncWriterFormat[kasten.Transacted] {
	return kasten.MakeCliFormatTransacted(
		u.Standort(),
		u.FormatColorWriter(),
		u.FormatSha(u.StoreObjekten().GetAbbrStore().AbbreviateSha),
		u.FormatKasten(),
		verb,
	)
}

func (u *Umwelt) FormatTypCheckedOut() schnittstellen.FuncWriterFormat[typ.External] {
	return typ.MakeCliFormatExternal(
		u.Standort(),
		u.FormatColorWriter(),
		u.FormatSha(u.StoreObjekten().GetAbbrStore().AbbreviateSha),
		u.FormatTyp(),
	)
}

func (u *Umwelt) FormatZettel() schnittstellen.FuncWriterFormat[zettel.Objekte] {
	return zettel.MakeCliFormat(
		u.FormatBezeichnung(),
		format.MakeFormatStringer[kennung.EtikettSet](),
		u.FormatTyp(),
	)
}

func (u *Umwelt) FormatZettelExternal() schnittstellen.FuncWriterFormat[zettel_external.Zettel] {
	return zettel_external.MakeCliFormat(
		u.Standort(),
		u.FormatColorWriter(),
		u.FormatHinweis(),
		u.FormatSha(u.StoreObjekten().GetAbbrStore().AbbreviateSha),
		u.FormatZettel(),
	)
}

func (u *Umwelt) FormatZettelExternalFD() schnittstellen.FuncWriterFormat[fd.FD] {
	return zettel_external.MakeCliFormatFD(
		u.Standort(),
		u.FormatColorWriter(),
	)
}

func (u *Umwelt) FormatZettelCheckedOut() schnittstellen.FuncWriterFormat[zettel_checked_out.Zettel] {
	return zettel_checked_out.MakeCliFormat(
		u.Standort(),
		u.FormatColorWriter(),
		u.FormatHinweis(),
		u.FormatSha(u.StoreObjekten().GetAbbrStore().AbbreviateSha),
		u.FormatZettel(),
	)
}

func (u *Umwelt) FormatZettelTransacted() schnittstellen.FuncWriterFormat[zettel.Transacted] {
	return zettel.MakeCliFormatTransacted(
		u.FormatHinweis(),
		u.FormatSha(u.StoreObjekten().GetAbbrStore().AbbreviateSha),
		u.FormatZettel(),
	)
}

func (u *Umwelt) FormatZettelTransactedDelta(verb string) schnittstellen.FuncWriterFormat[zettel.Transacted] {
	return zettel.MakeCliFormatTransactedDelta(
		verb,
		u.FormatZettelTransacted(),
	)
}

func (u *Umwelt) FormatFileNotRecognized() schnittstellen.FuncWriterFormat[fd.FD] {
	return store_fs.MakeCliFormatNotRecognized(
		u.FormatColorWriter(),
		u.Standort(),
		u.FormatSha(u.StoreObjekten().GetAbbrStore().AbbreviateSha),
	)
}

func (u *Umwelt) FormatFileRecognized() schnittstellen.FuncWriterFormat[store_fs.FileRecognized] {
	return store_fs.MakeCliFormatRecognized(
		u.FormatColorWriter(),
		u.Standort(),
		u.FormatSha(u.StoreObjekten().GetAbbrStore().AbbreviateSha),
		u.FormatZettel(),
	)
}

func (u *Umwelt) FormatDirDeleted() schnittstellen.FuncWriterFormat[store_fs.Dir] {
	return store_fs.MakeCliFormatDirDeleted(
		u.FormatColorWriter(),
		u.Standort(),
	)
}

func (u *Umwelt) FormatFDDeleted() schnittstellen.FuncWriterFormat[fd.FD] {
	return store_fs.MakeCliFormatFDDeleted(
		u.FormatColorWriter(),
		u.Standort(),
		u.FormatZettelExternalFD(),
	)
}
