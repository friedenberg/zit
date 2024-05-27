package umwelt

import (
	"io"

	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/echo/standort"
	"code.linenisgreat.com/zit/src/echo/thyme"
	"code.linenisgreat.com/zit/src/juliett/konfig"
	"code.linenisgreat.com/zit/src/juliett/query"
	"code.linenisgreat.com/zit/src/mike/store"
)

func (u *Umwelt) Sonnenaufgang() thyme.Time {
	return u.sonnenaufgang
}

func (u *Umwelt) GetKonfig() *konfig.Compiled {
	return &u.konfig
}

func (u *Umwelt) Schlummernd() *query.Schlummernd {
	return &u.schlummernd
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

func (u *Umwelt) GetStore() *store.Store {
	return &u.store
}
