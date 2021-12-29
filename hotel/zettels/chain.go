package zettels

type Chain struct {
	Hinweis _Hinweis
	//stored in reverse (latest is at 0)
	Zettels []_StoredZettel
}

func (zs zettels) AllInChain(id _Id) (c Chain, err error) {
	var s _Sha

	if s, c.Hinweis, err = zs.TailFromId(id); err != nil {
		err = _Error(err)
		return
	}

	shas := make(map[string]bool)

	for {
		if s.IsNull() {
			break
		}

		if _, ok := shas[s.String()]; ok {
			err = _Errorf("loop detected in history for sha '%s'", s)
			return
		}

		shas[s.String()] = true

		var sz _NamedZettel

		if sz, err = zs.Read(s); err != nil {
			err = _Error(err)
			return
		}

		c.Zettels = append(
			c.Zettels,
			sz.Stored,
		)

		s = sz.Mutter
	}

	return
}

func (zs zettels) TailFromId(id _Id) (s _Sha, h _Hinweis, err error) {
	ok := false

	if s, ok = id.(_Sha); ok {
		if h, err = zs.hinweisen.ReadSha(s); err != nil {
			err = _Error(err)
			return
		}
	} else {
		if h, ok = id.(_Hinweis); !ok {
			err = _Errorf("unsupported id: '%q'", id)
			return
		}

		if s, err = zs.hinweisen.Read(h); err != nil {
			err = _Error(err)
			return
		}
	}

	return
}
