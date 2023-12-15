package commands

import (
	"flag"
	"io"

	"github.com/friedenberg/zit/src/alfa/angeboren"
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/id"
	"github.com/friedenberg/zit/src/charlie/age"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/sku_fmt"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/lima/bestandsaufnahme"
	"github.com/friedenberg/zit/src/oscar/umwelt"
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

	ofo := objekte_format.Options{IncludeTai: true, IncludeVerzeichnisse: true}

	bf := bestandsaufnahme.MakeAkteFormat(u.Konfig().GetStoreVersion(), ofo)

	var rc io.ReadCloser

	// setup besty reader
	{
		o := standort.FileReadOptions{
			Age:             *ag,
			Path:            c.Bestandsaufnahme,
			CompressionType: c.CompressionType,
		}

		if rc, err = standort.NewFileReader(o); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.DeferredCloser(&err, rc)
	}

	// scanner := sku_formats.MakeFormatBestandsaufnahmeScanner(
	// 	rc,
	// 	objekte_format.FormatForVersion(u.Konfig().GetStoreVersion()),
	// 	ofo,
	// )

	// for scanner.Scan() {
	// }

	besty := bestandsaufnahme.MakeAkte()

	if _, err = bf.ParseAkte(rc, besty); err != nil {
		err = errors.Wrap(err)
		return
	}

	u.Lock()
	defer u.Unlock()

	if err = u.StoreObjekten().GetBestandsaufnahmeStore().Create(besty); err != nil {
		err = errors.Wrap(err)
		return
	}

	f1 := u.StoreObjekten().GetReindexFunc()

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
	sk *sku.Transacted,
	ag *age.Age,
) (err error) {
	akteSha := sk.GetAkteSha()

	if u.Standort().HasAkte(u.Konfig().GetStoreVersion(), akteSha) {
		return
	}

	p := id.Path(akteSha, c.Akten)

	o := standort.FileReadOptions{
		Age:             *ag,
		Path:            p,
		CompressionType: c.CompressionType,
	}

	var rc sha.ReadCloser

	if rc, err = standort.NewFileReader(o); err != nil {
		if errors.IsNotExist(err) {
			err = errors.TodoRecoverable(
				"make recoverable: sku missing akte: %s",
				sku_fmt.String(sk),
			)
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	defer errors.DeferredCloser(&err, rc)

	var aw sha.WriteCloser

	if aw, err = u.Standort().AkteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, aw)

	var n int64

	if n, err = io.Copy(aw, rc); err != nil {
		errors.TodoRecoverable("%s: Sku: %s", err, sku_fmt.String(sk))
		err = nil
		return
	}

	shaRc := rc.GetShaLike()

	if !shaRc.EqualsSha(akteSha) {
		errors.TodoRecoverable(
			"sku akte mismatch: %s while akten had %s",
			sku_fmt.String(sk),
			shaRc,
		)
	}

	errors.Err().Printf("copied Akte %s (%d bytes)", akteSha, n)

	return
}
