package remote_http

import (
	"io"
	"net/http"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
)

func (client *Client) HasBlob(sh interfaces.Sha) (ok bool) {
	var request *http.Request

	{
		var err error

		if request, err = http.NewRequestWithContext(
			client.GetEnv().Context,
			"HEAD",
			"/blobs",
			strings.NewReader(sh.GetShaLike().GetShaString()),
		); err != nil {
			client.GetEnv().CancelWithError(err)
		}
	}

	var response *http.Response

	{
		var err error

		if response, err = client.Do(request); err != nil {
			client.GetEnv().CancelWithError(err)
		}
	}

	ok = response.StatusCode == http.StatusNoContent

	return
}

func (client *Client) BlobWriter() (w interfaces.ShaWriteCloser, err error) {
	err = todo.Implement()
	return
}

func (client *Client) BlobReader(
	sh interfaces.Sha,
) (r interfaces.ShaReadCloser, err error) {
	var request *http.Request

	if request, err = http.NewRequestWithContext(
		client.GetEnv().Context,
		"GET",
		"/blobs",
		strings.NewReader(sh.GetShaLike().GetShaString()),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	var response *http.Response

	if response, err = client.Do(request); err != nil {
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
