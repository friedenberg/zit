package zettel_checked_out

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/delta/zettel"
	"github.com/friedenberg/zit/src/golf/zettel_external"
)

// (same|diff) [path@sha !typ "bez"]
// (same|diff) [path.akte_ext@sha]
func MakeCliFormat(
	s standort.Standort,
	sf collections.WriterFuncFormat[sha.Sha],
	zf collections.WriterFuncFormat[zettel.Zettel],
) collections.WriterFuncFormat[Zettel] {
	return func(w io.Writer, z *Zettel) (n int64, err error) {
		var wtsZettel, wtsAkte collections.Writer

		switch {
		case !z.External.ZettelFD.IsEmpty():
			wtsZettel = func(w io.Writer) (n int64, err error) {
				diff := "diff"

				if z.Internal.Named.Stored.Sha.Equals(z.External.Named.Stored.Sha) {
					diff = "same"
				}

				return collections.WriteFormats(
					w,
					collections.MakeWriterLiteral("(%s) ", diff),
					collections.MakeWriterFormatFunc(
						zettel_external.MakeCliFormatZettel(
							s,
							sf,
							zf,
						),
						&z.External,
					),
				)
			}

		case !z.External.AkteFD.IsEmpty():
			wtsAkte = func(w io.Writer) (n int64, err error) {
				diff := "diff"

				if z.Internal.Named.Stored.Zettel.Akte.Equals(z.External.Named.Stored.Zettel.Akte) {
					diff = "same"
				}

				return collections.WriteFormats(
					w,
					collections.MakeWriterLiteral("(%s) ", diff),
					collections.MakeWriterFormatFunc(
						zettel_external.MakeCliFormatAkte(
							s,
							sf,
						),
						&z.External,
					),
				)
			}

		default:
			err = errors.Errorf("zettel external in unknown state: %q", z.External)
			return
		}

		ws := []collections.Writer{}

		if wtsZettel != nil {
			ws = append(ws, wtsZettel)

			if wtsAkte != nil {
				ws = append(ws, collections.MakeWriterLiteral("\n"))
			}
		}

		if wtsAkte != nil {
			ws = append(ws, wtsAkte)
		}

		return collections.WriteFormats(
			w,
			ws...,
		)
	}
}

func MakeCliFormatFresh(
	s standort.Standort,
	sf collections.WriterFuncFormat[sha.Sha],
	zf collections.WriterFuncFormat[zettel.Zettel],
) collections.WriterFuncFormat[Zettel] {
	return func(w io.Writer, z *Zettel) (n int64, err error) {
		var wtsZettel, wtsAkte collections.Writer

		switch {
		case !z.External.ZettelFD.IsEmpty():
			wtsZettel = func(w io.Writer) (n int64, err error) {
				return collections.WriteFormats(
					w,
					collections.MakeWriterLiteral("(checked out) "),
					collections.MakeWriterFormatFunc(
						zettel_external.MakeCliFormatZettel(
							s,
							sf,
							zf,
						),
						&z.External,
					),
				)
			}

		case !z.External.AkteFD.IsEmpty():
			wtsAkte = func(w io.Writer) (n int64, err error) {
				return collections.WriteFormats(
					w,
					collections.MakeWriterLiteral("(checked out) "),
					collections.MakeWriterFormatFunc(
						zettel_external.MakeCliFormatAkte(
							s,
							sf,
						),
						&z.External,
					),
				)
			}

		default:
			err = errors.Errorf("zettel external in unknown state: %q", z.External)
			return
		}

		ws := []collections.Writer{}

		if wtsZettel != nil {
			ws = append(ws, wtsZettel)

			if wtsAkte != nil {
				ws = append(ws, collections.MakeWriterLiteral("\n"))
			}
		}

		if wtsAkte != nil {
			ws = append(ws, wtsAkte)
		}

		return collections.WriteFormats(
			w,
			ws...,
		)
	}
}
