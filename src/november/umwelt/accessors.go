package umwelt

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/standort"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/hotel/transacted"
	"github.com/friedenberg/zit/src/india/konfig"
	"github.com/friedenberg/zit/src/lima/store_objekten"
	"github.com/friedenberg/zit/src/mike/store_fs"
)

func (u *Umwelt) Sonnenaufgang() kennung.Time {
	return u.sonnenaufgang
}

func (u *Umwelt) Konfig() konfig.Compiled {
	return u.konfig
}

func (u *Umwelt) KonfigPtr() *konfig.Compiled {
	return &u.konfig
}

func (u *Umwelt) In() io.Reader {
	return u.in
}

func (u *Umwelt) Out() schnittstellen.WriterAndStringWriter {
	return u.out
}

func (u *Umwelt) Err() schnittstellen.WriterAndStringWriter {
	return u.err
}

func (u *Umwelt) Standort() standort.Standort {
	return u.standort
}

func (u *Umwelt) StoreObjekten() *store_objekten.Store {
	return u.storeObjekten
}

func (u *Umwelt) StoreWorkingDirectory() *store_fs.Store {
	return u.storeWorkingDirectory
}

func (u *Umwelt) ZettelVerzeichnissePool() schnittstellen.Pool[transacted.Zettel, *transacted.Zettel] {
	return u.zettelVerzeichnissePool
}
