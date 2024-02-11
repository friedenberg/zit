package kennung

import "code.linenisgreat.com/zit/src/delta/thyme"

type Clock interface {
	GetTime() thyme.Time
	GetTai() Tai
}
