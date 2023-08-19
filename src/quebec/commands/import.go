package commands

import (
	"flag"
	"io"

	"github.com/friedenberg/zit/src/alfa/angeboren"
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/id"
	"github.com/friedenberg/zit/src/charlie/age"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/echo/age_io"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/kennung_index"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/hotel/sku_formats"
	"github.com/friedenberg/zit/src/india/bestandsaufnahme"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type Import struct {
	Bestandsaufnahme string
	Akten            string
	AgeIdentity      string
	CompressionType  angeboren.CompressionType
	Proto            zettel.ProtoZettel
}

func init() {
	registerCommand(
		"import",
		func(f *flag.FlagSet) Command {
			c := &Import{
				Proto:           zettel.MakeEmptyProtoZettel(),
				CompressionType: angeboren.CompressionTypeDefault,
			}

			f.StringVar(&c.Bestandsaufnahme, "bestandsaufnahme", "", "")
			f.StringVar(&c.Akten, "akten", "", "")
			f.StringVar(&c.AgeIdentity, "age-identity", "", "")
			c.CompressionType.AddToFlagSet(f)

			c.Proto.AddToFlagSet(f)

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
			Age:             *ag,
			Path:            c.Bestandsaufnahme,
			CompressionType: c.CompressionType,
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
		Age:             *ag,
		Path:            p,
		CompressionType: c.CompressionType,
	}

	var rc sha.ReadCloser

	if rc, err = age_io.NewFileReader(o); err != nil {
		if errors.IsNotExist(err) {
			err = errors.TodoRecoverable(
				"make recoverable: sku missing akte: %s",
				sku_formats.String(sk),
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
		errors.TodoRecoverable("%s: Sku: %s", err, sku_formats.String(sk))
		err = nil
		return
	}

	shaRc := rc.GetShaLike()

	if !shaRc.EqualsSha(akteSha) {
		errors.TodoRecoverable(
			"sku akte mismatch: %s while akten had %s",
			sku_formats.String(sk),
			shaRc,
		)
	}

	errors.Err().Printf("copied Akte %s (%d bytes)", akteSha, n)

	return
}
