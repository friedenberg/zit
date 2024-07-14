package blob_store

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/charlie/script_config"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func MakeTextFormatter(
	options checkout_options.TextFormatterOptions,
	af interfaces.BlobReaderFactory,
	k Config,
) textFormatter {
	return MakeTextFormatterWithBlobFormatter(options, af, k, nil)
}

func MakeTextFormatterWithBlobFormatter(
	options checkout_options.TextFormatterOptions,
	af interfaces.BlobReaderFactory,
	k Config,
	formatter script_config.RemoteScript,
) textFormatter {
	return textFormatter{
		k:                 k,
		fMetadateiAndBlob: object_metadata.MakeTextFormatterMetadataInlineBlob(options, af, formatter),
		fMetadateiOnly:    object_metadata.MakeTextFormatterMetadataOnly(options, af, formatter),
		fBlobOnly:         object_metadata.MakeTextFormatterExcludeMetadata(options, af, formatter),
	}
}

type textFormatter struct {
	k                                            Config
	fMetadateiAndBlob, fMetadateiOnly, fBlobOnly object_metadata.TextFormatter
}

func (tf textFormatter) WriteStringFormat(w io.Writer, s *sku.Transacted) (n int64, err error) {
	if genres.Config.EqualsGenre(s.GetGenre()) {
		n, err = tf.fBlobOnly.FormatMetadata(w, s)
	} else if tf.k.IsInlineType(s.GetType()) {
		n, err = tf.fMetadateiAndBlob.FormatMetadata(w, s)
	} else {
		n, err = tf.fMetadateiOnly.FormatMetadata(w, s)
	}

	return
}

func (tf textFormatter) WriteStringFormatWithMode(
	w io.Writer,
	s *sku.Transacted,
	mode checkout_mode.Mode,
) (n int64, err error) {
	if genres.Config.EqualsGenre(s.GetGenre()) || mode == checkout_mode.ModeAkteOnly {
		n, err = tf.fBlobOnly.FormatMetadata(w, s)
	} else if tf.k.IsInlineType(s.GetType()) {
		n, err = tf.fMetadateiAndBlob.FormatMetadata(w, s)
	} else {
		n, err = tf.fMetadateiOnly.FormatMetadata(w, s)
	}

	return
}
