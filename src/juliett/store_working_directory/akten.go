package store_working_directory

import (
	"io"
	"os"
	"path/filepath"

	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/charlie/open_file_guard"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/delta/id"
	"github.com/friedenberg/zit/src/echo/age_io"
	"github.com/friedenberg/zit/src/india/store_objekten"
	"github.com/friedenberg/zit/src/india/zettel_checked_out"
)

func (s Store) ReadExternalZettelFromAktePath(p string) (cz zettel_checked_out.CheckedOut, err error) {
	if p, err = filepath.Abs(p); err != nil {
		err = errors.Error(err)
		return
	}

	if p, err = filepath.Rel(s.path, p); err != nil {
		err = errors.Error(err)
		return
	}

	cz.External.AktePath = p

	head, tail := id.HeadTailFromFileName(p)

	if cz.External.Named.Hinweis, err = hinweis.Make(head + "/" + tail); err != nil {
		err = errors.Error(err)
		return
	}

	if cz.Internal, err = s.storeObjekten.Read(cz.External.Named.Hinweis); err != nil {
		if errors.Is(err, store_objekten.ErrNotFound{}) {
			err = nil
		} else {
			err = errors.Error(err)
			return
		}
	}

	cz.External.Named.Stored = cz.Internal.Named.Stored

	var akteSha sha.Sha

	if akteSha, err = s.AkteShaFromPath(p); err != nil {
		err = errors.Error(err)
		return
	}

	cz.External.AktePath = p
	cz.External.Named.Stored.Zettel.Akte = akteSha
	cz.Matches.Akten, _ = s.storeObjekten.ReadAkteSha(akteSha)

	cz.DetermineState()

	return
}

func (s Store) AkteShaFromPath(p string) (sh sha.Sha, err error) {
	var aw age_io.Writer

	if aw, err = s.storeObjekten.AkteWriter(); err != nil {
		err = errors.Error(err)
		return
	}

	var f *os.File

	if f, err = open_file_guard.Open(p); err != nil {
		err = errors.Error(err)
		return
	}

	defer open_file_guard.Close(f)

	if _, err = io.Copy(aw, f); err != nil {
		err = errors.Error(err)
		return
	}

	sh = aw.Sha()

	return
}
