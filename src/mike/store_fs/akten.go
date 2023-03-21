package store_fs

import (
	"io"
	"os"
	"path/filepath"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/id"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/delta/age_io"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/hotel/objekte_store"
	"github.com/friedenberg/zit/src/juliett/zettel"
)

func (s Store) ReadExternalZettelFromAktePath(
	p string,
) (cz zettel.CheckedOut, err error) {
	errors.TodoP3("use cache")

	if p, err = filepath.Abs(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	if p, err = filepath.Rel(s.Cwd(), p); err != nil {
		err = errors.Wrap(err)
		return
	}

	head, tail := id.HeadTailFromFileName(p)

	if cz.External.Sku.Kennung, err = kennung.MakeHinweis(head + "/" + tail); err != nil {
		err = errors.Wrap(err)
		return
	}

	var zt *zettel.Transacted

	if zt, err = s.storeObjekten.Zettel().ReadOne(
		&cz.External.Sku.Kennung,
	); err != nil {
		if errors.Is(err, objekte_store.ErrNotFound{}) {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	if zt != nil {
		cz.Internal = *zt
	}

	errors.TodoP4("capture this as a function")
	cz.External.Objekte = cz.Internal.Objekte
	cz.External.Sku.ObjekteSha = cz.Internal.Sku.ObjekteSha
	cz.External.Sku.Kennung = cz.Internal.Sku.Kennung

	var akteSha sha.Sha

	if akteSha, err = s.AkteShaFromPath(p); err != nil {
		err = errors.Wrapf(err, "path: %s", p)
		return
	}

	errors.TodoP2("add mod time")
	cz.External.Sku.FDs.Akte.Path = p
	cz.External.Objekte.Akte = akteSha
	// cz.Matches.Akten, _ = s.storeObjekten.ReadAkteSha(akteSha)

	cz.DetermineState()

	return
}

func (s Store) AkteShaFromPath(p string) (sh sha.Sha, err error) {
	errors.TodoP3("use cache")

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

	defer errors.Deferred(&err, f.Close)

	if _, err = io.Copy(aw, f); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh = sha.Make(aw.Sha())

	return
}
