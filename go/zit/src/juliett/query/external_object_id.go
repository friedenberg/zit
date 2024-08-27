package query

// type ExternalObjectId struct {
// 	Exact   bool
// 	Virtual bool
// 	Debug   bool

// 	sku.ExternalObjectId
// }

// // TODO support exact
// func (exp ExternalObjectId) ContainsSku(sk *sku.Transacted) (ok bool) {
// 	defer sk.Metadata.Cache.QueryPath.PushOnReturn(exp, &ok)

// 	skMe := sk.GetMetadata()

// 	switch exp.GetGenre() {
// 	case genres.Tag:
// 		var idx int

// 		if exp.Exact {
// 			idx, ok = skMe.Cache.TagPaths.All.ContainsObjectIdTagExact(
// 				exp.ExternalObjectId.String(),
// 			)
// 		} else {
// 			idx, ok = skMe.Cache.TagPaths.All.ContainsObjectIdTag(
// 				exp.GetObjectId(),
// 			)
// 		}

// 		ui.Log().Print(exp, idx, ok, skMe.Cache.TagPaths.All, sk)

// 		if ok {
// 			// if k.Exact {
// 			// 	ewp := me.Verzeichnisse.Etiketten.All[idx]
// 			// 	ui.Debug().Print(ewp, sk)
// 			// }

// 			ps := skMe.Cache.TagPaths.All[idx]
// 			sk.Metadata.Cache.QueryPath.Push(ps.Parents)
// 			return
// 		}

// 		return

// 	case genres.Type:
// 		if ids.Contains(skMe.GetType(), exp.GetObjectId()) {
// 			ok = true
// 			return
// 		}
// 	}

// 	idl := &sk.ObjectId

// 	if !ids.Contains(idl, exp.GetObjectId()) {
// 		return
// 	}

// 	ok = true

// 	return
// }

// func (k ExternalObjectId) String() string {
// 	var sb strings.Builder

// 	if k.Exact {
// 		sb.WriteRune('=')
// 	}

// 	if k.Virtual {
// 		sb.WriteRune('%')
// 	}

// 	sb.WriteString(ids.FormattedString(k.GetObjectId()))

// 	return sb.String()
// }
