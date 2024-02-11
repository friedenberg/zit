package checkout_mode

import (
	"strings"

	"code.linenisgreat.com/zit-go/src/alfa/errors"
)

type Mode int

type Getter interface {
	GetCheckoutMode() (Mode, error)
}

const (
	ModeObjekteOnly = Mode(iota)
	ModeObjekteAndAkte
	ModeAkteOnly
)

func (m Mode) String() string {
	switch m {
	case ModeObjekteOnly:
		return "objekte-only"

	case ModeAkteOnly:
		return "akte-only"

	case ModeObjekteAndAkte:
		return "both"

	default:
		return "unknown"
	}
}

func (m *Mode) Set(v string) (err error) {
	v = strings.ToLower(strings.TrimSpace(v))

	switch v {
	case "objekte-only":
		fallthrough
	case "objekte":
		fallthrough
	case "zettel":
		fallthrough
	case "zettel-only":
		*m = ModeObjekteOnly

	case "akte":
		fallthrough
	case "akte-only":
		*m = ModeAkteOnly

	case "both":
		*m = ModeObjekteAndAkte

	default:
		err = errors.Errorf("unsupported checkout mode: %s", v)
		return
	}

	return
}

func (m Mode) IncludesAkte() bool {
	switch m {
	case ModeObjekteAndAkte, ModeAkteOnly:
		return true

	default:
		return false
	}
}

func (m Mode) IncludesObjekte() bool {
	switch m {
	case ModeObjekteAndAkte, ModeObjekteOnly:
		return true

	default:
		return false
	}
}
