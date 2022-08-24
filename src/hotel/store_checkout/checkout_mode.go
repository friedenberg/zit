package store_checkout

import (
	"strings"

	"github.com/friedenberg/zit/src/bravo/errors"
)

type CheckoutMode int

const (
	CheckoutModeZettelOnly = CheckoutMode(iota)
	CheckoutModeZettelAndAkte
	CheckoutModeAkteOnly
)

func (m CheckoutMode) String() string {
	switch m {
	case CheckoutModeZettelOnly:
		return "zettel-only"

	case CheckoutModeAkteOnly:
		return "akte-only"

	case CheckoutModeZettelAndAkte:
		return "both"

	default:
		return "unknown"
	}
}

func (m *CheckoutMode) Set(v string) (err error) {
	v = strings.ToLower(strings.TrimSpace(v))

	switch v {
	case "zettel-only":
		*m = CheckoutModeZettelOnly

	case "akte-only":
		*m = CheckoutModeAkteOnly

	case "both":
		*m = CheckoutModeZettelAndAkte

	default:
		err = errors.Errorf("unsupported checkout mode: %s", v)
		return
	}

	return
}
