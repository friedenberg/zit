package zettel_external

import (
	"io"

	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/foxtrot/zettel"
)

// [path@sha !typ "bez"]
// [path.akte_ext@sha]
func MakeCliFormatZettel(
	s standort.Standort,
	sf collections.WriterFuncFormat[sha.Sha],
	zf collections.WriterFuncFormat[zettel.Zettel],
) collections.WriterFuncFormat[Zettel] {
	return func(w io.Writer, z *Zettel) (n int64, err error) {
		return collections.WriteFormats(
			w,
			collections.MakeWriterLiteral("["),
			s.MakeWriterRelativePath(z.ZettelFD.Path),
			collections.MakeWriterLiteral("@"),
			collections.MakeWriterFormatFunc(sf, &z.Named.Stored.Sha),
			collections.MakeWriterLiteral(" "),
			collections.MakeWriterFormatFunc(zf, &z.Named.Stored.Zettel),
			collections.MakeWriterLiteral("]"),
		)
	}
}

// [path.akte_ext@sha]
func MakeCliFormatAkte(
	s standort.Standort,
	sf collections.WriterFuncFormat[sha.Sha],
) collections.WriterFuncFormat[Zettel] {
	return func(w io.Writer, z *Zettel) (n int64, err error) {
		return collections.WriteFormats(
			w,
			collections.MakeWriterLiteral("["),
			s.MakeWriterRelativePath(z.AkteFD.Path),
			collections.MakeWriterLiteral("@"),
			collections.MakeWriterFormatFunc(sf, &z.Named.Stored.Zettel.Akte),
			collections.MakeWriterLiteral("]"),
		)
	}
}
