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
				diff := "changed"

				if z.Internal.Named.Stored.Sha.Equals(z.External.Named.Stored.Sha) {
					diff = "same"
				}

				return format.Write(
					w,
					format.MakeFormatString("(%s) ", diff),
					format.MakeWriter(zef, &z.External),
				)
			}
		}

		if !z.External.AkteFD.IsEmpty() {
			wtsAkte = func(w io.Writer) (n int64, err error) {
				diff := "changed"

				if z.Internal.Named.Stored.Zettel.Akte.Equals(z.External.Named.Stored.Zettel.Akte) {
					diff = "same"
				}

				return format.Write(
					w,
					format.MakeFormatString("(%s) ", diff),
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
					format.MakeFormatString("(checked out) "),
					format.MakeWriter(zef, &z.External),
				)
			}
		}

		if !z.External.AkteFD.IsEmpty() {
			wtsAkte = func(w io.Writer) (n int64, err error) {
				return format.Write(
					w,
					format.MakeFormatString("(checked out) "),
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
