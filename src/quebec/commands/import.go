package commands

import (
	"flag"
	"fmt"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/id"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/age"
	"github.com/friedenberg/zit/src/delta/age_io"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/kennung_index"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/hotel/sku_formats"
	"github.com/friedenberg/zit/src/india/bestandsaufnahme"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type Import struct {
	Bestandsaufnahme string
	Akten            string
	AgeIdentity      string
}

func init() {
	registerCommand(
		"import",
		func(f *flag.FlagSet) Command {
			c := &Import{}

			f.StringVar(&c.Bestandsaufnahme, "bestandsaufnahme", "", "")
			f.StringVar(&c.Akten, "akten", "", "")
			f.StringVar(&c.AgeIdentity, "age-identity", "", "")

			return c
		},
	)
}

func (c Import) Run(u *umwelt.Umwelt, args ...string) (err error) {
	if c.Bestandsaufnahme == "" {
		err = errors.Errorf("empty Bestandsaufnahme")
		return
	}

	if c.Akten == "" {
		err = errors.Errorf("empty Akten")
		return
	}

	var ag *age.Age

	if ag, err = age.MakeFromIdentity(c.AgeIdentity); err != nil {
		err = errors.Wrapf(err, "age-identity: %q", c.AgeIdentity)
		return
	}

	bf := bestandsaufnahme.MakeAkteFormat(
		u.Konfig().GetStoreVersion(),
		u.StoreObjekten(),
	)

	var rc io.ReadCloser

	// setup besty reader
	{
		o := age_io.FileReadOptions{
			Age:  *ag,
			Path: c.Bestandsaufnahme,
		}

		if rc, err = age_io.NewFileReader(o); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.DeferredCloser(&err, rc)
	}

	var besty bestandsaufnahme.Akte
	besty.Reset()

	if _, err = bf.ParseAkte(rc, &besty); err != nil {
		err = errors.Wrap(err)
		return
	}

	u.Lock()
	defer u.Unlock()

	if err = u.StoreObjekten().GetBestandsaufnahmeStore().Create(&besty); err != nil {
		err = errors.Wrap(err)
		return
	}

	var ti kennung_index.KennungIndex[kennung.Typ, *kennung.Typ]

	if ti, err = u.StoreObjekten().GetTypenIndex(); err != nil {
		err = errors.Wrap(err)
		return
	}

	f1 := u.StoreObjekten().GetReindexFunc(ti)

	for {
		sk, ok := besty.Skus.Pop()

		if !ok {
			break
		}

		if err = c.importAkteIfNecessary(u, sk, ag); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = f1(sk); err != nil {
			err = errors.Wrapf(err, "Sku: %s", sk)
			return
		}
	}

	return
}

func (c Import) importAkteIfNecessary(
	u *umwelt.Umwelt,
	sk sku.SkuLike,
	ag *age.Age,
) (err error) {
	akteSha := sk.GetAkteSha()

	if u.Standort().HasAkte(u.Konfig().GetStoreVersion(), akteSha) {
		return
	}

	p := id.Path(akteSha, c.Akten)

	o := age_io.FileReadOptions{
		Age:  *ag,
		Path: p,
	}

	var rc sha.ReadCloser

	if rc, err = age_io.NewFileReader(o); err != nil {
		if errors.IsNotExist(err) {
			err = errors.Todo(
				fmt.Sprintf(
					"make recoverable: sku missing akte: %s",
					sku_formats.String(sk),
				),
			)
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	defer errors.DeferredCloser(&err, rc)

	var aw sha.WriteCloser

	if aw, err = u.StoreObjekten().AkteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, aw)

	var n int64

	if n, err = io.Copy(aw, rc); err != nil {
		err = errors.Wrap(err)
		return
	}

	shaRc := rc.GetShaLike()

	if !shaRc.EqualsSha(akteSha) {
		err = errors.Errorf("sku had sha %s while akten had %s", akteSha, shaRc)
		return
	}

	errors.Err().Printf("copied Akte %s (%d bytes)", akteSha, n)

	return
}
