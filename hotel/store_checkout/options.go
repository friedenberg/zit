package checkout_store

import "github.com/friedenberg/zit/foxtrot/zettel"

type CheckinOptions struct {
	IgnoreMissingHinweis bool
	AddMdExtension       bool
	IncludeAkte          bool
	Format               zettel.Format
}
