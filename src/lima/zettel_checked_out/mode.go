package zettel_checked_out

import (
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
)

// TODO-P3 move to zettel
type Mode int

const (
	ModeZettelOnly = Mode(iota)
	ModeZettelAndAkte
	ModeAkteOnly
)

func (m Mode) String() string {
	switch m {
	case ModeZettelOnly:
		return "zettel-only"

	case ModeAkteOnly:
		return "akte-only"

	case ModeZettelAndAkte:
		return "both"

	default:
		return "unknown"
	}
}

func (m *Mode) Set(v string) (err error) {
	v = strings.ToLower(strings.TrimSpace(v))

	switch v {
	case "zettel":
		fallthrough
	case "zettel-only":
		*m = ModeZettelOnly

	case "akte":
		fallthrough
	case "akte-only":
		*m = ModeAkteOnly

	case "both":
		*m = ModeZettelAndAkte

	default:
		err = errors.Errorf("unsupported checkout mode: %s", v)
		return
	}

	return
}

func (m Mode) IncludesAkte() bool {
	switch m {
	case ModeZettelAndAkte, ModeAkteOnly:
		return true

	default:
		return false
	}
}

func (m Mode) IncludesZettel() bool {
	switch m {
	case ModeZettelAndAkte, ModeZettelOnly:
		return true

	default:
		return false
	}
}
