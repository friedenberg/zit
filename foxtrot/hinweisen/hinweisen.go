package hinweisen

import (
	"path"
)

type Hinweisen interface {
	Read(h _Hinweis) (sha _Sha, err error)
	ReadSha(s _Sha) (h _Hinweis, err error)
	ReadManyStrings(args ...string) (shas []_Sha, hins []_Hinweis, err error)
	All() (shas []_Sha, hins []_Hinweis, err error)
	Store(s _Sha) (h _Hinweis, err error)
	Update(h _Hinweis, s _Sha) (err error)
	Flush() error
}

type hinweisen struct {
	basePath string
	storeH   _Store
	storeS   _Store
	factory  *factory
}

func New(age _Age, basePath string) (s *hinweisen, err error) {
	s = &hinweisen{
		basePath: basePath,
	}

	if s.factory, err = newFactory(basePath); err != nil {
		err = _Error(err)
		return
	}

	if s.storeS, err = _NewStore(path.Join(basePath, "Zettel-Hinweis"), s); err != nil {
		err = _Error(err)
		return
	}

	if s.storeH, err = _NewStore(path.Join(basePath, "Hinweis"), s); err != nil {
		err = _Error(err)
		return
	}

	return
}

func (hn hinweisen) NewShard(p string, id string) (s _Shard, err error) {
	if s, err = _NewShard(path.Join(p, id), nil, &_ShardGeneric{}); err != nil {
		err = _Error(err)
		return
	}

	return
}

func (zs *hinweisen) Flush() (err error) {
	if err = zs.storeH.Flush(); err != nil {
		err = _Error(err)
		return
	}

	if err = zs.storeS.Flush(); err != nil {
		err = _Error(err)
		return
	}

	if err = zs.factory.Flush(); err != nil {
		err = _Error(err)
		return
	}

	return
}

func (hn *hinweisen) Store(sha _Sha) (h _Hinweis, err error) {
	var ss _Shard

	if ss, err = hn.storeS.Shard(sha.Head()); err != nil {
		err = _Error(err)
		return
	}

	var ok bool
	var stringH string

	// the zettel is already mapped to a hinweis,
	// so just short circuit and return that
	if stringH, ok = ss.Read(sha.String()); ok {
		if h, err = _MakeBlindHinweis(stringH); err != nil {
			err = _Error(err)
			return
		}

		return
	}

	if h, err = hn.factory.Make(); err != nil {
		err = _Error(err)
		return
	}

	var sh _Shard

	if sh, err = hn.storeH.Shard(h.Head()); err != nil {
		err = _Error(err)
		return
	}

	if _, ok = sh.Read(h.String()); ok {
		err = _Errorf("hinweis already stored: %s", h)
		return
	}

	sh.Set(h.String(), sha.String())
	ss.Set(sha.String(), h.String())

	return
}

func (hn *hinweisen) Update(h _Hinweis, s _Sha) (err error) {
	var sh _Shard

	if sh, err = hn.storeH.Shard(h.Head()); err != nil {
		err = _Error(err)
		return
	}

	if _, err = hn.Read(h); err != nil {
		err = _Errorf("hinweis '%s' does not yet exist: %w", h, err)
		return
	}

	var ss _Shard

	if ss, err = hn.storeS.Shard(s.Head()); err != nil {
		err = _Error(err)
		return
	}

	sh.Set(h.String(), s.String())
	ss.Set(s.String(), h.String())

	return
}

func (hn hinweisen) Read(h _Hinweis) (s _Sha, err error) {
	var sh _Shard

	if sh, err = hn.storeH.Shard(h.Head()); err != nil {
		err = _Error(err)
		return
	}

	var ok bool
	var shaString string

	if shaString, ok = sh.Read(h.String()); !ok {
		err = _Errorf("hinweis '%s' does not exist", h)
		return
	}

	if err = s.Set(shaString); err != nil {
		err = _Error(err)
		return
	}

	return
}

func (hn hinweisen) ReadSha(s _Sha) (h _Hinweis, err error) {
	var ss _Shard

	if ss, err = hn.storeS.Shard(s.Head()); err != nil {
		err = _Error(err)
		return
	}

	var ok bool
	var hString string

	if hString, ok = ss.Read(s.String()); !ok {
		err = _Errorf("hinweis for sha '%s' does not exist", s)
		return
	}

	if h, err = _MakeBlindHinweis(hString); err != nil {
		err = _Error(err)
		return
	}

	return
}

func (zs *hinweisen) ReadManyStrings(args ...string) (shas []_Sha, hins []_Hinweis, err error) {
	shas = make([]_Sha, len(args))
	hins = make([]_Hinweis, len(args))

	for i, a := range args {
		var h _Hinweis

		if h, err = _MakeBlindHinweis(a); err != nil {
			err = _Error(err)
			return
		}

		var sha _Sha

		if sha, err = zs.Read(h); err != nil {
			err = _Error(err)
			return
		}

		shas[i] = sha
		hins[i] = h
	}

	return
}

func (hn *hinweisen) All() (shas []_Sha, hins []_Hinweis, err error) {
	shas = make([]_Sha, 0)
	hins = make([]_Hinweis, 0)

	var es []_Entry

	if es, err = hn.storeH.All(); err != nil {
		err = _Error(err)
		return
	}

	for _, e := range es {
		var h _Hinweis

		if h, err = _MakeBlindHinweis(e.Key); err != nil {
			err = _Error(err)
			return
		}

		var sha _Sha

		if err = sha.Set(e.Value); err != nil {
			err = _Error(err)
			return
		}

		hins = append(hins, h)
		shas = append(shas, sha)
	}

	return
}
