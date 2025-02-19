package env_repo

import (
	"bufio"
	"encoding/gob"
	"io"
	"os"
	"path/filepath"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/repo_type"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/golf/config_immutable_io"
)

func (env *Env) Genesis(bb BigBang) {
	env.ConfigLoaded.Type = bb.Type
	env.ConfigLoaded.ImmutableConfig = bb.Config

	if err := env.MakeDir(
		env.DirObjectId(),
		env.DirCache(),
		env.DirLostAndFound(),
		env.DirInventoryLists(),
		env.DirBlobs(),
	); err != nil {
		env.CancelWithError(err)
	}

	writeFile(env.FileInventoryListLog(), "")

	{
		// if err := s.config.ImmutableConfig.GetBlobStoreImmutableConfig().GetAgeEncryption().AddIdentityOrGenerateIfNecessary(
		// 	bb.AgeIdentity,
		// ); err != nil {
		// 	if !errors.IsExist(err) {
		// 		s.CancelWithError(err)
		// 	}
		// }

		{
			var f *os.File

			{
				var err error

				if f, err = files.CreateExclusiveWriteOnly(
					env.FileConfigPermanent(),
				); err != nil {
					env.CancelWithError(err)
				}

				defer env.MustClose(f)
			}

			encoder := config_immutable_io.Coder{}

			if _, err := encoder.EncodeTo(&env.ConfigLoaded, f); err != nil {
				env.CancelWithError(err)
			}
		}
	}

	if env.ConfigLoaded.ImmutableConfig.GetRepoType() == repo_type.TypeWorkingCopy {
		if err := env.readAndTransferLines(
			bb.Yin,
			filepath.Join(env.DirObjectId(), "Yin"),
		); err != nil {
			env.CancelWithError(err)
		}

		if err := env.readAndTransferLines(
			bb.Yang,
			filepath.Join(env.DirObjectId(), "Yang"),
		); err != nil {
			env.CancelWithError(err)
		}

		writeFile(env.FileConfigMutable(), "")
		writeFile(env.FileCacheDormant(), "")
	}

	if err := env.setupStores(); err != nil {
		env.CancelWithError(err)
	}
}

func (s Env) readAndTransferLines(in, out string) (err error) {
	if in == "" {
		return
	}

	var fi, fo *os.File

	if fi, err = files.Open(in); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, fi.Close)

	if fo, err = files.CreateExclusiveWriteOnly(out); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, fo.Close)

	r := bufio.NewReader(fi)
	w := bufio.NewWriter(fo)

	defer errors.Deferred(&err, w.Flush)

	for {
		var l string
		l, err = r.ReadString('\n')

		if errors.Is(err, io.EOF) {
			err = nil
			break
		}

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		// TODO-P2 sterilize line
		w.WriteString(l)
	}

	return
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
