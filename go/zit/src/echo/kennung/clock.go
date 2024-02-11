package kennung

import "code.linenisgreat.com/zit-go/src/delta/thyme"

type Clock interface {
	GetTime() thyme.Time
	GetTai() Tai
}
