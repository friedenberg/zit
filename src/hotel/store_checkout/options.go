package checkout_store

import (
	"github.com/friedenberg/zit/src/foxtrot/zettel"
)

type CheckinOptions struct {
	IgnoreMissingHinweis bool
	AddMdExtension       bool
	IncludeAkte          bool
	Format               zettel.Format
}
