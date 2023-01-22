package remote_transfers

import "github.com/friedenberg/zit/src/alfa/errors"

var (
	ErrPullRemoteHasHigherVersion = errors.New("pull remote has higher version")
)
