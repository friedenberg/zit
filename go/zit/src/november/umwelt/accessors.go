package umwelt

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/thyme"
	"code.linenisgreat.com/zit/go/zit/src/echo/fs_home"
	"code.linenisgreat.com/zit/go/zit/src/juliett/konfig"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/mike/store"
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

func (u *Umwelt) Out() interfaces.WriterAndStringWriter {
	return u.out
}

func (u *Umwelt) Err() interfaces.WriterAndStringWriter {
	return u.err
}

func (u *Umwelt) Standort() fs_home.Standort {
	return u.fs_home
}

func (u *Umwelt) GetStore() *store.Store {
	return &u.store
}
