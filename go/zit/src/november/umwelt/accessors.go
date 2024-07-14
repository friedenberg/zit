package umwelt

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/thyme"
	"code.linenisgreat.com/zit/go/zit/src/echo/fs_home"
	"code.linenisgreat.com/zit/go/zit/src/india/dormant_index"
	"code.linenisgreat.com/zit/go/zit/src/juliett/config"
	"code.linenisgreat.com/zit/go/zit/src/mike/store"
)

func (u *Umwelt) Sonnenaufgang() thyme.Time {
	return u.sonnenaufgang
}

func (u *Umwelt) GetKonfig() *config.Compiled {
	return &u.config
}

func (u *Umwelt) GetDormantIndex() *dormant_index.Index {
	return &u.dormantIndex
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

func (u *Umwelt) Standort() fs_home.Home {
	return u.fs_home
}

func (u *Umwelt) GetStore() *store.Store {
	return &u.store
}
