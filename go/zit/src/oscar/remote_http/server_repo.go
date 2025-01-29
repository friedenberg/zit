package remote_http

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"

	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
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

	typedInventoryListStore := server.Repo.GetTypedInventoryListBlobStore()

	var blobWriter sha.WriteCloser

	{
		var err error

		if blobWriter, err = blobStore.BlobWriter(); err != nil {
			response.Error(err)
			return
		}
	}

	seq := typedInventoryListStore.IterInventoryListBlobSkusFromReader(
		listSku.GetType(),
		bufio.NewReader(io.TeeReader(request.Body, blobWriter)),
	)

	b := bytes.NewBuffer(nil)

	for sk, err := range seq {
		server.Repo.GetEnv().ContinueOrPanicOnDone()

		if err != nil {
			response.Error(err)
			return
		}

		blobSha := sk.GetBlobSha()

		if blobStore.HasBlob(blobSha) {
			continue
		}

		sh := sha.GetPool().Get()
		sha.GetPool().Put(sh)
		sh.ResetWithShaLike(blobSha)
		fmt.Fprintf(b, "%s\n", sh)
	}

	if err := blobWriter.Close(); err != nil {
		response.Error(err)
		return
	}

	expected := sha.Make(listSku.GetBlobSha())
	actual := blobWriter.GetShaLike()

	if err := expected.AssertEqualsShaLike(actual); err != nil {
		response.Error(err)
		return
	}

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
