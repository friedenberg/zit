package commands

import (
	"flag"
	"io"

	"code.linenisgreat.com/zit/src/alfa/angeboren"
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/id"
	"code.linenisgreat.com/zit/src/charlie/age"
	"code.linenisgreat.com/zit/src/charlie/sha"
	"code.linenisgreat.com/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/src/delta/standort"
	"code.linenisgreat.com/zit/src/golf/objekte_format"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/kilo/zettel"
	"code.linenisgreat.com/zit/src/lima/bestandsaufnahme"
	"code.linenisgreat.com/zit/src/mike/store"
	"code.linenisgreat.com/zit/src/oscar/umwelt"
)

type Import struct {
	Bestandsaufnahme string
	Akten            string
	AgeIdentity      age.Identity
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
			f.Var(&c.AgeIdentity, "age-identity", "")
			c.CompressionType.AddToFlagSet(f)

			c.Proto.AddToFlagSet(f)

			return c
		},
	)
}

func (c Import) Run(u *umwelt.Umwelt, args ...string) (err error) {
	hasConflicts := false

	if c.Bestandsaufnahme == "" {
		err = errors.Errorf("empty Bestandsaufnahme")
		return
	}

	if c.Akten == "" {
		err = errors.Errorf("empty Akten")
		return
	}

	var ag age.Age

	if err = ag.AddIdentity(c.AgeIdentity); err != nil {
		err = errors.Wrapf(err, "age-identity: %q", &c.AgeIdentity)
		return
	}

	coPrinter := u.PrinterCheckedOut()

	ofo := objekte_format.Options{Tai: true, Verzeichnisse: true}

	bf := bestandsaufnahme.MakeAkteFormat(u.Konfig().GetStoreVersion(), ofo)

	var rc io.ReadCloser

	// setup besty reader
	{
		o := standort.FileReadOptions{
			Age:             &ag,
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

	var co *sku.CheckedOut

	for {
		sk, ok := besty.Skus.Pop()

		if !ok {
			break
		}

		if co, err = u.Store().Import(
			sk,
		); err != nil {
			err = errors.Wrapf(err, "Sku: %s", sk)
			return
		}

		if co.State == checked_out_state.StateError {
			if co.Error == store.ErrNeedsMerge {
				hasConflicts = true
			}

			if err = coPrinter(co); err != nil {
				err = errors.Wrap(err)
				return
			}

			continue
		}

		if err = c.importAkteIfNecessary(u, co, &ag, coPrinter); err != nil {
			err = errors.Wrap(err)
			return
		}

		if co.State == checked_out_state.StateError {
			if err = coPrinter(co); err != nil {
				err = errors.Wrap(err)
				return
			}

			continue
		}
	}

	if hasConflicts {
		err = store.ErrNeedsMerge
	}

	return
}

func (c Import) importAkteIfNecessary(
	u *umwelt.Umwelt,
	co *sku.CheckedOut,
	ag *age.Age,
	coErrPrinter schnittstellen.FuncIter[*sku.CheckedOut],
) (err error) {
	akteSha := co.External.GetAkteSha()

	if u.Standort().HasAkte(u.Konfig().GetStoreVersion(), akteSha) {
		return
	}

	p := id.Path(akteSha, c.Akten)

	o := standort.FileReadOptions{
		Age:             ag,
		Path:            p,
		CompressionType: c.CompressionType,
	}

	var rc sha.ReadCloser

	if rc, err = standort.NewFileReader(o); err != nil {
		if errors.IsNotExist(err) {
			co.SetError(errors.New("akte missing"))
			err = coErrPrinter(co)
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
		co.SetError(errors.New("akte copy failed"))
		err = coErrPrinter(co)
		return
	}

	shaRc := rc.GetShaLike()

	if !shaRc.EqualsSha(akteSha) {
		co.SetError(errors.New("akte sha mismatch"))
		err = coErrPrinter(co)
		errors.TodoRecoverable(
			"sku akte mismatch: %s while akten had %s",
			co.Internal.GetAkteSha(),
			shaRc,
		)
	}

	errors.Err().Printf("copied Akte %s (%d bytes)", akteSha, n)

	return
}
