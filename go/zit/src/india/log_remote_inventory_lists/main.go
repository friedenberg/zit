package log_remote_inventory_lists

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/repo_signing"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_repo"
)

type Entry = interfaces.Sha

type Log interface {
	errors.Flusher
	initialize(env_repo.Env, repo_signing.PublicKey)
	Append(Entry) error
	Exists(Entry) error
}

func Make(env env_repo.Env, pubkey repo_signing.PublicKey) Log {
	sv := env.GetConfigPrivate().ImmutableConfig.GetStoreVersion()

	switch sv := sv.GetInt(); {
	case sv < 8:
		env.CancelWithErrorf("unsupported store version: %s")
		return nil

	default:
		log := &v0{}
		log.initialize(env, pubkey)
		return log
	}
}
