package typed_blob_store

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/charlie/script_config"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_repo"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
)

func MakeTextFormatter(
	envRepo env_repo.Env,
	options checkout_options.TextFormatterOptions,
	inlineTypeChecker ids.InlineTypeChecker,
	checkoutMode checkout_mode.Mode,
) textFormatter {
	return MakeTextFormatterWithBlobFormatter(
		envRepo,
		options,
		inlineTypeChecker,
		nil,
		checkoutMode,
	)
}

func MakeTextFormatterWithBlobFormatter(
	envRepo env_repo.Env,
	options checkout_options.TextFormatterOptions,
	inlineTypeChecker ids.InlineTypeChecker,
	formatter script_config.RemoteScript,
	checkoutMode checkout_mode.Mode,
) textFormatter {
	return textFormatter{
		options:           options,
		InlineTypeChecker: inlineTypeChecker,
		TextFormatterFamily: object_metadata.MakeTextFormatterFamily(
			object_metadata.Dependencies{
				EnvDir:        envRepo,
				BlobStore:     envRepo,
				BlobFormatter: formatter,
			},
		),
		checkoutMode: checkoutMode,
	}
}

type textFormatter struct {
	ids.InlineTypeChecker
	options checkout_options.TextFormatterOptions
	object_metadata.TextFormatterFamily
	checkoutMode checkout_mode.Mode
}

func (formatter textFormatter) EncodeStringTo(
	object *sku.Transacted,
	writer io.Writer,
) (n int64, err error) {
	context := object_metadata.TextFormatterContext{
		PersistentFormatterContext: object,
		TextFormatterOptions:       formatter.options,
	}

	switch formatter.checkoutMode {
	case checkout_mode.MetadataOnly:
		n, err = formatter.MetadataOnly.FormatMetadata(writer, context)

	default:
		if genres.Config.EqualsGenre(object.GetGenre()) {
			n, err = formatter.BlobOnly.FormatMetadata(writer, context)
		} else if formatter.InlineTypeChecker.IsInlineType(object.GetType()) {
			n, err = formatter.InlineBlob.FormatMetadata(writer, context)
		} else {
			n, err = formatter.MetadataOnly.FormatMetadata(writer, context)
		}
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
	} else if tf.InlineTypeChecker.IsInlineType(sk.GetType()) {
		n, err = tf.InlineBlob.FormatMetadata(w, ctx)
	} else {
		n, err = tf.MetadataOnly.FormatMetadata(w, ctx)
	}

	return
}
