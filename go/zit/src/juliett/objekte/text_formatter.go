package objekte

import (
	"io"

	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/src/charlie/script_config"
	"code.linenisgreat.com/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/india/akten"
)

func MakeTextFormatter(
	af schnittstellen.AkteReaderFactory,
	k akten.Konfig,
) textFormatter {
	return MakeTextFormatterWithAkteFormatter(af, k, nil)
}

func MakeTextFormatterWithAkteFormatter(
	af schnittstellen.AkteReaderFactory,
	k akten.Konfig,
	akteFormatter script_config.RemoteScript,
) textFormatter {
	return textFormatter{
		k:                 k,
		fMetadateiAndAkte: metadatei.MakeTextFormatterMetadateiInlineAkte(af, akteFormatter),
		fMetadateiOnly:    metadatei.MakeTextFormatterMetadateiOnly(af, akteFormatter),
		fAkteOnly:         metadatei.MakeTextFormatterExcludeMetadatei(af, akteFormatter),
	}
}

type textFormatter struct {
	k                                            akten.Konfig
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
