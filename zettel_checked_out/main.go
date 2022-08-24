package zettel_checked_out

import "github.com/friedenberg/zit/src/golf/stored_zettel"

type CheckedOut struct {
	Internal stored_zettel.Transacted
	External stored_zettel.External
}
