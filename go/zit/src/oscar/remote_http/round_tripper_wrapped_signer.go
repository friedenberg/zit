package remote_http

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
)

type roundTripperWrappedSigner struct {
	ed25519.PublicKey
	http.RoundTripper
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
		request.Header.Add("X-Zit-Challenge-Nonce", nonceString)
	}

	if response, err = roundTripper.RoundTripper.RoundTrip(request); err != nil {
		err = errors.Wrap(err)
		return
	}

	if len(roundTripper.PublicKey) > 0 && false {
		sig := response.Header.Get("X-Zit-Challenge-Response")

		if err = ed25519.VerifyWithOptions(
			roundTripper.PublicKey,
			nonceBytes,
			[]byte(sig),
			&ed25519.Options{},
		); err != nil {
			err = errors.Wrapf(err, "invalid signature: %q", sig)
			return
		}
	}

	return
}
