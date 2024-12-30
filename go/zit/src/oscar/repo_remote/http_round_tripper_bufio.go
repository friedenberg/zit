package repo_remote

import (
	"bufio"
	"net/http"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
)

type HTTPRoundTripperBufio struct {
	*bufio.Writer
	*bufio.Reader
}

func (roundTripper *HTTPRoundTripperBufio) RoundTrip(
	request *http.Request,
) (response *http.Response, err error) {
	if err = request.Write(roundTripper.Writer); err != nil {
		err = errors.Errorf("failed to write to socket: %w", err)
		return
	}

	if err = roundTripper.Flush(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if response, err = http.ReadResponse(
		roundTripper.Reader,
		request,
	); err != nil {
		err = errors.Errorf("failed to read response: %w", err)
		return
	}

	return
}
