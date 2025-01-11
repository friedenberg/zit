package repo_local

import (
	"encoding/gob"
	"io"
	"os"
	"path"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/echo/repo_layout"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/builtin_types"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/config_mutable_cli"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/golf/mutable_config_blobs"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/hotel/type_blobs"
)

func Genesis(
	bb repo_layout.BigBang,
	context *errors.Context,
	config config_mutable_cli.Config,
	options env.Options,
) (u *Repo, err error) {
	dirLayout := dir_layout.MakeDefaultAndInitialize(
		context,
		config.Debug,
		bb.OverrideXDGWithCwd,
	)

	env := env.Make(
		context,
		config,
		dirLayout,
		options,
	)

	u = Make(env, OptionsEmpty)

	repoLayout := u.GetRepoLayout()
	repoLayout.Genesis(bb)

	if err = u.dormantIndex.Flush(
		u.GetRepoLayout(),
		u.PrinterHeader(),
		u.config.DryRun,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	u.Must(u.Reset)
	u.Must(repoLayout.ResetCache)

	if err = u.initDefaultTypeAndConfig(bb); err != nil {
		err = errors.Wrap(err)
		return
	}

	u.Must(u.Lock)
	u.Must(u.GetStore().ResetIndexes)
	u.Must(u.Unlock)

	return
}

func (repo *Repo) initDefaultTypeAndConfig(bb repo_layout.BigBang) (err error) {
	if err = repo.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, repo.Unlock)

	var defaultTypeObjectId ids.Type

	if defaultTypeObjectId, err = repo.initDefaultTypeIfNecessaryAfterLock(
		bb,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = repo.initDefaultConfigIfNecessaryAfterLock(
		bb,
		defaultTypeObjectId,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (repo *Repo) initDefaultTypeIfNecessaryAfterLock(
	bb repo_layout.BigBang,
) (defaultTypeObjectId ids.Type, err error) {
	if bb.ExcludeDefaultType {
		return
	}

	defaultTypeObjectId = ids.MustType("md")
	defaultTypeBlob := type_blobs.Default()

	var k ids.ObjectId

	if err = k.SetWithIdLike(defaultTypeObjectId); err != nil {
		err = errors.Wrap(err)
		return
	}

	var sh interfaces.Sha

	if sh, _, err = repo.GetStore().GetBlobStore().GetTypeV1().SaveBlobText(
		&defaultTypeBlob,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	o := sku.GetTransactedPool().Get()
	defer sku.GetTransactedPool().Put(o)

	if err = o.ObjectId.SetWithIdLike(defaultTypeObjectId); err != nil {
		err = errors.Wrap(err)
		return
	}

	o.Metadata.Blob.ResetWithShaLike(sh)
	o.GetMetadata().Type = builtin_types.DefaultOrPanic(genres.Type)

	if err = repo.GetStore().CreateOrUpdate(
		o,
		sku.GetStoreOptionsCreate(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (repo *Repo) initDefaultConfigIfNecessaryAfterLock(
	bb repo_layout.BigBang,
	defaultTypeObjectId ids.Type,
) (err error) {
	if bb.ExcludeDefaultConfig {
		return
	}

	var sh interfaces.Sha
	var tipe ids.Type

	if sh, tipe, err = writeDefaultMutableConfig(
		repo,
		defaultTypeObjectId,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	newConfig := sku.GetTransactedPool().Get()

	if err = newConfig.ObjectId.SetWithIdLike(ids.Config{}); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = newConfig.SetBlobSha(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	newConfig.Metadata.Type.ResetWith(tipe)

	if err = repo.GetStore().CreateOrUpdate(
		newConfig,
		sku.GetStoreOptionsCreate(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func writeDefaultMutableConfig(
	u *Repo,
	dt ids.Type,
) (sh interfaces.Sha, tipe ids.Type, err error) {
	defaultMutableConfig := mutable_config_blobs.Default(dt)
	tipe = defaultMutableConfig.Type

	f := u.GetStore().GetConfigBlobFormat()

	var aw sha.WriteCloser

	if aw, err = u.GetRepoLayout().BlobWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, aw)

	if _, err = f.FormatParsedBlob(aw, defaultMutableConfig.Blob); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh = sha.Make(aw.GetShaLike())

	return
}

func mkdirAll(elements ...string) {
	err := os.MkdirAll(path.Join(elements...), os.ModeDir|0o755)
	errors.PanicIfError(err)
}

func writeFile(p string, contents any) {
	var f *os.File
	var err error

	if f, err = files.CreateExclusiveWriteOnly(p); err != nil {
		if errors.IsExist(err) {
			ui.Err().Printf("%s already exists, not overwriting", p)
			err = nil
		} else {
		}

		return
	}

	defer errors.PanicIfError(err)
	defer errors.DeferredCloser(&err, f)

	if s, ok := contents.(string); ok {
		_, err = io.WriteString(f, s)
	} else {
		enc := gob.NewEncoder(f)
		err = enc.Encode(contents)
	}
}
