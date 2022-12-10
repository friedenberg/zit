package umwelt

import (
	"io"

	"github.com/friedenberg/zit/src/echo/standort"
	"github.com/friedenberg/zit/src/juliett/konfig_compiled"
	"github.com/friedenberg/zit/src/lima/store_objekten"
	"github.com/friedenberg/zit/src/lima/zettel_verzeichnisse"
	"github.com/friedenberg/zit/src/mike/store_fs"
)

func (u *Umwelt) Konfig() konfig_compiled.Compiled {
	return u.konfig
}

func (u *Umwelt) KonfigPtr() *konfig_compiled.Compiled {
	return &u.konfig
}

func (u *Umwelt) Out() io.Writer {
	return u.out
}

func (u *Umwelt) In() io.Reader {
	return u.in
}

func (u *Umwelt) Err() io.Writer {
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

func (u *Umwelt) ZettelVerzeichnissePool() zettel_verzeichnisse.Pool {
	return u.zettelVerzeichnissePool
}
