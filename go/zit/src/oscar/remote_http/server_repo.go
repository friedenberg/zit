package remote_http

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/charlie/repo_signing"
	"code.linenisgreat.com/zit/go/zit/src/charlie/tridex"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/india/log_remote_inventory_lists"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
)

func (server *Server) writeInventoryList(
	request Request,
	listSku *sku.Transacted,
) (response Response) {
	if listSku.GetGenre() != genres.InventoryList {
		response.Error(genres.MakeErrUnsupportedGenre(listSku.GetGenre()))
		return
	}

	blobStore := server.Repo.GetBlobStore()

	if blobStore.HasBlob(listSku.GetBlobSha()) {
		response.StatusCode = http.StatusFound
		return
	}

	expected := sha.Make(listSku.GetBlobSha())

	pubBase64 := request.request.Header.Get(headerRepoPublicKey)

	var logEntry log_remote_inventory_lists.Entry

	if pubBase64 != "" {
		{
			var err error

			if logEntry.PublicKey, err = base64.URLEncoding.DecodeString(pubBase64); err != nil {
				response.Error(err)
				return
			}
		}

		logEntry.EntryType = log_remote_inventory_lists.EntryTypeReceived
		logEntry.Transacted = listSku

		sig := request.request.Header.Get(headerSha256Sig)

		if err := repo_signing.VerifyBase64Signature(
			logEntry.PublicKey,
			expected.GetShaBytes(),
			sig,
		); err != nil {
			response.Error(err)
			return
		}
	}

	if len(logEntry.PublicKey) > 0 {
		if err := server.logRemoteInventoryLists.Exists(
			logEntry,
		); collections.IsErrNotFound(err) && err != nil {
			err = nil
		} else if err != nil {
			response.Error(err)
			return
		} else {
			response.StatusCode = http.StatusFound
			return
		}
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

	seqInventoryListSkus := typedInventoryListStore.IterInventoryListBlobSkusFromReader(
		listSku.GetType(),
		bufio.NewReader(io.TeeReader(request.Body, blobWriter)),
	)

	b := bytes.NewBuffer(nil)
	writtenNeededBlobs := tridex.Make()

	{
		count := 0

		for sk, err := range seqInventoryListSkus {
			server.Repo.GetEnv().ContinueOrPanicOnDone()

			if err != nil {
				response.Error(err)
				return
			}

			blobSha := sk.GetBlobSha()

			var ok bool
			ok, err = server.blobCache.HasBlob(blobSha)
			if err != nil {
				response.Error(err)
				return
			}

			blobShaString := blobSha.String()

			if ok || writtenNeededBlobs.ContainsExpansion(blobShaString) {
				continue
			}

			ui.Log().Printf("missing blob: %s", blobSha)

			fmt.Fprintf(b, "%s\n", blobSha)
			writtenNeededBlobs.Add(blobShaString)
			count++
		}

		ui.Err().Printf("missing blobs: %d", count)
	}

	if err := blobWriter.Close(); err != nil {
		response.Error(err)
		return
	}

	actual := blobWriter.GetShaLike()

	if err := expected.AssertEqualsShaLike(actual); err != nil {
		ui.Err().Printf(
			"received list has different sha: expected: %s, actual: %s",
			expected,
			actual,
		)

		// response.ErrorWithStatus(http.StatusBadRequest, err)
		// return
	}

	ui.Log().Printf("list sha matches: %s", expected)

	// TODO make merge conflicts impossible

	response.StatusCode = http.StatusCreated
	response.Body = io.NopCloser(b)

	if err := server.Repo.GetObjectStore().Commit(
		listSku,
		sku.CommitOptions{},
	); err != nil {
		response.Error(err)
		return
	}

	return
}
