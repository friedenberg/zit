package kennung

import "code.linenisgreat.com/zit/src/echo/thyme"

type Clock interface {
	GetTime() thyme.Time
	GetTai() Tai
}
