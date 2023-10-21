package umwelt

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/delta/thyme"
	"github.com/friedenberg/zit/src/juliett/konfig"
	"github.com/friedenberg/zit/src/mike/store_util"
	"github.com/friedenberg/zit/src/november/store_objekten"
)

func (u *Umwelt) Sonnenaufgang() thyme.Time {
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

func (u *Umwelt) StoreUtil() store_util.StoreUtil {
	return u.storeUtil
}
