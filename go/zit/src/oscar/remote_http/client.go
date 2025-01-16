package remote_http

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/alfa/repo_type"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/inventory_list_blobs"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
)

type Client struct {
	*env.Env
	http.Client
	// Repo repo.WorkingCopy
	Repo *local_working_copy.Repo
	// *local_working_copy.Repo
}

func (repo *Client) GetRepoType() repo_type.Type {
	return repo_type.TypeUnknown
}

func (repo *Client) GetInventoryListStore() sku.InventoryListStore {
	return repo
}

func (repo *Client) GetBlobStore() interfaces.BlobStore {
	return &HTTPBlobStore{client: repo}
}

func (repo *Client) MakeExternalQueryGroup(
	builderOptions query.BuilderOptions,
	externalQueryOptions sku.ExternalQueryOptions,
	args ...string,
) (qg *query.Group, err error) {
	err = todo.Implement()
	return
}

func (client *Client) MakeInventoryList(
	qg *query.Group,
) (list *sku.List, err error) {
	var request *http.Request

	if request, err = http.NewRequestWithContext(
		client.Env.Context,
		"GET",
		"/inventory_lists",
		strings.NewReader(qg.String()),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	var response *http.Response

	if response, err = client.Do(request); err != nil {
		err = errors.Errorf("failed to read response: %w", err)
		return
	}

	bf := client.Repo.GetStore().GetInventoryListStore().FormatForVersion(
		client.Repo.GetConfig().GetStoreVersion(),
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

func (remoteHTTP *Client) PullQueryGroupFromRemote(
	remote repo.Repo,
	qg *query.Group,
	options repo.RemoteTransferOptions,
) (err error) {
	return remoteHTTP.pullQueryGroupFromWorkingCopy(
		remote.(repo.WorkingCopy),
		qg,
		options,
	)
}

func (remoteHTTP *Client) pullQueryGroupFromWorkingCopy(
	remote repo.WorkingCopy,
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
			remoteHTTP.Context,
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

	remoteHTTP.ContinueOrPanicOnDone()

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

func (remote *Client) WriteBlobToRemote(
	local repo.WorkingCopy,
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
		remote.Context,
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

func (remote *Client) ReadObjectHistory(
	oid *ids.ObjectId,
) (skus []*sku.Transacted, err error) {
	err = todo.Implement()
	return
}

type HTTPBlobStore struct {
	client *Client
}

func (blobStore *HTTPBlobStore) GetBlobStore() interfaces.BlobStore {
	return blobStore
}

func (blobStore *HTTPBlobStore) HasBlob(sh interfaces.Sha) (ok bool) {
	var request *http.Request

	{
		var err error

		if request, err = http.NewRequestWithContext(
			blobStore.client.Context,
			"HEAD",
			"/blobs",
			strings.NewReader(sh.GetShaLike().GetShaString()),
		); err != nil {
			blobStore.client.CancelWithError(err)
		}
	}

	var response *http.Response

	{
		var err error

		if response, err = blobStore.client.Do(request); err != nil {
			blobStore.client.CancelWithError(err)
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
		blobStore.client.Context,
		"GET",
		"/blobs",
		strings.NewReader(sh.GetShaLike().GetShaString()),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	var response *http.Response

	if response, err = blobStore.client.Do(request); err != nil {
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
