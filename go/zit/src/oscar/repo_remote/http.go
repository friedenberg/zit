package repo_remote

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/inventory_list_blobs"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
	"code.linenisgreat.com/zit/go/zit/src/november/read_write_repo_local"
)

type HTTP struct {
	http.Client
	*read_write_repo_local.Repo
}

func (repo *HTTP) GetReadWriteRepo() repo.ReadWrite {
	return repo
}

func (repo *HTTP) GetBlobStore() interfaces.BlobStore {
	return &HTTPBlobStore{repo: repo}
}

func (repo *HTTP) MakeQueryGroup(
	builderOptions query.BuilderOptions,
	repoId ids.RepoId,
	externalQueryOptions sku.ExternalQueryOptions,
	args ...string,
) (qg *query.Group, err error) {
	err = todo.Implement()
	return
}

func (repo *HTTP) MakeInventoryList(
	qg *query.Group,
) (list *sku.List, err error) {
	var request *http.Request

	if request, err = http.NewRequestWithContext(
		repo.Repo.Context,
		"GET",
		"/inventory_lists",
		strings.NewReader(qg.String()),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	var response *http.Response

	if response, err = repo.Do(request); err != nil {
		err = errors.Errorf("failed to read response: %w", err)
		return
	}

	bf := repo.Repo.GetStore().GetInventoryListStore().FormatForVersion(
		repo.Repo.GetConfig().GetStoreVersion(),
	)

	list = sku.MakeList()

	if err = inventory_list_blobs.ReadInventoryListBlob(
		bf,
		bufio.NewReader(response.Body),
		list,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// func (remoteHTTP *HTTP) PullQueryGroupFromRemote2(
// 	remote repo.ReadWrite,
// 	options repo.RemoteTransferOptions,
// 	queryStrings ...string,
// ) (err error) {
// 	var qg *query.Group

// 	if qg, err = remoteHTTP.MakeQueryGroup(queryStrings...); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	if err = remoteHTTP.PullQueryGroupFromRemote(
// 		remote,
// 		qg,
// 		options,
// 	); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	return
// }

func (remoteHTTP *HTTP) PullQueryGroupFromRemote(
	remote repo.ReadWrite,
	qg *query.Group,
	options repo.RemoteTransferOptions,
) (err error) {
	var list *sku.List

	if list, err = remote.MakeInventoryList(qg); err != nil {
		err = errors.Wrap(err)
		return
	}

	// TODO local / remote version negotiation

	bf := remoteHTTP.Repo.GetStore().GetInventoryListStore().FormatForVersion(
		immutable_config.CurrentStoreVersion,
	)

	b := bytes.NewBuffer(nil)

	if _, err = bf.WriteInventoryListBlob(list, b); err != nil {
		err = errors.Wrap(err)
		return
	}

	var response *http.Response

	{
		var request *http.Request

		if request, err = http.NewRequestWithContext(
			remoteHTTP.Repo.Context,
			"POST",
			"/inventory_lists",
			b,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		if options.AllowMergeConflicts {
			request.Header.Add("x-zit-remote_transfer_options-allow_merge_conflicts", "true")
		}

		if response, err = remoteHTTP.Do(request); err != nil {
			err = errors.Errorf("failed to read response: %w", err)
			return
		}
	}

	if response.StatusCode >= 300 {
		var sb strings.Builder

		if _, err = io.Copy(&sb, response.Body); err != nil {
		}

		err = errors.Errorf("remote responded with error: %q", &sb)
		return
	}

	br := bufio.NewReader(response.Body)

	remoteHTTP.Repo.ContinueOrPanicOnDone()

	var shas sha.Slice

	if _, err = shas.ReadFrom(br); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = response.Body.Close(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if options.IncludeBlobs {
		for _, expected := range shas {
			if err = remoteHTTP.WriteBlobToRemote(remote, expected); err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	ui.Log().Print("done")

	return
}

func (remote *HTTP) WriteBlobToRemote(
	local repo.ReadWrite,
	expected *sha.Sha,
) (err error) {
	var actual sha.Sha

	// Closed by the http client's transport (our roundtripper calling
	// request.Write)
	var rc interfaces.ShaReadCloser

	if rc, err = local.GetBlobStore().BlobReader(
		expected,
	); err != nil {
		if dir_layout.IsErrBlobMissing(err) {
			// TODO make an option to collect this error at the present it, and an
			// option to fetch it from another remote store
			ui.Err().Printf("Blob missing locally: %q", expected)
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	var request *http.Request

	if request, err = http.NewRequestWithContext(
		remote.Repo.Context,
		"POST",
		"/blobs",
		rc,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	request.TransferEncoding = []string{"chunked"}

	var response *http.Response

	if response, err = remote.Do(request); err != nil {
		err = errors.Errorf("failed to read response: %w", err)
		return
	}

	var shString strings.Builder

	if _, err = io.Copy(&shString, response.Body); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = response.Body.Close(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = actual.Set(strings.TrimSpace(shString.String())); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = expected.AssertEqualsShaLike(&actual); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (remote *HTTP) ReadObjectHistory(
	oid *ids.ObjectId,
) (skus []*sku.Transacted, err error) {
	err = todo.Implement()
	return
}

type HTTPBlobStore struct {
	repo *HTTP
}

func (blobStore *HTTPBlobStore) GetBlobStore() interfaces.BlobStore {
	return blobStore
}

func (blobStore *HTTPBlobStore) HasBlob(sh interfaces.Sha) (ok bool) {
	var request *http.Request

	{
		var err error

		if request, err = http.NewRequestWithContext(
			blobStore.repo.Repo.Context,
			"HEAD",
			"/blobs",
			strings.NewReader(sh.GetShaLike().GetShaString()),
		); err != nil {
			blobStore.repo.Repo.CancelWithError(err)
		}
	}

	var response *http.Response

	{
		var err error

		if response, err = blobStore.repo.Do(request); err != nil {
			blobStore.repo.Repo.CancelWithError(err)
		}
	}

	ok = response.StatusCode == http.StatusNoContent

	return
}

func (blobStore *HTTPBlobStore) BlobWriter() (w interfaces.ShaWriteCloser, err error) {
	err = todo.Implement()
	return
}

func (blobStore *HTTPBlobStore) BlobReader(
	sh interfaces.Sha,
) (r interfaces.ShaReadCloser, err error) {
	var request *http.Request

	if request, err = http.NewRequestWithContext(
		blobStore.repo.Repo.Context,
		"GET",
		"/blobs",
		strings.NewReader(sh.GetShaLike().GetShaString()),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	var response *http.Response

	if response, err = blobStore.repo.Do(request); err != nil {
		err = errors.Wrap(err)
		return
	}

	// TODO refactor this into a common structure
	if response.StatusCode >= 300 {
		var sb strings.Builder

		if _, err = io.Copy(&sb, response.Body); err != nil {
		}

		err = errors.Errorf("remote responded with error: %q", &sb)
		return
	}

	r = sha.MakeReadCloser(response.Body)

	return
}
