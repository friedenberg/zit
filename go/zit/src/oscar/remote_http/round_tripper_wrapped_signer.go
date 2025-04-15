package remote_http

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/charlie/repo_signing"
)

const (
	headerChallengeNonce    = "X-Zit-Challenge-Nonce"
	headerChallengeResponse = "X-Zit-Challenge-Response"
	headerRepoPublicKey     = "X-Zit-Repo-Public_Key"
	headerSha256Sig         = "X-Zit-Sha256-Sig"
)

type RoundTripperBufioWrappedSigner struct {
	repo_signing.PublicKey
	roundTripperBufio
}

// TODO extract signing into an agnostic middleware
func (roundTripper *RoundTripperBufioWrappedSigner) RoundTrip(
	request *http.Request,
) (response *http.Response, err error) {
	nonceBytes := make([]byte, 32)

	if _, err = rand.Read(nonceBytes); err != nil {
		err = errors.Wrap(err)
		return
	}

	nonceString := base64.URLEncoding.EncodeToString(nonceBytes)

	if len(roundTripper.PublicKey) > 0 {
		request.Header.Add(headerChallengeNonce, nonceString)
	}

	if response, err = roundTripper.roundTripperBufio.RoundTrip(
		request,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if len(roundTripper.PublicKey) > 0 {
		if err = repo_signing.VerifyBase64Signature(
			roundTripper.PublicKey,
			nonceBytes,
			response.Header.Get(headerChallengeResponse),
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
