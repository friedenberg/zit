package store

// type UnsureMatchMap struct {
// 	sku.UnsureMatchType
// 	lookup map[sha.Bytes]sku.CheckedOutLikeMutableSet
// }

// func (s *Store) QueryUnsure(
// 	qg *query.Group,
// 	o sku.UnsureMatchOptions,
// 	f sku.IterMatching,
// ) (err error) {
// 	matchMaps := o.MakeMatchMap()

// 	if err = s.cwdFiles.QueryUnsure(
// 		qg,
// 		sku.MakeUnsureMatchMapsCollector(matchMaps),
// 	); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	// TODO create a new query group for all of history
// 	qg.SetIncludeHistory()

// 	if matchMaps.Len() == 0 {
// 		return
// 	}

// 	if err = s.QueryWithKasten(
// 		qg,
// 		sku.MakeUnsureMatchMapsMatcher(
// 			matchMaps,
// 			f,
// 		),
// 	); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	return
// }
