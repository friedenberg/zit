package store_fs

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
)

//go:generate stringer -type=PathOption
type PathOption int

const (
	PathOptionDefault = PathOption(iota)
	PathOptionLeft
	PathOptionMiddle
	PathOptionRight
	PathOptionTempLocal
	PathOptionTempOS
)

type CheckoutOptions struct {
	Path            PathOption
	ForceInlineBlob bool
	AllowConflicted bool
	checkout_options.TextFormatterOptions
}

func GetCheckoutOptionsFromOptions(
	options checkout_options.Options,
) CheckoutOptions {
	return GetCheckoutOptionsFromOptionsWithoutMode(options.OptionsWithoutMode)
}

func GetCheckoutOptionsFromOptionsWithoutMode(
	options checkout_options.OptionsWithoutMode,
) (fsOptions CheckoutOptions) {
	switch t := options.StoreSpecificOptions.(type) {
	case nil:
	case CheckoutOptions:
		fsOptions = t

	default:
		panic(fmt.Sprintf("expected %T or nil but got %T", fsOptions, t))
	}

	if options.IgnoreWorkspace {
		fsOptions.Path = PathOptionTempLocal
	}

	return
}
