package remote_http

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"

	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/builtin_types"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
)

func (server *Server) writeInventoryListLocalWorkingCopy(
	repo *local_working_copy.Repo,
	request Request,
	listSku *sku.Transacted,
) (response Response) {
	listSkuType := builtin_types.GetOrPanic(builtin_types.InventoryListTypeV1).Type

	blobStore := server.Repo.GetBlobStore()

	if listSku != nil {
		if listSku.GetGenre() != genres.InventoryList {
			response.Error(genres.MakeErrUnsupportedGenre(listSku.GetGenre()))
			return
		}

		if blobStore.HasBlob(listSku.GetBlobSha()) {
			response.StatusCode = http.StatusFound
			return
		}

		listSkuType = listSku.GetType()
	}

	typedInventoryListStore := server.Repo.GetTypedInventoryListBlobStore()

	var blobWriter sha.WriteCloser

	{
		var err error

		if blobWriter, err = blobStore.BlobWriter(); err != nil {
			response.Error(err)
			return
		}
	}

	var list *sku.List

	{
		var err error

		if list, err = typedInventoryListStore.ReadInventoryListBlob(
			listSkuType,
			bufio.NewReader(io.TeeReader(request.Body, blobWriter)),
		); err != nil {
			response.Error(err)
			return
		}
	}

	b := bytes.NewBuffer(nil)

	// TODO make option to read from headers
	// TODO add remote blob store
	importerOptions := sku.ImporterOptions{
		// TODO
		CheckedOutPrinter: repo.PrinterCheckedOutConflictsForRemoteTransfers(),
	}

	if request.Headers.Get("x-zit-remote_transfer_options-allow_merge_conflicts") == "true" {
		importerOptions.AllowMergeConflicts = true
	}

	importerOptions.BlobCopierDelegate = func(
		result sku.BlobCopyResult,
	) (err error) {
		server.Repo.GetEnv().ContinueOrPanicOnDone()

		if result.N != -1 {
			return
		}

		fmt.Fprintf(b, "%s\n", result.GetBlobSha())

		return
	}

	importer := server.Repo.MakeImporter(
		importerOptions,
		sku.GetStoreOptionsRemoteTransfer(),
	)

	if err := server.Repo.ImportList(
		list,
		importer,
	); err != nil {
		response.Error(err)
		return
	}

	response.StatusCode = http.StatusCreated
	response.Body = io.NopCloser(b)

	return
}
