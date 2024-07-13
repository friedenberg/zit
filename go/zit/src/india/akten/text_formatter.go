package akten

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/charlie/script_config"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func MakeTextFormatter(
	options checkout_options.TextFormatterOptions,
	af interfaces.BlobReaderFactory,
	k Konfig,
) textFormatter {
	return MakeTextFormatterWithAkteFormatter(options, af, k, nil)
}

func MakeTextFormatterWithAkteFormatter(
	options checkout_options.TextFormatterOptions,
	af interfaces.BlobReaderFactory,
	k Konfig,
	akteFormatter script_config.RemoteScript,
) textFormatter {
	return textFormatter{
		k:                 k,
		fMetadateiAndAkte: metadatei.MakeTextFormatterMetadateiInlineAkte(options, af, akteFormatter),
		fMetadateiOnly:    metadatei.MakeTextFormatterMetadateiOnly(options, af, akteFormatter),
		fAkteOnly:         metadatei.MakeTextFormatterExcludeMetadatei(options, af, akteFormatter),
	}
}

type textFormatter struct {
	k                                            Konfig
	fMetadateiAndAkte, fMetadateiOnly, fAkteOnly metadatei.TextFormatter
}

func (tf textFormatter) WriteStringFormat(w io.Writer, s *sku.Transacted) (n int64, err error) {
	if gattung.Konfig.EqualsGattung(s.GetGattung()) {
		n, err = tf.fAkteOnly.FormatMetadatei(w, s)
	} else if tf.k.IsInlineTyp(s.GetTyp()) {
		n, err = tf.fMetadateiAndAkte.FormatMetadatei(w, s)
	} else {
		n, err = tf.fMetadateiOnly.FormatMetadatei(w, s)
	}

	return
}

func (tf textFormatter) WriteStringFormatWithMode(
	w io.Writer,
	s *sku.Transacted,
	mode checkout_mode.Mode,
) (n int64, err error) {
	if gattung.Konfig.EqualsGattung(s.GetGattung()) || mode == checkout_mode.ModeAkteOnly {
		n, err = tf.fAkteOnly.FormatMetadatei(w, s)
	} else if tf.k.IsInlineTyp(s.GetTyp()) {
		n, err = tf.fMetadateiAndAkte.FormatMetadatei(w, s)
	} else {
		n, err = tf.fMetadateiOnly.FormatMetadatei(w, s)
	}

	return
}
