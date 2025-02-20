package remote_http

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
)

const (
	headerChallengeNonce    = "X-Zit-Challenge-Nonce"
	headerChallengeResponse = "X-Zit-Challenge-Response"
)

type roundTripperWrappedSigner struct {
	ed25519.PublicKey
	roundTripperBufio
}

func (roundTripper *roundTripperWrappedSigner) RoundTrip(
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

	if response, err = roundTripper.roundTripperBufio.RoundTrip(request); err != nil {
		err = errors.Wrap(err)
		return
	}

	if len(roundTripper.PublicKey) > 0 {
		sigBase64 := response.Header.Get(headerChallengeResponse)

		var sig []byte

		if sig, err = base64.URLEncoding.DecodeString(sigBase64); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = ed25519.VerifyWithOptions(
			roundTripper.PublicKey,
			nonceBytes,
			sig,
			&ed25519.Options{},
		); err != nil {
			err = errors.Wrapf(err, "invalid signature: %q", sigBase64)
			return
		}
	}

	return
}
