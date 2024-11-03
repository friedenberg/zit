package blob_store

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/charlie/script_config"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func MakeTextFormatter(
	dirLayout dir_layout.DirLayout,
	options checkout_options.TextFormatterOptions,
	k Config,
) textFormatter {
	return MakeTextFormatterWithBlobFormatter(dirLayout, options, k, nil)
}

func MakeTextFormatterWithBlobFormatter(
	dirLayout dir_layout.DirLayout,
	options checkout_options.TextFormatterOptions,
	k Config,
	formatter script_config.RemoteScript,
) textFormatter {
	return textFormatter{
		options: options,
		k:       k,
		TextFormatterFamily: object_metadata.MakeTextFormatterFamily(
			dirLayout,
			formatter,
		),
	}
}

type textFormatter struct {
	k       Config
	options checkout_options.TextFormatterOptions
	object_metadata.TextFormatterFamily
}

func (tf textFormatter) WriteStringFormat(w io.Writer, z *sku.Transacted) (n int64, err error) {
	s := object_metadata.TextFormatterContext{
		PersistentFormatterContext: z,
		TextFormatterOptions:       tf.options,
	}
	if genres.Config.EqualsGenre(z.GetGenre()) {
		n, err = tf.BlobOnly.FormatMetadata(w, s)
	} else if tf.k.IsInlineType(z.GetType()) {
		n, err = tf.InlineBlob.FormatMetadata(w, s)
	} else {
		n, err = tf.MetadataOnly.FormatMetadata(w, s)
	}

	return
}

func (tf textFormatter) WriteStringFormatWithMode(
	w io.Writer,
	sk *sku.Transacted,
	mode checkout_mode.Mode,
) (n int64, err error) {
	ctx := object_metadata.TextFormatterContext{
		PersistentFormatterContext: sk,
		TextFormatterOptions:       tf.options,
	}

	if genres.Config.EqualsGenre(sk.GetGenre()) || mode == checkout_mode.BlobOnly {
		n, err = tf.BlobOnly.FormatMetadata(w, ctx)
	} else if tf.k.IsInlineType(sk.GetType()) {
		n, err = tf.InlineBlob.FormatMetadata(w, ctx)
	} else {
		n, err = tf.MetadataOnly.FormatMetadata(w, ctx)
	}

	return
}
