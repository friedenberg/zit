package umwelt

import (
	"github.com/friedenberg/zit/src/bravo/format"
	"github.com/friedenberg/zit/src/charlie/bezeichnung"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/etikett"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/echo/typ"
	"github.com/friedenberg/zit/src/foxtrot/zettel"
	"github.com/friedenberg/zit/src/hotel/zettel_named"
	"github.com/friedenberg/zit/src/india/zettel_external"
	"github.com/friedenberg/zit/src/india/zettel_transacted"
	"github.com/friedenberg/zit/src/juliett/zettel_checked_out"
	store_fs "github.com/friedenberg/zit/src/mike/store_fs"
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

func (u *Umwelt) FormatTyp() format.FormatWriterFunc[typ.Typ] {
	return typ.MakeCliFormat(u.FormatColorWriter())
}

func (u *Umwelt) FormatZettel() format.FormatWriterFunc[zettel.Zettel] {
	return zettel.MakeCliFormat(
		u.FormatBezeichnung(),
		format.MakeFormatStringer[etikett.Set](),
		u.FormatTyp(),
	)
}

func (u *Umwelt) FormatZettelNamed() format.FormatWriterFunc[zettel_named.Zettel] {
	return zettel_named.MakeCliFormat(
		u.FormatHinweis(),
		u.FormatSha(),
		u.FormatZettel(),
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

// TODO support tty-colored output
func (u *Umwelt) FormatZettelCheckedOut() format.FormatWriterFunc[zettel_checked_out.Zettel] {
	return zettel_checked_out.MakeCliFormat(
		u.Standort(),
		u.FormatZettelExternal(),
		u.FormatZettelExternalAkte(),
	)
}

// TODO support tty-colored output
func (u *Umwelt) FormatZettelCheckedOutFresh() format.FormatWriterFunc[zettel_checked_out.Zettel] {
	return zettel_checked_out.MakeCliFormatFresh(
		u.Standort(),
		u.FormatZettelExternal(),
		u.FormatZettelExternalAkte(),
	)
}

// TODO support tty-colored output
func (u *Umwelt) FormatZettelTransacted() format.FormatWriterFunc[zettel_transacted.Zettel] {
	return zettel_transacted.MakeCliFormat(
		u.FormatZettelNamed(),
	)
}

func (u *Umwelt) FormatFileNotRecognized() format.FormatWriterFunc[store_fs.File] {
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
		u.FormatZettelNamed(),
	)
}

func (u *Umwelt) FormatDirDeleted() format.FormatWriterFunc[store_fs.Dir] {
	return store_fs.MakeCliFormatDirDeleted(
		u.FormatColorWriter(),
		u.Standort(),
	)
}
