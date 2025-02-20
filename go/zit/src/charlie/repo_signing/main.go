package repo_signing

import (
	"crypto/ed25519"
	"encoding/base64"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
)

type (
	PrivateKey = ed25519.PrivateKey
	PublicKey  = ed25519.PublicKey
)

var NewKeyFromSeed = ed25519.NewKeyFromSeed

func SignBase64(key PrivateKey, message []byte) (signature string, err error) {
	var sig []byte

	if sig, err = key.Sign(
		nil,
		message,
		&ed25519.Options{},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	signature = base64.URLEncoding.EncodeToString(sig)

	return
}

func VerifyBase64Signature(
	publicKey PublicKey,
	message []byte,
	signatureBase64 string,
) (err error) {
	var sig []byte

	if sig, err = base64.URLEncoding.DecodeString(signatureBase64); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = ed25519.VerifyWithOptions(
		publicKey,
		message,
		sig,
		&ed25519.Options{},
	); err != nil {
		err = errors.Wrapf(err, "invalid signature: %q", signatureBase64)
		return
	}

	return
}
