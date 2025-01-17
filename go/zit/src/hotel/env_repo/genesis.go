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
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/triple_hyphen_io"
)

func (s *Env) Genesis(bb BigBang) {
	s.config.Type = bb.Type
	s.config.ImmutableConfig = bb.Config

	if err := s.MakeDir(
		s.DirObjectId(),
		s.DirCache(),
		s.DirLostAndFound(),
	); err != nil {
		s.CancelWithError(err)
	}

	for _, g := range []genres.Genre{genres.Blob, genres.InventoryList} {
		var d string
		var err error

		if d, err = s.DirObjectGenre(g); err != nil {
			if genres.IsErrUnsupportedGenre(err) {
				err = nil
				continue
			} else {
				s.CancelWithError(err)
			}
		}

		if err := s.MakeDir(d); err != nil {
			s.CancelWithError(err)
		}
	}

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
					s.FileConfigPermanent(),
				); err != nil {
					s.CancelWithError(err)
				}

				defer s.MustClose(f)
			}

			thw := triple_hyphen_io.Writer{
				Metadata: metadata{config: &s.config},
				Blob:     &s.config,
			}

			if _, err := thw.WriteTo(f); err != nil {
				s.CancelWithError(err)
			}
		}
	}

	if s.config.ImmutableConfig.GetRepoType() == repo_type.TypeWorkingCopy {
		if err := s.readAndTransferLines(
			bb.Yin,
			filepath.Join(s.DirObjectId(), "Yin"),
		); err != nil {
			s.CancelWithError(err)
		}

		if err := s.readAndTransferLines(
			bb.Yang,
			filepath.Join(s.DirObjectId(), "Yang"),
		); err != nil {
			s.CancelWithError(err)
		}

		writeFile(s.FileConfigMutable(), "")
		writeFile(s.FileCacheDormant(), "")
	}

	if err := s.setupStores(); err != nil {
		s.CancelWithError(err)
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
