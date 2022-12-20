package store_fs

import "github.com/friedenberg/zit/src/mike/zettel_checked_out"

//TODO-P2 remove
type CheckoutMode = zettel_checked_out.Mode

const (
	CheckoutModeZettelOnly    = zettel_checked_out.ModeZettelOnly
	CheckoutModeZettelAndAkte = zettel_checked_out.ModeZettelAndAkte
	CheckoutModeAkteOnly      = zettel_checked_out.ModeAkteOnly
)

