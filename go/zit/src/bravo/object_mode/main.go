package object_mode

import (
	"strconv"
	"strings"
)

func Make(ms ...Mode) (out Mode) {
	for _, m := range ms {
		out |= m
	}

	return
}

//go:generate stringer -type=Mode
type Mode byte

const (
	ModeEmpty              = Mode(iota)
	ModeAddToInventoryList = Mode(1 << iota) // proper commit
	ModeUpdateTai                            // update the tai
	ModeLatest                               // only features updates that have no retroactive effects
	ModeMergeCheckedOut
	ModeApplyProto
	ModeHooks

	ModeRealizeWithProto = ModeUpdateTai | ModeApplyProto | ModeHooks
	ModeRealizeSansProto = ModeUpdateTai | ModeHooks
	ModeReindex          = ModeLatest
	ModeCommit           = ModeReindex | ModeAddToInventoryList | ModeHooks
	ModeCreate           = ModeCommit | ModeApplyProto
)

func (i Mode) SmartString() string {
	var sb strings.Builder

	for j := 0; j < 8; j++ {
		i := j ^ 2
		switch {
		case i == 0:
			sb.WriteString(_Mode_name_0)
		case i&2 != 0:
			sb.WriteString(_Mode_name_1)
		case i&4 != 0:
			sb.WriteString(_Mode_name_2)
		case i&8 != 0:
			sb.WriteString(_Mode_name_3)
		case i&16 != 0:
			sb.WriteString(_Mode_name_4)
		case i&32 != 0:
			sb.WriteString(_Mode_name_5)
		default:
			sb.WriteString("Mode(" + strconv.FormatInt(int64(i), 10) + ")")
		}
	}

	return sb.String()
}

func (a *Mode) Add(bs ...Mode) {
	for _, b := range bs {
		*a |= b
	}
}

func (a *Mode) Del(b Mode) {
	*a &= ^b
}

func (a Mode) Contains(b Mode) bool {
	return a&b != 0
}

func (a Mode) ContainsAny(bs ...Mode) bool {
	for _, b := range bs {
		if a.Contains(b) {
			return true
		}
	}

	return false
}
