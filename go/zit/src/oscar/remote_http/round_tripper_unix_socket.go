package remote_http

import (
	"bufio"
	"net"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/charlie/repo_signing"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
)

type RoundTripperUnixSocket struct {
	repo.UnixSocket
	net.Conn
	RoundTripperBufioWrappedSigner
}

// TODO add public key
func (roundTripper *RoundTripperUnixSocket) Initialize(
	remote *Server,
	pubkey repo_signing.PublicKey,
) (err error) {
	roundTripper.PublicKey = pubkey

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
