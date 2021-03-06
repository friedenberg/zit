package zettels

import (
	"path"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/alfa/logz"
)

type Zettels interface {
	AllInChain(id _Id) (c Chain, err error)
	All() (map[string]_NamedZettel, error)
	Query(_NamedZettelFilter) (map[string]_NamedZettel, error)

	ReadZettel(sha _Sha) (z _StoredZettel, err error)
	Read(id _Id) (z _NamedZettel, err error)
	Create(_Zettel) (z _NamedZettel, err error)
	CreateWithHinweis(_Zettel, _Hinweis) (z _NamedZettel, err error)
	Update(z _NamedZettel) (stored _NamedZettel, err error)
	Revert(h _Hinweis) (named _NamedZettel, err error)
	UpdateNoKinder(z _NamedZettel) (err error)

	//TODO move to user_ops
	Checkout(options CheckinOptions, args ...string) (czs []_ZettelCheckedOut, err error)

	Delete(id _Id) (zettel _NamedZettel, err error)

	Flush() error

	Etiketten() _Etiketten
	Hinweisen() _Hinweisen
	Age() _Age

	Konfig() _Konfig

	_AkteReaderFactory
	_AkteWriterFactory
}

type zettels struct {
	umwelt    *_Umwelt
	store     _Store
	basePath  string
	age       _Age
	etiketten _Etiketten
	hinweisen _Hinweisen
}

func New(u *_Umwelt, age _Age) (s *zettels, err error) {
	s = &zettels{
		umwelt:   u,
		basePath: u.DirZit(),
		age:      age,
	}

	if s.hinweisen, err = _NewHinweisen(age, s.basePath); err != nil {
		err = errors.Error(err)
		return
	}

	if s.etiketten, err = _NewEtiketten(u.Konfig, age, s.basePath); err != nil {
		err = errors.Error(err)
		return
	}

	zp := path.Join(s.basePath, "Objekte", "Zettel")

	if s.store, err = _NewStore(zp, s); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (zs *zettels) Age() _Age {
	return zs.age
}

func (zs *zettels) Hinweisen() _Hinweisen {
	return zs.hinweisen
}

func (zs *zettels) Etiketten() _Etiketten {
	return zs.etiketten
}

func (zs *zettels) Konfig() _Konfig {
	return zs.umwelt.Konfig
}

//TODO-P0,D4 make flushing atomic
func (zs *zettels) Flush() (err error) {
	logz.Print("flushing zettels")
	if err = zs.store.Flush(); err != nil {
		err = errors.Error(err)
		return
	}

	logz.Print("flushing hinweisen")
	if err = zs.Hinweisen().Flush(); err != nil {
		err = errors.Error(err)
		return
	}

	logz.Print("flushing etiketten")
	if err = zs.Etiketten().Flush(); err != nil {
		err = errors.Error(err)
		return
	}

	logz.Print("done flushing")

	return
}

func (zs zettels) NewShard(p string, id string) (s _Shard, err error) {
	if s, err = _NewShard(path.Join(p, id), zs.age, _ShardGeneric{}); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

//TODO-P2,D2 move to store_with_lock
func (zs zettels) CreateWithHinweis(in _Zettel, h _Hinweis) (z _NamedZettel, err error) {
	if in.IsEmpty() {
		err = errors.Normal(errors.Errorf("zettel is empty"))
		return
	}

	z.Stored.Zettel = in

	if z.Sha, err = zs.storeBaseZettel(z.Stored); err != nil {
		err = errors.Error(err)
		return
	}

	if err = zs.hinweisen.StoreExisting(h, z.Sha); err != nil {
		err = errors.Error(err)
		return
	}

	z.Hinweis = h

	var named _NamedZettel

	if named, err = zs.Read(z.Sha); err != nil {
		err = errors.Error(err)
		return
	}

	logz.Print(z)
	logz.Print(named)

	if !z.Equals(named) {
		err = errors.Errorf(
			"stored zettel doesn't equal stored zettel:\n%s\n%s",
			z,
			named,
		)

		return
	}

	return
}

//TODO-P1,D2 move to store_with_lock
func (zs zettels) Create(in _Zettel) (z _NamedZettel, err error) {
	if in.IsEmpty() {
		err = errors.Normal(errors.Errorf("zettel is empty"))
		return
	}

	z.Stored.Zettel = in

	if z.Sha, err = zs.storeBaseZettel(z.Stored); err != nil {
		err = errors.Error(err)
		return
	}

	var existing _NamedZettel

	if existing, err = zs.Read(z.Sha); err == nil {
		z = existing
		return
	}

	if z.Hinweis, err = zs.hinweisen.StoreNew(z.Sha); err != nil {
		err = errors.Error(err)
		return
	}

	var named _NamedZettel

	if named, err = zs.Read(z.Sha); err != nil {
		err = errors.Error(err)
		return
	}

	if !z.Equals(named) {
		err = errors.Errorf(
			"stored zettel doesn't equal stored zettel:\n%s\n%s",
			z,
			named,
		)

		return
	}

	return
}

//TODO-P1,D3 move to store_with_lock
func (zs zettels) Update(in _NamedZettel) (z _NamedZettel, err error) {
	if in.Zettel.IsEmpty() {
		err = errors.Normal(errors.Errorf("zettel is empty"))
		return
	}

	var mutter _NamedZettel

	if mutter, err = zs.Read(in.Hinweis); err != nil {
		err = errors.Error(err)
		return
	}

	z = in

	if z.Zettel.Equals(mutter.Zettel) {
		_Errf("[%s %s] (unchanged)\n", z.Hinweis, z.Sha)
		return
	}

	z.Mutter = mutter.Sha

	if z.Mutter.IsNull() {
		err = errors.Errorf("mutter cannot be null")
		return
	}

	logz.Printf("updating zettel %s %s", z.Zettel, z.Sha)
	if z.Sha, err = zs.storeBaseZettel(z.Stored); err != nil {
		err = errors.Error(err)
		return
	}

	logz.Printf("updating mutter %s %s", z.Mutter, z.Sha)
	if err = zs.updateMutterIfNecessary(z.Mutter, z.Sha); err != nil {
		err = errors.Error(err)
		return
	}

	logz.Printf("updating %s %s", z.Hinweis, z.Sha)
	if err = zs.hinweisen.Update(z.Hinweis, z.Sha); err != nil {
		err = errors.Error(err)
		return
	}

	_Errf("[%s %s] (updated)\n", z.Hinweis, z.Sha)

	return
}

//TODO-P1,D3 move to store_with_lock
func (zs zettels) Revert(h _Hinweis) (named _NamedZettel, err error) {
	if named, err = zs.Read(h); err != nil {
		err = errors.Error(err)
		return
	}

	if named.Mutter.IsNull() {
		err = errors.Normal(errors.Errorf("cannot revert %s as it has no Mutter", h))
		return
	}

	var mutter _NamedZettel

	if mutter, err = zs.Read(named.Mutter); err != nil {
		err = errors.Error(err)
		return
	}

	named.Zettel = mutter.Zettel

	if named, err = zs.Update(named); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (zs zettels) UpdateNoKinder(zettel _NamedZettel) (err error) {
	if err = zs.update(zettel.Stored); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (zs zettels) Delete(id _Id) (zettel _NamedZettel, err error) {
	if zettel, err = zs.readNamedZettel(id); err != nil {
		err = errors.Error(err)
		return
	}

	var s _Shard

	if s, err = zs.store.Shard(zettel.Sha.Head()); err != nil {
		err = errors.Error(err)
		return
	}

	s.Remove(zettel.Sha.String())

	return
}

func (zs zettels) ReadZettel(sha _Sha) (z _StoredZettel, err error) {
	if z, err = zs.readStoredZettel(sha); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (zs zettels) Read(id _Id) (sz _NamedZettel, err error) {
	if sz, err = zs.readNamedZettel(id); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

//TODO move to store_with_lock
func (zs zettels) All() (ns map[string]_NamedZettel, err error) {
	ns = make(map[string]_NamedZettel)

	var es []_Entry

	if es, err = zs.store.All(); err != nil {
		err = errors.Error(err)
		return
	}

OUTER:
	for _, e := range es {
		var sha _Sha

		if err = sha.Set(e.Key); err != nil {
			err = errors.Error(err)
			return
		}

		var named _NamedZettel

		if named, err = zs.Read(sha); err != nil {
			err = errors.Error(err)
			return
		}

		logz.Print(named)

		if !named.Kinder.IsNull() {
			continue OUTER
		}

		prefixes := named.Zettel.Etiketten.Expanded(_EtikettExpanderRight{})

		logz.Print(zs.umwelt.Konfig.Tags)
		logz.Print(prefixes)
		logz.Print(named.Zettel.Etiketten)

	INNER:
		for tn, tv := range zs.umwelt.Konfig.Tags {
			if !tv.Hide {
				logz.Print("not hidden, checking next tag")
				continue INNER
			}

			logz.Print("checking for hide matches")
			logz.Print(tn)
			logz.Print(prefixes)
			logz.Print(prefixes.ContainsString(tn))

			if prefixes.ContainsString(tn) {
				logz.Printf("hiding %s due to %s", named.Hinweis, tn)
				continue OUTER
			}
		}

		if otherZettel, ok := ns[named.Hinweis.String()]; ok {
			err = errors.Errorf(
				"two separate zettels with hinweis: %s:\n%s\n%s",
				named.Hinweis,
				otherZettel.Sha,
				named.Sha,
			)
			return
		}

		ns[named.Hinweis.String()] = named
	}

	return
}

//TODO swap query and all methods for performance reasons
func (zs zettels) Query(filter _NamedZettelFilter) (ns map[string]_NamedZettel, err error) {
	var ns1 map[string]_NamedZettel

	if ns1, err = zs.All(); err != nil {
		err = errors.Error(err)
		return
	}

	if filter == nil {
		ns = ns1
		return
	}

	ns = make(map[string]_NamedZettel)

	for n, z := range ns1 {
		if filter.IncludeNamedZettel(z) {
			ns[n] = z
		}
	}

	return
}
