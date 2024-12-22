package env

import (
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/india/box_format"
	"code.linenisgreat.com/zit/go/zit/src/mike/store"
)

func (env *Local) MakeImporter(
	remoteBlobStore dir_layout.BlobStore,
	printCopies bool,
) (importer store.Importer) {
	coPrinter := env.PrinterCheckedOut(box_format.CheckedOutHeaderState{})

	// TODO setup parent negotiator
	importer = store.Importer{
		Store:           env.GetStore(),
		ErrPrinter:      coPrinter,
		RemoteBlobStore: remoteBlobStore,
	}

	if printCopies {
		importer.BlobCopierDelegate = func(result store.BlobCopyResult) error {
			// TODO switch to Err and fix test
			return ui.Out().Printf(
				"copied Blob %s (%d bytes)",
				result.GetBlobSha(),
				result.N,
			)
		}
	}

	return
}
