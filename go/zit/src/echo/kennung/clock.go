package kennung

import "github.com/friedenberg/zit/src/delta/thyme"

type Clock interface {
	GetTime() thyme.Time
	GetTai() Tai
}
