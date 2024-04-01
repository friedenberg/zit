package organize_text_mode

import (
	"strings"

	"code.linenisgreat.com/zit/src/alfa/errors"
)

type Mode int

const (
	ModeInteractive = Mode(iota)
	ModeCommitDirectly
	ModeOutputOnly
	ModeUnknown = -1
)

func (m *Mode) Set(v string) (err error) {
	switch strings.ToLower(v) {
	case "interactive":
		*m = ModeInteractive
	case "commit-directly":
		*m = ModeCommitDirectly
	case "output-only":
		*m = ModeOutputOnly
	default:
		*m = ModeUnknown
		err = errors.Errorf("unsupported mode: %s", v)
	}

	return
}

func (m Mode) String() string {
	switch m {
	case ModeInteractive:
		return "interactive"
	case ModeCommitDirectly:
		return "commit-directly"
	case ModeOutputOnly:
		return "output-only"
	default:
		return "unknown"
	}
}
