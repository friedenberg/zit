package umwelt

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/echo/bezeichnung"
	"github.com/friedenberg/zit/src/echo/sha_cli_format"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/hotel/etikett"
	"github.com/friedenberg/zit/src/hotel/kasten"
	"github.com/friedenberg/zit/src/hotel/typ"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/kilo/zettel_external"
	"github.com/friedenberg/zit/src/mike/store_fs"
)

//                                                    _
//    ___ ___  _ __ ___  _ __   ___  _ __   ___ _ __ | |_ ___
//   / __/ _ \| '_ ` _ \| '_ \ / _ \| '_ \ / _ \ '_ \| __/ __|
//  | (_| (_) | | | | | | |_) | (_) | | | |  __/ | | | |_\__ \
//   \___\___/|_| |_| |_| .__/ \___/|_| |_|\___|_| |_|\__|___/
//                      |_|

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

func (u *Umwelt) FormatKasten() schnittstellen.FuncWriterFormat[kennung.Kasten] {
	return kasten.MakeCliFormat(u.FormatColorWriter())
}

func (u *Umwelt) FormatEtikett() schnittstellen.FuncWriterFormat[kennung.Etikett] {
	return etikett.MakeCliFormat(u.FormatColorWriter())
}

//    ___  _     _      _    _
//   / _ \| |__ (_) ___| | _| |_ ___ _ __
//  | | | | '_ \| |/ _ \ |/ / __/ _ \ '_ \
//  | |_| | |_) | |  __/   <| ||  __/ | | |
//   \___/|_.__// |\___|_|\_\\__\___|_| |_|
//            |__/

func (u *Umwelt) FormatEtikettTransacted() schnittstellen.FuncWriterFormat[etikett.Transacted] {
	return etikett.MakeCliFormatTransacted(
		u.Standort(),
		u.FormatColorWriter(),
		u.FormatSha(u.StoreObjekten().GetAbbrStore().AbbreviateSha),
		u.FormatEtikett(),
	)
}

func (u *Umwelt) FormatTypTransacted() schnittstellen.FuncWriterFormat[typ.Transacted] {
	return typ.MakeCliFormatTransacted(
		u.Standort(),
		u.FormatColorWriter(),
		u.FormatSha(u.StoreObjekten().GetAbbrStore().AbbreviateSha),
		u.FormatTyp(),
	)
}

func (u *Umwelt) FormatKastenTransacted() schnittstellen.FuncWriterFormat[kasten.Transacted] {
	return kasten.MakeCliFormatTransacted(
		u.Standort(),
		u.FormatColorWriter(),
		u.FormatSha(u.StoreObjekten().GetAbbrStore().AbbreviateSha),
		u.FormatKasten(),
	)
}

func (u *Umwelt) FormatTypCheckedOut() schnittstellen.FuncWriterFormat[typ.CheckedOut] {
	return typ.MakeCliFormatCheckedOut(
		u.Standort(),
		u.FormatColorWriter(),
		u.FormatSha(u.StoreObjekten().GetAbbrStore().AbbreviateSha),
		u.FormatTyp(),
	)
}

func (u *Umwelt) FormatKastenCheckedOut() schnittstellen.FuncWriterFormat[kasten.CheckedOut] {
	return kasten.MakeCliFormatCheckedOut(
		u.Standort(),
		u.FormatColorWriter(),
		u.FormatSha(u.StoreObjekten().GetAbbrStore().AbbreviateSha),
		u.FormatKasten(),
	)
}

func (u *Umwelt) FormatEtikettCheckedOut() schnittstellen.FuncWriterFormat[etikett.CheckedOut] {
	return etikett.MakeCliFormatCheckedOut(
		u.Standort(),
		u.FormatColorWriter(),
		u.FormatSha(u.StoreObjekten().GetAbbrStore().AbbreviateSha),
		u.FormatEtikett(),
	)
}

func (u *Umwelt) FormatMetadateiGattung(
	g schnittstellen.GattungGetter,
) schnittstellen.FuncWriterFormat[metadatei.Metadatei] {
	return metadatei.MakeCliFormat(
		u.FormatBezeichnung(),
		format.MakeFormatStringer[kennung.Etikett](
			collections.StringCommaSeparated[kennung.Etikett],
		),
		u.FormatTyp(),
	)
}

func (u *Umwelt) FormatZettelExternal() schnittstellen.FuncWriterFormat[zettel.External] {
	return zettel_external.MakeCliFormat(
		u.Standort(),
		u.FormatColorWriter(),
		u.FormatHinweis(),
		u.FormatSha(u.StoreObjekten().GetAbbrStore().AbbreviateSha),
		u.FormatMetadateiGattung(gattung.Zettel),
	)
}

func (u *Umwelt) FormatExternalFD() schnittstellen.FuncWriterFormat[kennung.FD] {
	return zettel_external.MakeCliFormatFD(
		u.Standort(),
		u.FormatColorWriter(),
	)
}

func (u *Umwelt) FormatZettelCheckedOut() schnittstellen.FuncWriterFormat[zettel.CheckedOut] {
	return zettel.MakeCliFormatCheckedOut(
		u.Standort(),
		u.FormatColorWriter(),
		u.FormatHinweis(),
		u.FormatSha(u.StoreObjekten().GetAbbrStore().AbbreviateSha),
		// u.FormatTyp(),
		u.FormatMetadateiGattung(gattung.Zettel),
	)
}

func (u *Umwelt) FormatZettelTransacted() schnittstellen.FuncWriterFormat[zettel.Transacted] {
	return zettel.MakeCliFormatTransacted(
		u.FormatHinweis(),
		u.FormatSha(u.StoreObjekten().GetAbbrStore().AbbreviateSha),
		u.FormatMetadateiGattung(gattung.Zettel),
	)
}

func (u *Umwelt) FormatZettelTransactedDelta() schnittstellen.FuncWriterFormat[zettel.Transacted] {
	return zettel.MakeCliFormatTransactedDelta(
		u.FormatZettelTransacted(),
	)
}

func (u *Umwelt) FormatFileNotRecognized() schnittstellen.FuncWriterFormat[kennung.FD] {
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
		u.FormatZettelTransacted(),
	)
}

func (u *Umwelt) FormatFDDeleted() schnittstellen.FuncWriterFormat[kennung.FD] {
	return store_fs.MakeCliFormatFDDeleted(
		u.Konfig().DryRun,
		u.FormatColorWriter(),
		u.Standort(),
		u.FormatExternalFD(),
	)
}
