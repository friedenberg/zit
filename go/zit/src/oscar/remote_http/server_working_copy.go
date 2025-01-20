package remote_http

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"

	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/inventory_list_blobs"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
)

func (server *Server) writeInventoryListLocalWorkingCopy(
	repo *local_working_copy.Repo,
	request Request,
) (response Response) {
	bf := server.Repo.GetInventoryListStore().FormatForVersion(
		server.Repo.GetImmutableConfig().ImmutableConfig.GetStoreVersion(),
	)

	list := sku.MakeList()

	if err := inventory_list_blobs.ReadInventoryListBlob(
		bf,
		bufio.NewReader(request.Body),
		list,
	); err != nil {
		response.Error(err)
		return
	}

	b := bytes.NewBuffer(nil)

	// TODO make option to read from headers
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

		sh := sha.GetPool().Get()
		sha.GetPool().Put(sh)
		sh.ResetWithShaLike(result.GetBlobSha())
		fmt.Fprintf(b, "%s\n", sh)

		return
	}

	// TODO
	importer := server.Repo.MakeImporter(
		importerOptions,
		sku.GetStoreOptionsRemoteTransfer(),
	)

	// TODO
	if err := server.Repo.ImportList(
		list,
		importer,
	); err != nil {
		response.Error(err)
		return
	}

	response.StatusCode = http.StatusCreated

	if b.Len() > 0 {
		response.Body = io.NopCloser(b)
	}

	return
}
