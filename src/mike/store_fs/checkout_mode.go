package store_fs

import "github.com/friedenberg/zit/src/juliett/zettel_checked_out"

type CheckoutMode = zettel_checked_out.Mode

const (
	CheckoutModeZettelOnly    = zettel_checked_out.ModeZettelOnly
	CheckoutModeZettelAndAkte = zettel_checked_out.ModeZettelAndAkte
	CheckoutModeAkteOnly      = zettel_checked_out.ModeAkteOnly
)
