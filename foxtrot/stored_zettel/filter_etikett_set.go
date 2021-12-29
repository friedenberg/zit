package stored_zettel

type FilterEtikettSet _EtikettSet

func (f FilterEtikettSet) IncludeNamedZettel(z Named) bool {
	ft := _EtikettSet(f)
	set := z.Zettel.Etiketten.IntersectPrefixes(ft)
	return set.Len() == ft.Len()
}

// func (f FilterEtikettSet) Set(v string) (err error) {
// 	ft := _EtikettSet(f)

// 	if err = ft.Set(v); err != nil {
// 		err = _Error(err)
// 		return
// 	}

// 	return
// }
