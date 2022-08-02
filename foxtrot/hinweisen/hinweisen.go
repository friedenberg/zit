package hinweisen

import (
	"path"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/alfa/logz"
	"github.com/friedenberg/zit/bravo/sha"
	"github.com/friedenberg/zit/charlie/age"
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/echo/sharded_store"
)

type Hinweisen interface {
	Read(h hinweis.Hinweis) (sha sha.Sha, err error)
	ReadSha(s sha.Sha) (h hinweis.Hinweis, err error)
	ReadString(s string) (sha sha.Sha, hin hinweis.Hinweis, err error)
	ReadManyStrings(args ...string) (shas []sha.Sha, hins []hinweis.Hinweis, err error)
	All() (shas []sha.Sha, hins []hinweis.Hinweis, err error)
	StoreNew(sha sha.Sha) (h hinweis.Hinweis, err error)
	StoreExisting(h hinweis.Hinweis, sha sha.Sha) (err error)
	Update(h hinweis.Hinweis, s sha.Sha) (err error)
	Flush() error
	Factory() *factory
}

type hinweisen struct {
	basePath string
	storeH   sharded_store.Store
	storeS   sharded_store.Store
	factory  *factory
}

func New(age age.Age, basePath string) (s *hinweisen, err error) {
	s = &hinweisen{
		basePath: basePath,
	}

	if s.factory, err = newFactory(basePath); err != nil {
		err = errors.Error(err)
		return
	}

	if s.storeS, err = sharded_store.NewStore(path.Join(basePath, "Zettel-Hinweis"), s); err != nil {
		err = errors.Error(err)
		return
	}

	if s.storeH, err = sharded_store.NewStore(path.Join(basePath, "Hinweis"), s); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (hn hinweisen) Factory() *factory {
	return hn.factory
}

func (hn hinweisen) NewShard(p string, id string) (s sharded_store.Shard, err error) {
	if s, err = sharded_store.NewShard(path.Join(p, id), nil, &sharded_store.ShardGeneric{}); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (zs *hinweisen) Flush() (err error) {
	if err = zs.storeH.Flush(); err != nil {
		err = errors.Error(err)
		return
	}

	if err = zs.storeS.Flush(); err != nil {
		err = errors.Error(err)
		return
	}

	if err = zs.factory.Flush(); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (hn *hinweisen) StoreNew(sha sha.Sha) (h hinweis.Hinweis, err error) {
	logz.Print("storing new")
	logz.PrintDebug(hn.factory)
	if h, err = hn.factory.Make(); err != nil {
		logz.Print("failed")
		err = errors.Error(err)
		return
	}
	logz.Print("succeeded")

	err = hn.StoreExisting(h, sha)

	return
}

func (hn *hinweisen) StoreExisting(h hinweis.Hinweis, sha sha.Sha) (err error) {
	var ss sharded_store.Shard

	if ss, err = hn.storeS.Shard(sha.Head()); err != nil {
		err = errors.Error(err)
		return
	}

	var ok bool
	var stringH string

	// the zettel is already mapped to a hinweis,
	// so just short circuit and return that
	if stringH, ok = ss.Read(sha.String()); ok {
		if h, err = hinweis.MakeBlindHinweis(stringH); err != nil {
			err = errors.Error(err)
			return
		}

		return
	}

	var sh sharded_store.Shard

	logz.PrintDebug(h)
	logz.PrintDebug(h.Head())
	logz.PrintDebug("wow")
	if sh, err = hn.storeH.Shard(h.Head()); err != nil {
		err = errors.Error(err)
		return
	}

	if _, ok = sh.Read(h.String()); ok {
		err = errors.Errorf("hinweis already stored: %s", h)
		return
	}

	sh.Set(h.String(), sha.String())
	ss.Set(sha.String(), h.String())

	return
}

func (hn *hinweisen) Update(h hinweis.Hinweis, s sha.Sha) (err error) {
	var sh sharded_store.Shard

	if sh, err = hn.storeH.Shard(h.Head()); err != nil {
		err = errors.Error(err)
		return
	}

	if _, err = hn.Read(h); err != nil {
		err = errors.Errorf("hinweis '%s' does not yet exist: %s", h, err)
		return
	}

	var ss sharded_store.Shard

	if ss, err = hn.storeS.Shard(s.Head()); err != nil {
		err = errors.Error(err)
		return
	}

	sh.Set(h.String(), s.String())
	ss.Set(s.String(), h.String())

	return
}

func (hn hinweisen) Read(h hinweis.Hinweis) (s sha.Sha, err error) {
	var sh sharded_store.Shard

	if sh, err = hn.storeH.Shard(h.Head()); err != nil {
		err = errors.Error(err)
		return
	}

	var ok bool
	var shaString string

	if shaString, ok = sh.Read(h.String()); !ok {
		err = errors.Wrapped(ErrDoesNotExist, "%s", h)
		return
	}

	if err = s.Set(shaString); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (hn hinweisen) ReadSha(s sha.Sha) (h hinweis.Hinweis, err error) {
	var ss sharded_store.Shard

	if ss, err = hn.storeS.Shard(s.Head()); err != nil {
		err = errors.Error(err)
		return
	}

	var ok bool
	var hString string

	if hString, ok = ss.Read(s.String()); !ok {
		err = errors.Errorf("hinweis for sha '%s' does not exist", s)
		return
	}

	if h, err = hinweis.MakeBlindHinweis(hString); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (zs *hinweisen) ReadString(s string) (sha sha.Sha, hin hinweis.Hinweis, err error) {
	if hin, err = hinweis.MakeBlindHinweis(s); err != nil {
		err = errors.Error(err)
		return
	}

	if sha, err = zs.Read(hin); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (zs *hinweisen) ReadManyStrings(args ...string) (shas []sha.Sha, hins []hinweis.Hinweis, err error) {
	shas = make([]sha.Sha, len(args))
	hins = make([]hinweis.Hinweis, len(args))

	for i, a := range args {
		var h hinweis.Hinweis

		if h, err = hinweis.MakeBlindHinweis(a); err != nil {
			err = errors.Error(err)
			return
		}

		var sha sha.Sha

		if sha, err = zs.Read(h); err != nil {
			err = errors.Error(err)
			return
		}

		shas[i] = sha
		hins[i] = h
	}

	return
}

func (hn *hinweisen) All() (shas []sha.Sha, hins []hinweis.Hinweis, err error) {
	shas = make([]sha.Sha, 0)
	hins = make([]hinweis.Hinweis, 0)

	var es []sharded_store.Entry

	if es, err = hn.storeH.All(); err != nil {
		err = errors.Error(err)
		return
	}

	for _, e := range es {
		var h hinweis.Hinweis

		if h, err = hinweis.MakeBlindHinweis(e.Key); err != nil {
			err = errors.Error(err)
			return
		}

		var sha sha.Sha

		if err = sha.Set(e.Value); err != nil {
			err = errors.Error(err)
			return
		}

		hins = append(hins, h)
		shas = append(shas, sha)
	}

	return
}
