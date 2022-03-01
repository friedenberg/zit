package zettels

func (zs *zettels) Checkin(options CheckinOptions, paths ...string) (checkedOut map[_Hinweis]_ZettelCheckedOut, err error) {
	if checkedOut, err = zs.ReadCheckedOut(options, paths...); err != nil {
		err = _Error(err)
		return
	}

	for _, z := range checkedOut {
		named := z.Internal
		named.Zettel = z.External.Zettel

		if _, err = zs.Update(named); err != nil {
			err = _Error(err)
			return
		}
	}

	return
}
