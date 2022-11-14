package zettel_checked_out

import (
	"io"

	"github.com/friedenberg/zit/src/bravo/format"
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/india/zettel_external"
)

// (same|changed) [path@sha !typ "bez"]
// (same|changed) [path.akte_ext@sha]
func MakeCliFormat(
	s standort.Standort,
	zef format.FormatWriterFunc[zettel_external.Zettel],
	aef format.FormatWriterFunc[zettel_external.Zettel],
) format.FormatWriterFunc[Zettel] {
	return func(w io.Writer, z *Zettel) (n int64, err error) {
		var wtsZettel, wtsAkte format.WriterFunc

		if !z.External.ZettelFD.IsEmpty() {
			wtsZettel = func(w io.Writer) (n int64, err error) {
				diff := format.StringChanged

				if z.Internal.Named.Stored.Sha.Equals(z.External.Named.Stored.Sha) {
					diff = format.StringSame
				}

				return format.Write(
					w,
					format.MakeFormatStringRightAlignedParen(diff),
					format.MakeWriter(zef, &z.External),
				)
			}
		}

		if !z.External.AkteFD.IsEmpty() {
			wtsAkte = func(w io.Writer) (n int64, err error) {
				diff := format.StringChanged

				if z.Internal.Named.Stored.Zettel.Akte.Equals(z.External.Named.Stored.Zettel.Akte) {
					diff = format.StringSame
				}

				return format.Write(
					w,
					format.MakeFormatStringRightAlignedParen(diff),
					format.MakeWriter(aef, &z.External),
				)
			}
		}

		ws := []format.WriterFunc{}

		if wtsZettel != nil {
			ws = append(ws, wtsZettel)

			if wtsAkte != nil {
				ws = append(ws, format.MakeFormatString("\n"))
			}
		}

		if wtsAkte != nil {
			ws = append(ws, wtsAkte)
		}

		return format.Write(w, ws...)
	}
}

func MakeCliFormatFresh(
	s standort.Standort,
	zef format.FormatWriterFunc[zettel_external.Zettel],
	aef format.FormatWriterFunc[zettel_external.Zettel],
) format.FormatWriterFunc[Zettel] {
	return func(w io.Writer, z *Zettel) (n int64, err error) {
		var wtsZettel, wtsAkte format.WriterFunc

		if !z.External.ZettelFD.IsEmpty() {
			wtsZettel = func(w io.Writer) (n int64, err error) {
				return format.Write(
					w,
					format.MakeFormatStringRightAlignedParen(format.StringCheckedOut),
					format.MakeWriter(zef, &z.External),
				)
			}
		}

		if !z.External.AkteFD.IsEmpty() {
			wtsAkte = func(w io.Writer) (n int64, err error) {
				return format.Write(
					w,
					format.MakeFormatStringRightAlignedParen(format.StringCheckedOut),
					format.MakeWriter(aef, &z.External),
				)
			}
		}

		ws := []format.WriterFunc{}

		if wtsZettel != nil {
			ws = append(ws, wtsZettel)

			if wtsAkte != nil {
				ws = append(ws, format.MakeFormatString("\n"))
			}
		}

		if wtsAkte != nil {
			ws = append(ws, wtsAkte)
		}

		return format.Write(w, ws...)
	}
}
