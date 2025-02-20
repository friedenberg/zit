package log_remote_inventory_lists

import (
	"bufio"
	"encoding/base64"
	"encoding/gob"
	"io"
	"os"
	"path/filepath"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/charlie/repo_signing"
	"code.linenisgreat.com/zit/go/zit/src/charlie/tridex"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_repo"
)

type v0 struct {
	env    env_repo.Env
	once   sync.Once
	path   string
	file   *os.File
	values interfaces.MutableTridex
}

func (log *v0) Flush() (err error) {
	if _, err = log.file.Seek(0, io.SeekStart); err != nil {
		err = errors.Wrap(err)
		return
	}

	bufferedWriter := bufio.NewWriter(log.file)

	enc := gob.NewEncoder(bufferedWriter)

	if err = enc.Encode(log.values); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = bufferedWriter.Flush(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = log.file.Close(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return nil
}

func (log *v0) initialize(env env_repo.Env, pubkey repo_signing.PublicKey) {
	gob.Register(tridex.Make())

	log.env = env
	log.values = tridex.Make()
	dir := env.DirCacheInventoryListLog()
	log.path = filepath.Join(dir, base64.URLEncoding.EncodeToString(pubkey))

	if err := log.env.MakeDir(dir); err != nil {
		env.CancelWithError(err)
		return
	}

	{
		var err error

		if log.file, err = files.CreateExclusiveWriteOnly(log.path); err != nil {
			if errors.IsExist(err) {
				if log.file, err = files.OpenExclusive(log.path); err != nil {
					env.CancelWithError(err)
					return
				}
			} else {
				env.CancelWithError(err)
				return
			}
		}
	}
}

func (log *v0) Append(entry Entry) (err error) {
	if err = log.readIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	log.values.Add(entry.GetShaString())

	return
}

func (log *v0) Exists(entry Entry) (err error) {
	if err = log.readIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !log.values.ContainsExpansion(entry.GetShaString()) {
		return collections.ErrNotFound
	}

	return
}

func (log *v0) readIfNecessary() (err error) {
	log.once.Do(
		func() {
			bufferedReader := bufio.NewReader(log.file)

			dec := gob.NewDecoder(bufferedReader)

			if err = dec.Decode(log.values); err != nil {
				err = errors.Wrap(err)
				return
			}
		},
	)

	return
}
