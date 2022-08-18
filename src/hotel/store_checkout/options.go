package checkout_store

import (
	"github.com/friedenberg/zit/src/foxtrot/zettel"
)

type CheckoutMode int

const (
	CheckoutModeZettelOnly = CheckoutMode(iota)
	CheckoutModeZettelAndAkte
	CheckoutModeAkteOnly
)

type CheckinOptions struct {
	IgnoreMissingHinweis bool
	AddMdExtension       bool
	IncludeAkte          bool
	Format               zettel.Format
}

type CheckoutOptions struct {
	CheckoutMode
	zettel.Format
}
