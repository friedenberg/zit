package umwelt

import (
	"io"

	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/delta/standort"
	"code.linenisgreat.com/zit/src/delta/thyme"
	"code.linenisgreat.com/zit/src/juliett/konfig"
	"code.linenisgreat.com/zit/src/mike/store_util"
)

func (u *Umwelt) Sonnenaufgang() thyme.Time {
	return u.sonnenaufgang
}

func (u *Umwelt) Konfig() *konfig.Compiled {
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

func (u *Umwelt) Store() *store_util.Store {
	return &u.store
}
