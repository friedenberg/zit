package zettel_checked_out

import "github.com/friedenberg/zit/src/juliett/zettel"

type Zettel = zettel.CheckedOut

// TODO-P0 remove
// type Zettel struct {
// 	zettel.CheckedOut
// }

// func (sz Zettel) String() string {
// 	return collections.MakeKey(
// 		sz.Internal.Sku.Kopf,
// 		sz.Internal.Sku.Mutter[0],
// 		sz.Internal.Sku.Mutter[1],
// 		sz.Internal.Sku.Schwanz,
// 		sz.Internal.Sku.Kennung,
// 		sz.Internal.Sku.ObjekteSha,
// 	)
// }

// func (a Zettel) EqualsAny(b any) bool {
// 	return values.Equals(a, b)
// }

// func (a Zettel) Equals(b Zettel) bool {
// 	if !a.Internal.Equals(b.Internal) {
// 		return false
// 	}

// 	if !a.External.Equals(b.External) {
// 		return false
// 	}

// 	if a.State != b.State {
// 		return false
// 	}

// 	return true
// }
