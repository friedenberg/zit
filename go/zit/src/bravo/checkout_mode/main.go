package checkout_mode

import (
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
)

type Mode int

type Getter interface {
	GetCheckoutMode() (Mode, error)
}

const (
	ModeNone = Mode(iota)
	ModeMetadataOnly
	ModeMetadataAndBlob
	ModeBlobOnly
)

func (m Mode) String() string {
	switch m {
	case ModeNone:
		return "none"

	case ModeMetadataOnly:
		return "objekte-only"

	case ModeBlobOnly:
		return "akte-only"

	case ModeMetadataAndBlob:
		return "both"

	default:
		return "unknown"
	}
}

func (m *Mode) Set(v string) (err error) {
	v = strings.ToLower(strings.TrimSpace(v))

	switch v {
	case "":
		*m = ModeNone

	case "objekte-only":
		fallthrough
	case "objekte":
		fallthrough
	case "zettel":
		fallthrough
	case "zettel-only":
		*m = ModeMetadataOnly

	case "akte":
		fallthrough
	case "akte-only":
		*m = ModeBlobOnly

	case "both":
		*m = ModeMetadataAndBlob

	default:
		err = errors.Errorf("unsupported checkout mode: %s", v)
		return
	}

	return
}

func (m Mode) IncludesAkte() bool {
	switch m {
	case ModeMetadataAndBlob, ModeBlobOnly:
		return true

	default:
		return false
	}
}

func (m Mode) IncludesObjekte() bool {
	switch m {
	case ModeMetadataAndBlob, ModeMetadataOnly:
		return true

	default:
		return false
	}
}
