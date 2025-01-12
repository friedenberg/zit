package repo_layout

import (
	"bufio"
	"encoding/gob"
	"flag"
	"io"
	"os"
	"path/filepath"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/age"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/echo/triple_hyphen_io"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/builtin_types"
)

type BigBang struct {
	ids.Type
	Config immutable_config.Latest

	AgeIdentity          age.Identity
	Yin                  string
	Yang                 string
	ExcludeDefaultType   bool
	ExcludeDefaultConfig bool
	OverrideXDGWithCwd   bool
}

func (bb *BigBang) SetFlagSet(f *flag.FlagSet) {
	f.Var(&bb.AgeIdentity, "age", "")
	f.BoolVar(&bb.OverrideXDGWithCwd, "override-xdg-with-cwd", false, "")
	f.StringVar(&bb.Yin, "yin", "", "File containing list of zettel id left parts")
	f.StringVar(&bb.Yang, "yang", "", "File containing list of zettel id right parts")

	bb.Type = builtin_types.GetOrPanic(builtin_types.ImmutableConfigV1).Type
	bb.Config = immutable_config.Default()
	bb.Config.BlobStore.SetFlagSet(f)
}

func (s *Layout) Genesis(bb BigBang) {
	s.Config.Type = bb.Type
	s.Config.Config = bb.Config

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

		if err := s.Age().AddIdentityOrGenerateIfNecessary(
			bb.AgeIdentity,
			s.FileAge(),
		); err != nil {
			s.CancelWithError(err)
		}

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
				Metadata: metadata{Config: &s.Config},
				Blob:     &s.Config,
			}

			if _, err := thw.WriteTo(f); err != nil {
				s.CancelWithError(err)
			}
		}

		writeFile(s.FileConfigMutable(), "")
		writeFile(s.FileCacheDormant(), "")
	}
}

func (s Layout) readAndTransferLines(in, out string) (err error) {
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
