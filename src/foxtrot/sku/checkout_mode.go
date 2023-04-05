package sku

import (
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type CheckoutMode int

type CheckoutModeGetter interface {
	GetCheckoutMode() (CheckoutMode, error)
}

const (
	CheckoutModeObjekteOnly = CheckoutMode(iota)
	CheckoutModeObjekteAndAkte
	CheckoutModeAkteOnly
)

func (m CheckoutMode) String() string {
	switch m {
	case CheckoutModeObjekteOnly:
		return "objekte-only"

	case CheckoutModeAkteOnly:
		return "akte-only"

	case CheckoutModeObjekteAndAkte:
		return "both"

	default:
		return "unknown"
	}
}

func (m *CheckoutMode) Set(v string) (err error) {
	v = strings.ToLower(strings.TrimSpace(v))

	switch v {
	case "zettel":
		fallthrough
	case "zettel-only":
		*m = CheckoutModeObjekteOnly

	case "akte":
		fallthrough
	case "akte-only":
		*m = CheckoutModeAkteOnly

	case "both":
		*m = CheckoutModeObjekteAndAkte

	default:
		err = errors.Errorf("unsupported checkout mode: %s", v)
		return
	}

	return
}

func (m CheckoutMode) IncludesAkte() bool {
	switch m {
	case CheckoutModeObjekteAndAkte, CheckoutModeAkteOnly:
		return true

	default:
		return false
	}
}

func (m CheckoutMode) IncludesObjekte() bool {
	switch m {
	case CheckoutModeObjekteAndAkte, CheckoutModeObjekteOnly:
		return true

	default:
		return false
	}
}
