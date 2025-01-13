package repo_remote

import (
	"bufio"
	"net"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
)

type HTTPRoundTripperUnixSocket struct {
	local_working_copy.UnixSocket
	net.Conn
	HTTPRoundTripperBufio
}

func (roundTripper *HTTPRoundTripperUnixSocket) Initialize(
	remote *local_working_copy.Repo,
) (err error) {
	if roundTripper.UnixSocket, err = remote.InitializeUnixSocket(
		net.ListenConfig{},
		"",
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if roundTripper.Conn, err = net.Dial("unix", roundTripper.Path); err != nil {
		err = errors.Wrap(err)
		return
	}

	roundTripper.Writer = bufio.NewWriter(roundTripper.Conn)
	roundTripper.Reader = bufio.NewReader(roundTripper.Conn)

	return
}
