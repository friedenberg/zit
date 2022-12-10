package zettel_checked_out

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/echo/standort"
	"github.com/friedenberg/zit/src/india/zettel_external"
)

// (same|changed) [path@sha !typ "bez"]
// (same|changed) [path.akte_ext@sha]
func MakeCliFormat(
	s standort.Standort,
	zef format.FormatWriterFunc[zettel_external.Zettel],
	aef format.FormatWriterFunc[zettel_external.Zettel],
	mode Mode,
) format.FormatWriterFunc[Zettel] {
	wzef := makeWriterFuncZettel(zef, false)
	waef := makeWriterFuncAkte(aef, false)

	switch mode {
	case ModeAkteOnly:
		return func(w io.Writer, z *Zettel) (n int64, err error) {
			return format.Write(
				w,
				format.MakeWriter(waef, z),
			)
		}

	case ModeZettelAndAkte:
		return func(w io.Writer, z *Zettel) (n int64, err error) {
			return format.Write(
				w,
				format.MakeWriter(wzef, z),
				format.MakeFormatString("\n"),
				format.MakeWriter(waef, z),
			)
		}

	case ModeZettelOnly:
		return func(w io.Writer, z *Zettel) (n int64, err error) {
			return format.Write(
				w,
				format.MakeWriter(wzef, z),
			)
		}

	default:
		return func(w io.Writer, z *Zettel) (n int64, err error) {
			err = errors.Errorf("unsupported checkout mode: %s", mode)
			return
		}
	}
}

func MakeCliFormatFresh(
	s standort.Standort,
	zef format.FormatWriterFunc[zettel_external.Zettel],
	aef format.FormatWriterFunc[zettel_external.Zettel],
	mode Mode,
) format.FormatWriterFunc[Zettel] {
	wzef := makeWriterFuncZettel(zef, true)
	waef := makeWriterFuncAkte(aef, true)

	switch mode {
	case ModeAkteOnly:
		return func(w io.Writer, z *Zettel) (n int64, err error) {
			return format.Write(
				w,
				format.MakeFormatStringRightAlignedParen(format.StringCheckedOut),
				format.MakeWriter(wzef, z),
			)
		}

	case ModeZettelAndAkte:
		return func(w io.Writer, z *Zettel) (n int64, err error) {
			return format.Write(
				w,
				format.MakeFormatStringRightAlignedParen(format.StringCheckedOut),
				format.MakeWriter(wzef, z),
				format.MakeFormatString("\n"),
				format.MakeFormatStringRightAlignedParen(format.StringCheckedOut),
				format.MakeWriter(waef, z),
			)
		}

	case ModeZettelOnly:
		return func(w io.Writer, z *Zettel) (n int64, err error) {
			return format.Write(
				w,
				format.MakeFormatStringRightAlignedParen(format.StringCheckedOut),
				format.MakeWriter(wzef, z),
			)
		}

	default:
		return func(w io.Writer, z *Zettel) (n int64, err error) {
			err = errors.Errorf("unsupported checkout mode: %s", mode)
			return
		}
	}
}

func makeWriterFuncZettel(
	zef format.FormatWriterFunc[zettel_external.Zettel],
	fresh bool,
) format.FormatWriterFunc[Zettel] {
	if fresh {
		return func(w io.Writer, z *Zettel) (n int64, err error) {
			return format.Write(
				w,
				format.MakeWriter(zef, &z.External),
			)
		}
	} else {
		return func(w io.Writer, z *Zettel) (n int64, err error) {
			diff := format.StringChanged

			if z.Internal.Sku.Sha.Equals(z.External.Sku.Sha) {
				diff = format.StringSame
			}

			return format.Write(
				w,
				format.MakeFormatStringRightAlignedParen(diff),
				format.MakeWriter(zef, &z.External),
			)
		}
	}
}

func makeWriterFuncAkte(
	aef format.FormatWriterFunc[zettel_external.Zettel],
	fresh bool,
) format.FormatWriterFunc[Zettel] {
	if fresh {
		return func(w io.Writer, z *Zettel) (n int64, err error) {
			return format.Write(
				w,
				format.MakeWriter(aef, &z.External),
			)
		}
	} else {
		return func(w io.Writer, z *Zettel) (n int64, err error) {
			diff := format.StringChanged

			if z.Internal.Objekte.Akte.Equals(z.External.Sku.Sha) {
				diff = format.StringSame
			}

			return format.Write(
				w,
				format.MakeFormatStringRightAlignedParen(diff),
				format.MakeWriter(aef, &z.External),
			)
		}
	}
}
