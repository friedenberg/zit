package store_fs

import (
	"io"
	"os"
	"path/filepath"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/delta/id"
	"github.com/friedenberg/zit/src/echo/age_io"
	"github.com/friedenberg/zit/src/juliett/zettel_checked_out"
	"github.com/friedenberg/zit/src/lima/store_objekten"
)

func (s Store) ReadExternalZettelFromAktePath(p string) (cz zettel_checked_out.Zettel, err error) {
	if p, err = filepath.Abs(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	if p, err = filepath.Rel(s.path, p); err != nil {
		err = errors.Wrap(err)
		return
	}

	head, tail := id.HeadTailFromFileName(p)

	if cz.External.Named.Hinweis, err = hinweis.Make(head + "/" + tail); err != nil {
		err = errors.Wrap(err)
		return
	}

	if cz.Internal, err = s.storeObjekten.ReadHinweisSchwanzen(cz.External.Named.Hinweis); err != nil {
		if errors.Is(err, store_objekten.ErrNotFound{}) {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	cz.External.Named.Stored = cz.Internal.Named.Stored

	var akteSha sha.Sha

	if akteSha, err = s.AkteShaFromPath(p); err != nil {
		err = errors.Wrapf(err, "path: %s", p)
		return
	}

	//TODO add mod time
	cz.External.AkteFD.Path = p
	cz.External.Named.Stored.Zettel.Akte = akteSha
	// cz.Matches.Akten, _ = s.storeObjekten.ReadAkteSha(akteSha)

	cz.DetermineState()

	return
}

func (s Store) AkteShaFromPath(p string) (sh sha.Sha, err error) {
	var aw age_io.Writer

	if aw, err = s.storeObjekten.AkteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var f *os.File

	if f, err = files.Open(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer files.Close(f)

	if _, err = io.Copy(aw, f); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh = aw.Sha()

	return
}