package objekte

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/checkout_mode"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/charlie/script_config"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/hotel/sku"
)

func MakeTextFormatter(
	af schnittstellen.AkteReaderFactory,
	k Konfig,
) textFormatter {
	return MakeTextFormatterWithAkteFormatter(af, k, nil)
}

func MakeTextFormatterWithAkteFormatter(
	af schnittstellen.AkteReaderFactory,
	k Konfig,
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
