package sha

import (
	"io"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
)

func DebugCompareBlobs(
	ctx errors.Context,
	blobStore interfaces.BlobStore,
	expectedSha interfaces.Sha,
	actual *strings.Builder,
) {
	var expected strings.Builder

	var blobReader interfaces.ShaReadCloser

	{
		var err error

		if blobReader, err = blobStore.BlobReader(expectedSha); err != nil {
			ctx.CancelWithError(err)
		}
	}

	if _, err := io.Copy(&expected, blobReader); err != nil {
		ctx.CancelWithError(err)
	}

	ui.Debug().Printf("expected: %q, actual: %q", &expected, actual)
}
