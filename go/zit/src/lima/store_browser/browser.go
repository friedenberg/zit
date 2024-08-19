package store_browser

import (
	"fmt"

	"code.linenisgreat.com/chrest/go/chrest"
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
)

func (s *Store) request(
	req chrest.BrowserRequest,
) (resp chrest.ResponseWithParsedJSONBody, err error) {
	if resp, err = s.store_browser.Request(req); err != nil {
		err = errors.Wrap(err)
		return
	}

	fmt.Println(resp.Response.Header.Get("X-Chrest-Startup-Time"))
	// TODO check launch time

	return
}
