package zettels

import (
	"path"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/alfa/logz"
	"github.com/friedenberg/zit/alfa/stdprinter"
	"github.com/friedenberg/zit/bravo/id"
	"github.com/friedenberg/zit/bravo/sha"
	"github.com/friedenberg/zit/charlie/age"
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/charlie/konfig"
	"github.com/friedenberg/zit/delta/umwelt"
	"github.com/friedenberg/zit/echo/sharded_store"
	"github.com/friedenberg/zit/echo/zettel"
	"github.com/friedenberg/zit/foxtrot/hinweisen"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
)

type Zettels interface {
	ReadZettel(sha sha.Sha) (z stored_zettel.Stored, err error)
	Read(id id.Id) (z stored_zettel.Named, err error)
	Create(zettel.Zettel) (z stored_zettel.Named, err error)
	CreateWithHinweis(zettel.Zettel, hinweis.Hinweis) (z stored_zettel.Named, err error)
	Update(z stored_zettel.Named) (stored stored_zettel.Named, err error)
	Revert(h hinweis.Hinweis) (named stored_zettel.Named, err error)
	UpdateNoKinder(z stored_zettel.Named) (err error)

	Delete(id id.Id) (zettel stored_zettel.Named, err error)

	Flush() error

	Hinweisen() hinweisen.Hinweisen
	Age() age.Age
	Umwelt() *umwelt.Umwelt

	Konfig() konfig.Konfig

	zettel.AkteReaderFactory
	zettel.AkteWriterFactory
}

type zettels struct {
	umwelt    *umwelt.Umwelt
	store     sharded_store.Store
	basePath  string
	age       age.Age
	hinweisen hinweisen.Hinweisen
}

func New(u *umwelt.Umwelt, age age.Age) (s *zettels, err error) {
	s = &zettels{
		umwelt:   u,
		basePath: u.DirZit(),
		age:      age,
	}

	// if s.hinweisen, err = hinweisen.New(age, s.basePath); err != nil {
	// 	err = errors.Error(err)
	// 	return
	// }

	zp := path.Join(s.basePath, "Objekte", "Zettel")

	if s.store, err = sharded_store.NewStore(zp, s); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (zs *zettels) Umwelt() *umwelt.Umwelt {
	return zs.umwelt
}

func (zs *zettels) Age() age.Age {
	return zs.age
}

func (zs *zettels) Hinweisen() hinweisen.Hinweisen {
	return zs.hinweisen
}

func (zs *zettels) Konfig() konfig.Konfig {
	return zs.umwelt.Konfig
}

//TODO-P0,D4 make flushing atomic
func (zs *zettels) Flush() (err error) {
	// logz.Print("flushing zettels")
	// if err = zs.store.Flush(); err != nil {
	// 	err = errors.Error(err)
	// 	return
	// }

	// logz.Print("flushing hinweisen")
	// if err = zs.Hinweisen().Flush(); err != nil {
	// 	err = errors.Error(err)
	// 	return
	// }

	// logz.Print("done flushing")

	return
}

func (zs zettels) NewShard(p string, id string) (s sharded_store.Shard, err error) {
	if s, err = sharded_store.NewShard(path.Join(p, id), zs.age, sharded_store.ShardGeneric{}); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

//TODO-P2,D2 move to store_with_lock
func (zs zettels) CreateWithHinweis(in zettel.Zettel, h hinweis.Hinweis) (z stored_zettel.Named, err error) {
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

	var named stored_zettel.Named

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
func (zs zettels) Create(in zettel.Zettel) (z stored_zettel.Named, err error) {
	z.Stored.Zettel = in

	if z.Sha, err = zs.storeBaseZettel(z.Stored); err != nil {
		err = errors.Error(err)
		return
	}

	var existing stored_zettel.Named

	if existing, err = zs.Read(z.Sha); err == nil {
		z = existing
		return
	}

	if z.Hinweis, err = zs.hinweisen.StoreNew(z.Sha); err != nil {
		err = errors.Error(err)
		return
	}

	var named stored_zettel.Named

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
func (zs zettels) Update(in stored_zettel.Named) (z stored_zettel.Named, err error) {
	if in.Zettel.IsEmpty() {
		err = errors.Normal(errors.Errorf("zettel is empty"))
		return
	}

	var mutter stored_zettel.Named

	if mutter, err = zs.Read(in.Hinweis); err != nil {
		err = errors.Error(err)
		return
	}

	z = in

	if z.Zettel.Equals(mutter.Zettel) {
		stdprinter.Errf("[%s %s] (unchanged)\n", z.Hinweis, z.Sha)
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

	stdprinter.Errf("[%s %s] (updated)\n", z.Hinweis, z.Sha)

	return
}

//TODO-P1,D3 move to store_with_lock
func (zs zettels) Revert(h hinweis.Hinweis) (named stored_zettel.Named, err error) {
	if named, err = zs.Read(h); err != nil {
		err = errors.Error(err)
		return
	}

	if named.Mutter.IsNull() {
		err = errors.Normal(errors.Errorf("cannot revert %s as it has no Mutter", h))
		return
	}

	var mutter stored_zettel.Named

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

func (zs zettels) UpdateNoKinder(zettel stored_zettel.Named) (err error) {
	if err = zs.update(zettel.Stored); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (zs zettels) Delete(id id.Id) (zettel stored_zettel.Named, err error) {
	if zettel, err = zs.readNamedZettel(id); err != nil {
		err = errors.Error(err)
		return
	}

	var s sharded_store.Shard

	if s, err = zs.store.Shard(zettel.Sha.Head()); err != nil {
		err = errors.Error(err)
		return
	}

	s.Remove(zettel.Sha.String())

	return
}

func (zs zettels) ReadZettel(sha sha.Sha) (z stored_zettel.Stored, err error) {
	if z, err = zs.readStoredZettel(sha); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (zs zettels) Read(id id.Id) (sz stored_zettel.Named, err error) {
	if sz, err = zs.readNamedZettel(id); err != nil {
		err = errors.Error(err)
		return
	}

	return
}
