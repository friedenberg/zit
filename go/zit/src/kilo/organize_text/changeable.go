package organize_text

// func (ot *Text) CompareMap(
// 	hinweis_expander func(string) (*kennung.Hinweis, error),
// ) (out compare_map.CompareMap, err error) {
// 	preExpansion := compare_map.CompareMap{
// 		Named:   make(compare_map.SetKeyToMetadatei),
// 		Unnamed: make(compare_map.SetKeyToMetadatei),
// 	}

// 	if err = ot.addToCompareMap(
// 		ot,
// 		ot.Metadatei,
// 		kennung.MakeEtikettSet(),
// 		&preExpansion,
// 	); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	out = compare_map.CompareMap{
// 		Named:   make(compare_map.SetKeyToMetadatei),
// 		Unnamed: preExpansion.Unnamed,
// 	}

// 	for h, v := range preExpansion.Named {
// 		var h1 schnittstellen.Stringer

// 		if h1, err = hinweis_expander(h); err == nil {
// 			h = h1.String()
// 		}

// 		err = nil

// 		out.Named[h] = v
// 	}

// 	return
// }

// func (ot *Text) GetSkus(
// 	original sku.TransactedSet,
// ) (out2 sku.TransactedSet, err error) {
// 	out := sku.MakeTransactedMutableSet()
// 	out2 = out

// 	if err = ot.addToSet(
// 		ot,
// 		out,
// 		original,
// 	); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	return
// }

// func (a *Assignment) addToSet(
// 	ot *Text,
// 	out sku.TransactedMutableSet,
// 	original sku.TransactedSet,
// ) (err error) {
// 	expanded := kennung.MakeEtikettMutableSet()

// 	if err = a.AllEtiketten(expanded); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	if err = a.Named.Each(
// 		func(o *obj) (err error) {
// 			var z *sku.Transacted
// 			ok := false

// 			if z, ok = out.Get(out.Key(&o.Transacted)); !ok {
// 				z = sku.GetTransactedPool().Get()

// 				if err = z.SetFromSkuLike(&o.Transacted); err != nil {
// 					err = errors.Wrap(err)
// 					return
// 				}

// 				if err = ot.EachPtr(
// 					z.Metadatei.AddEtikettPtr,
// 				); err != nil {
// 					err = errors.Wrap(err)
// 					return
// 				}

// 				if !ot.Metadatei.Typ.IsEmpty() {
// 					z.Metadatei.Typ.ResetWith(ot.Metadatei.Typ)
// 				}

// 				if err = out.Add(z); err != nil {
// 					err = errors.Wrap(err)
// 					return
// 				}

// 				zPrime, hasOriginal := original.Get(original.Key(&o.Transacted))

// 				if hasOriginal {
// 					z.Metadatei.Akte.ResetWith(&zPrime.Metadatei.Akte)
// 					z.Metadatei.Typ.ResetWith(zPrime.Metadatei.Typ)
// 				}

// 				if !ot.Metadatei.Typ.IsEmpty() {
// 					z.Metadatei.Typ.ResetWith(ot.Metadatei.Typ)
// 				}
// 			}

// 			if o.Kennung.String() == "" {
// 				panic(fmt.Sprintf("%s: Kennung is nil", o))
// 			}

// 			if err = z.Metadatei.Bezeichnung.Set(
// 				o.Metadatei.Bezeichnung.String(),
// 			); err != nil {
// 				err = errors.Wrap(err)
// 				return
// 			}

// 			if !o.Metadatei.Typ.IsEmpty() {
// 				if err = z.Metadatei.Typ.Set(
// 					o.Metadatei.Typ.String(),
// 				); err != nil {
// 					err = errors.Wrap(err)
// 					return
// 				}
// 			}

// 			z.Metadatei.Comments = append(
// 				z.Metadatei.Comments,
// 				o.Metadatei.Comments...,
// 			)

// 			if err = o.Metadatei.GetEtiketten().EachPtr(
// 				z.Metadatei.AddEtikettPtr,
// 			); err != nil {
// 				err = errors.Wrap(err)
// 				return
// 			}

// 			if err = expanded.EachPtr(
// 				z.Metadatei.AddEtikettPtr,
// 			); err != nil {
// 				err = errors.Wrap(err)
// 				return
// 			}

// 			return
// 		},
// 	); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	if err = a.Unnamed.Each(
// 		func(z *obj) (err error) {
// 			// TODO unnamed

// 			return
// 		},
// 	); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	for _, c := range a.Children {
// 		if err = c.addToSet(ot, out, original); err != nil {
// 			err = errors.Wrap(err)
// 			return
// 		}
// 	}

// 	return
// }
