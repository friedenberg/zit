package umwelt

import (
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/charlie/bezeichnung"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/etikett"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/echo/typ"
	"github.com/friedenberg/zit/src/foxtrot/zettel"
	"github.com/friedenberg/zit/src/hotel/zettel_named"
	"github.com/friedenberg/zit/src/india/zettel_transacted"
	"github.com/friedenberg/zit/src/juliett/zettel_checked_out"
	store_fs "github.com/friedenberg/zit/src/mike/store_fs"
)

func (u *Umwelt) FormatSha() collections.WriterFuncFormat[sha.Sha] {
	return sha.MakeCliFormat(
		u.StoreObjekten(),
	)
}

func (u *Umwelt) FormatHinweis() collections.WriterFuncFormat[hinweis.Hinweis] {
	var a hinweis.Abbr

	if u.konfig.PrintAbbreviatedHinweisen {
		a = u.StoreObjekten()
	}

	return hinweis.MakeCliFormat(
		a,
		0,
		0,
	)
}

func (u *Umwelt) FormatZettel() collections.WriterFuncFormat[zettel.Zettel] {
	return zettel.MakeCliFormat(
		bezeichnung.MakeCliFormat(),
		collections.MakeWriterFormatStringer[etikett.Set](),
		collections.MakeWriterFormatStringer[typ.Typ](),
	)
}

func (u *Umwelt) FormatZettelNamed() collections.WriterFuncFormat[zettel_named.Zettel] {
	return zettel_named.MakeCliFormat(
		u.FormatHinweis(),
		u.FormatSha(),
		u.FormatZettel(),
	)
}

// TODO support tty-colored output
func (u *Umwelt) FormatZettelCheckedOut() collections.WriterFuncFormat[zettel_checked_out.Zettel] {
	return zettel_checked_out.MakeCliFormat(
		u.Standort(),
		u.FormatSha(),
		u.FormatZettel(),
	)
}

// TODO support tty-colored output
func (u *Umwelt) FormatZettelCheckedOutFresh() collections.WriterFuncFormat[zettel_checked_out.Zettel] {
	return zettel_checked_out.MakeCliFormatFresh(
		u.Standort(),
		u.FormatSha(),
		u.FormatZettel(),
	)
}

// TODO support tty-colored output
func (u *Umwelt) FormatZettelTransacted() collections.WriterFuncFormat[zettel_transacted.Zettel] {
	return zettel_transacted.MakeCliFormat(
		u.FormatZettelNamed(),
	)
}

// TODO support tty-colored output
func (u *Umwelt) FormatFileNotRecognized() collections.WriterFuncFormat[store_fs.File] {
	return store_fs.MakeCliFormatNotRecognized(
		u.Standort(),
		u.FormatSha(),
	)
}

func (u *Umwelt) FormatFileRecognized() collections.WriterFuncFormat[store_fs.FileRecognized] {
	return store_fs.MakeCliFormatRecognized(
		u.Standort(),
		u.FormatSha(),
		u.FormatZettelNamed(),
	)
}
