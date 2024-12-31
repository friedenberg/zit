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
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/inventory_list_blobs"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
	"code.linenisgreat.com/zit/go/zit/src/november/repo_local"
)

type HTTP struct {
	http.Client
	remote *repo_local.Repo
}

func (remote *HTTP) GetRepo() repo.Repo {
	return remote
}

func (remote *HTTP) GetBlobStore() interfaces.BlobStore {
	return &HTTPBlobStore{remote: remote}
}

func (remote *HTTP) MakeQueryGroup(
	metaBuilder any,
	repoId ids.RepoId,
	externalQueryOptions sku.ExternalQueryOptions,
	args ...string,
) (qg *query.Group, err error) {
	err = todo.Implement()
	return
}

func (remote *HTTP) MakeInventoryList(
	qg *query.Group,
) (list *sku.List, err error) {
	var request *http.Request

	if request, err = http.NewRequestWithContext(
		remote.remote.Context,
		"GET",
		"/inventory_lists",
		strings.NewReader(qg.String()),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	var response *http.Response

	if response, err = remote.Do(request); err != nil {
		err = errors.Errorf("failed to read response: %w", err)
		return
	}

	bf := remote.remote.GetStore().GetInventoryListStore().FormatForVersion(
		remote.remote.GetConfig().GetStoreVersion(),
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

func (remoteHTTP *HTTP) PullQueryGroupFromRemote(
	remote repo.Repo,
	qg *query.Group,
	printCopies bool,
) (err error) {
	var list *sku.List

	if list, err = remote.MakeInventoryList(qg); err != nil {
		err = errors.Wrap(err)
		return
	}

	// TODO local / remote version negotiation

	bf := remoteHTTP.remote.GetStore().GetInventoryListStore().FormatForVersion(
		immutable_config.CurrentStoreVersion,
	)

	b := bytes.NewBuffer(nil)

	if _, err = bf.WriteInventoryListBlob(list, b); err != nil {
		err = errors.Wrap(err)
		return
	}

	var request *http.Request

	if request, err = http.NewRequestWithContext(
		remoteHTTP.remote.Context,
		"POST",
		"/inventory_lists",
		b,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	var response *http.Response

	if response, err = remoteHTTP.Do(request); err != nil {
		err = errors.Errorf("failed to read response: %w", err)
		return
	}

	if response.StatusCode >= 300 {
		var sb strings.Builder

		if _, err = io.Copy(&sb, response.Body); err != nil {
		}

		err = errors.Errorf("remote responded with error: %q", &sb)
		return
	}

	br := bufio.NewReader(response.Body)
	eof := false

	remoteHTTP.remote.ContinueOrPanicOnDone()

	var shas []*sha.Sha

	for !eof {
		remoteHTTP.remote.ContinueOrPanicOnDone()

		var line string
		line, err = br.ReadString('\n')

		if err == io.EOF {
			err = nil
			eof = true
		} else if err != nil {
			err = errors.Wrap(err)
			return
		}

		if line == "" {
			continue
		}

		sh := sha.GetPool().Get()

		if err = sh.Set(strings.TrimSpace(line)); err != nil {
			err = errors.Wrap(err)
			return
		}

		shas = append(shas, sh)
	}

	if err = response.Body.Close(); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, expected := range shas {
		var actual sha.Sha

		// Closed by the http client's transport (our roundtripper calling
		// request.Write)
		var rc interfaces.ShaReadCloser

		if rc, err = remote.GetBlobStore().BlobReader(
			expected,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		if request, err = http.NewRequestWithContext(
			remoteHTTP.remote.Context,
			"POST",
			"/blobs",
			rc,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		request.TransferEncoding = []string{"chunked"}

		var response *http.Response

		if response, err = remoteHTTP.Do(request); err != nil {
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
	}

	ui.Log().Print("done")

	return
}

func (remote *HTTP) ReadObjectHistory(
	oid *ids.ObjectId,
) (skus []*sku.Transacted, err error) {
	err = todo.Implement()
	return
}

type HTTPBlobStore struct {
	remote *HTTP
}

func (blobStore *HTTPBlobStore) GetBlobStore() interfaces.BlobStore {
	return blobStore
}

func (blobStore *HTTPBlobStore) HasBlob(sh interfaces.Sha) (ok bool) {
	var request *http.Request

	{
		var err error

		if request, err = http.NewRequestWithContext(
			blobStore.remote.remote.Context,
			"HEAD",
			"/blobs",
			strings.NewReader(sh.GetShaLike().GetShaString()),
		); err != nil {
			blobStore.remote.remote.CancelWithError(err)
		}
	}

	var response *http.Response

	{
		var err error

		if response, err = blobStore.remote.Do(request); err != nil {
			blobStore.remote.remote.CancelWithError(err)
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
		blobStore.remote.remote.Context,
		"GET",
		"/blobs",
		strings.NewReader(sh.GetShaLike().GetShaString()),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	var response *http.Response

	if response, err = blobStore.remote.Do(request); err != nil {
		err = errors.Wrap(err)
		return
	}

	if response.StatusCode != http.StatusOK {
		err = errors.Errorf("remote error: %d", response.StatusCode)
		return
	}

	r = sha.MakeReadCloser(response.Body)

	return
}
