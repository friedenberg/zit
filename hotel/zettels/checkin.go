package zettels

func (zs *zettels) Checkin(options CheckinOptions, paths ...string) (daZees map[_Hinweis]ExternalZettel, err error) {
	if daZees, err = zs.ReadExternal(options, paths...); err != nil {
		err = _Error(err)
		return
	}

	for h, z := range daZees {
		var named _NamedZettel

		if named, err = zs.Read(h); err != nil {
			err = _Error(err)
			return
		}

		// if options.IgnoreAkte {
		// 	c.Zettel.Akte = named.Zettel.Akte
		// 	c.Zettel.AkteExt = named.Zettel.AkteExt
		// }

		named.Zettel = z.Zettel

		if _, err = zs.Update(named); err != nil {
			err = _Error(err)
			return
		}
	}

	return
}
