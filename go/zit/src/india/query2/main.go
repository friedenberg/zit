package query2

import (
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/matcher_proto"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/india/matcher"
)

func MakeGroupFromCheckedOutSet(
	cos sku.CheckedOutSet,
) (q matcher_proto.QueryGroup, err error) {
	return matcher.MakeGroupFromCheckedOutSet(cos)
}

func MakeGroup(
	k schnittstellen.Konfig,
	cwd matcher_proto.Cwd,
	ex kennung.Abbr,
	hidden matcher_proto.Matcher,
	feg schnittstellen.FileExtensionGetter,
	dg kennung.Gattung,
	ki kennung.Index,
) matcher_proto.QueryGroupBuilder {
	return matcher.MakeGroup(k, cwd, ex, hidden, feg, dg, ki)
}

func MakeGroupAll(
	k schnittstellen.Konfig,
	cwd matcher_proto.Cwd,
	ex kennung.Abbr,
	hidden matcher_proto.Matcher,
	feg schnittstellen.FileExtensionGetter,
	ki kennung.Index,
) matcher_proto.QueryGroup {
	return matcher.MakeGroupAll(k, cwd, ex, hidden, feg, ki)
}
