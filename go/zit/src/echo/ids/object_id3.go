package ids

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/charlie/box"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
)

// TODO parse this directly
// /repo/tag
// /repo/zettel-id
// /repo/!type
// /browser/one/uno
// /browser/bookmark-1
// /browser/!md
// /browser/!md
func (oid *objectId2) ReadFromSeq(
	seq box.Seq,
) (err error) {
	switch {
	case seq.Len() == 0:
		err = errors.Errorf("empty seq")
		return

		// tag
	case seq.MatchAll(box.TokenTypeLiteral):
		oid.g = genres.Tag
		oid.right.Write(seq.At(0).Contents)
		return

		// !type
	case seq.MatchAll(box.TokenMatcherOp(box.OpType), box.TokenTypeLiteral):
		oid.g = genres.Type
		oid.right.Write(seq.At(1).Contents)
		return

		// %tag
	case seq.MatchAll(box.TokenMatcherOp(box.OpVirtual), box.TokenTypeLiteral):
		oid.g = genres.Tag
		oid.middle = box.OpVirtual
		oid.right.Write(seq.At(1).Contents)
		return

		// /repo
	case seq.MatchAll(box.TokenMatcherOp(box.OpPathSeparator), box.TokenTypeLiteral):
		oid.g = genres.Repo
		oid.middle = box.OpPathSeparator
		oid.right.Write(seq.At(1).Contents)
		return

		// @sha
	case seq.MatchAll(box.TokenMatcherOp('@'), box.TokenTypeLiteral):
		oid.g = genres.Blob
		oid.middle = '@'
		oid.right.Write(seq.At(1).Contents)
		return

		// zettel/id
	case seq.MatchAll(
		box.TokenTypeLiteral,
		box.TokenMatcherOp(box.OpPathSeparator),
		box.TokenTypeLiteral,
	):
		oid.g = genres.Zettel
		oid.right.Write(seq.At(0).Contents)
		oid.middle = box.OpPathSeparator
		oid.right.Write(seq.At(2).Contents)
		return

	default:
		err = errors.Errorf("unsupported seq: %q", seq)
		return
	}
}
