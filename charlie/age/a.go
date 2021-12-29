package age

import (
	_age "filippo.io/age"
	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/bravo/open_file_guard"
)

const (
	FileName = "AgeIdentity"
)

type (
	_AgeIdentity       = _age.Identity
	_AgeRecipient      = _age.Recipient
	_AgeX25519Identity = _age.X25519Identity
)

var (
	_Error                  = errors.Error
	_GenerateX25519Identity = _age.GenerateX25519Identity
	_ParseX25519Identity    = _age.ParseX25519Identity
	_ParseX25519Recipient   = _age.ParseX25519Recipient
	_ReadStringAll          = open_file_guard.ReadAllString
	_AgeDecrypt             = _age.Decrypt
	_AgeEncrypt             = _age.Encrypt
)
