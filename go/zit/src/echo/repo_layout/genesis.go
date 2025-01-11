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
)

type BigBang struct {
	AgeIdentity          age.Identity
	Yin                  string
	Yang                 string
	Config               immutable_config.Latest
	ExcludeDefaultType   bool
	ExcludeDefaultConfig bool
	OverrideXDGWithCwd   bool
}

func (e *BigBang) SetFlagSet(f *flag.FlagSet) {
	f.Var(&e.AgeIdentity, "age", "")
	f.BoolVar(&e.OverrideXDGWithCwd, "override-xdg-with-cwd", false, "")
	f.StringVar(&e.Yin, "yin", "", "File containing list of zettel id left parts")
	f.StringVar(&e.Yang, "yang", "", "File containing list of zettel id right parts")

	e.Config.SetFlagSet(f)
}

func (s Layout) Genesis(bb BigBang) {
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

		writeFile(s.FileConfigPermanent(), bb.Config)
		writeFile(s.FileConfigMutable(), "")
		writeFile(s.FileCacheDormant(), "")
	}

	return
}

func (s Layout) readAndTransferLines(in, out string) (err error) {
	ui.TodoP4("move to user operations")

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
