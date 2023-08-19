package ts

import "github.com/friedenberg/zit/src/echo/kennung"

const (
	Epoch          = kennung.Epoch
	FormatDateTime = kennung.FormatDateTime
	// FormatDateTai  = "%y-%m-%d %H:%M"
)

type Time = kennung.Time

var (
	Now           = kennung.Now
	Tyme          = kennung.Tyme
	TimeWithIndex = kennung.TimeWithIndex
)
