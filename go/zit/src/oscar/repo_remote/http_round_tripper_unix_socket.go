package repo_remote

import (
	"bufio"
	"net"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/november/repo_local"
)

type HTTPRoundTripperUnixSocket struct {
	repo_local.UnixSocket
	net.Conn
	HTTPRoundTripperBufio
}

func (roundTripper *HTTPRoundTripperUnixSocket) Initialize(
	remote *repo_local.Repo,
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
