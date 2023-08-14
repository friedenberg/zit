package umwelt

import (
	"io/ioutil"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/log"
	"github.com/friedenberg/zit/src/charlie/age"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/charlie/collections_ptr"
	"github.com/friedenberg/zit/src/charlie/file_lock"
	"github.com/friedenberg/zit/src/charlie/standort"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/erworben"
	"github.com/friedenberg/zit/src/hotel/objekte"
	"github.com/friedenberg/zit/src/hotel/objekte_store"
	"github.com/friedenberg/zit/src/india/konfig"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/kilo/organize_text"
	"github.com/friedenberg/zit/src/kilo/store_util"
	"github.com/friedenberg/zit/src/lima/store_objekten"
	"github.com/friedenberg/zit/src/mike/store_fs"
)

type Umwelt struct {
	sonnenaufgang kennung.Time

	in  *os.File
	out *os.File
	err *os.File

	inIsTty  bool
	outIsTty bool
	errIsTty bool

	standort    standort.Standort
	erworbenCli erworben.Cli
	konfig      konfig.Compiled

	storesInitialized     bool
	lock                  *file_lock.Lock
	storeUtil             store_util.StoreUtil
	storeObjekten         *store_objekten.Store
	age                   *age.Age
	storeWorkingDirectory *store_fs.Store

	zettelVerzeichnissePool schnittstellen.Pool[zettel.Transacted, *zettel.Transacted]
}

func Make(kCli erworben.Cli, options Options) (u *Umwelt, err error) {
	u = &Umwelt{
		in:                      os.Stdin,
		out:                     os.Stdout,
		err:                     os.Stderr,
		zettelVerzeichnissePool: collections.MakePool[zettel.Transacted, *zettel.Transacted](),
		erworbenCli:             kCli,
	}

	u.konfig.Reset()

	if files.IsTty(u.in) {
		u.inIsTty = true
	}

	if files.IsTty(u.out) {
		u.outIsTty = true
	}

	if files.IsTty(u.err) {
		u.errIsTty = true
	}

	err = u.Initialize(options)

	return
}

func (u *Umwelt) Reset() (err error) {
	return u.Initialize(OptionsEmpty)
}

func (u *Umwelt) Initialize(options Options) (err error) {
	if err = u.Flush(); err != nil {
		err = errors.Wrap(err)
		return
	}

	u.sonnenaufgang = kennung.Now()

	errors.TodoP4("find a better place for this")
	{
		if u.erworbenCli.Verbose {
			errors.SetVerbose()
		} else {
			log.SetOutput(ioutil.Discard)
		}

		if u.erworbenCli.Todo {
			errors.SetTodoOn()
		}

		standortOptions := standort.Options{
			BasePath: u.erworbenCli.BasePath,
		}

		if standortOptions.BasePath == "" {
			if standortOptions.BasePath, err = os.Getwd(); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		if u.standort, err = standort.Make(standortOptions); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = os.MkdirAll(
			u.standort.DirTempLocal(),
			os.ModeDir|0o755,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	{
		fa := u.standort.FileAge()

		if files.Exists(fa) {
			if u.age, err = age.MakeFromIdentityFile(fa); err != nil {
				errors.Wrap(err)
				return
			}
		} else {
			u.age = &age.Age{}
		}
	}

	{
		var k *konfig.Compiled

		if k, err = konfig.Make(
			u.standort,
			u.erworbenCli,
		); err != nil {
			if options.GetAllowKonfigReadError() {
				err = nil
			} else {
				err = errors.Wrap(err)
				return
			}
		}

		u.konfig = *k
	}

	u.konfig.ApplyPrintOptionsKonfig(u.konfig.PrintOptions)
	u.lock = file_lock.New(u.standort.DirZit("Lock"))

	// for _, rb := range u.konfig.Transacted.Objekte.Akte.Recipients {
	// 	if err = u.age.AddBech32PivYubikeyEC256(rb); err != nil {
	// 		errors.Wrap(err)
	// 		return
	// 	}
	// }

	log.Log().Printf("store version: %s", u.Konfig().GetStoreVersion())

	if u.storeUtil, err = store_util.MakeStoreUtil(
		u.lock,
		*u.age,
		u.KonfigPtr(),
		u.standort,
		objekte_format.FormatForVersion(u.Konfig().GetStoreVersion()),
	); err != nil {
		err = errors.Wrapf(err, "failed to initialize store util")
		return
	}

	if u.storeObjekten, err = store_objekten.Make(
		u.storeUtil,
		u.zettelVerzeichnissePool,
	); err != nil {
		err = errors.Wrapf(err, "failed to initialize zettel meta store")
		return
	}

	errors.Log().Print("initing checkout store")
	if u.storeWorkingDirectory, err = store_fs.New(
		u.Sonnenaufgang(),
		u.Konfig(),
		u.standort,
		u.storeObjekten,
	); err != nil {
		errors.Log().Print(err)
		err = errors.Wrap(err)
		return
	}

	errors.Log().Print("done initing checkout store")

	ptl := u.PrinterTransactedLike()

	lw := objekte_store.LogWriter[objekte.TransactedLikePtr]{
		New:       ptl,
		Updated:   ptl,
		Unchanged: ptl,
		Archived:  ptl,
	}

	u.storeObjekten.Zettel().SetLogWriter(lw)
	u.storeObjekten.Konfig().SetLogWriter(lw)
	u.storeObjekten.Typ().SetLogWriter(lw)
	u.storeObjekten.Etikett().SetLogWriter(lw)
	u.storeObjekten.Kasten().SetLogWriter(lw)

	u.storeWorkingDirectory.SetCheckedOutLogPrinter(
		u.PrinterCheckedOutLike(),
	)

	u.storesInitialized = true

	return
}

// TODO-P2 remove this
func (u Umwelt) DefaultEtiketten() (etiketten kennung.EtikettSet, err error) {
	f := collections_ptr.MakeFlagCommas[kennung.Etikett](
		collections_ptr.SetterPolicyAppend,
	)

	for _, e := range u.konfig.EtikettenToAddToNew {
		if err = f.Set(e); err != nil {
			err = errors.Wrapf(err, "Etikett: %s", e)
			err = errors.Wrapf(
				err,
				"Etiketten: %s",
				u.konfig.EtikettenToAddToNew,
			)
			return
		}
	}

	etiketten = f.GetSetPtrLike()

	return
}

func (u Umwelt) Flush() error {
	return u.age.Close()
}

func (u *Umwelt) MakeKennungIndex() kennung.Index {
	return kennung.Index{
		Etiketten: u.StoreObjekten().GetKennungIndex().GetEtikett,
	}
}

func (u *Umwelt) MakeKennungHidden() kennung.Matcher {
	h := kennung.MakeMatcherOrDoNotMatchOnEmpty()

	i := u.MakeKennungIndex()

	u.Konfig().EtikettenHidden.EachPtr(
		func(e *kennung.Etikett) (err error) {
			impl := u.Konfig().GetImplicitEtiketten(e)

			if err = impl.EachPtr(
				func(e *kennung.Etikett) (err error) {
					return h.Add(kennung.MakeMatcherContains(e, i))
				},
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = h.Add(kennung.MakeMatcherContains(e, i)); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	)

	return h
}

func (u *Umwelt) MakeKennungExpanders() (out kennung.Abbr) {
	out.Etikett.Expand = u.StoreObjekten().GetAbbrStore().Etiketten().ExpandStringString
	out.Typ.Expand = u.StoreObjekten().GetAbbrStore().Typen().ExpandStringString
	out.Kasten.Expand = u.StoreObjekten().GetAbbrStore().Kisten().ExpandStringString
	out.Hinweis.Expand = u.StoreObjekten().GetAbbrStore().Hinweis().ExpandStringString
	out.Sha.Expand = u.StoreObjekten().GetAbbrStore().Shas().ExpandStringString

	if u.Konfig().PrintOptions.PrintAbbreviatedKennungen {
		out.Etikett.Abbreviate = u.StoreObjekten().GetAbbrStore().Etiketten().Abbreviate
		out.Typ.Abbreviate = u.StoreObjekten().GetAbbrStore().Typen().Abbreviate
		out.Kasten.Abbreviate = u.StoreObjekten().GetAbbrStore().Kisten().Abbreviate
	}

	if u.Konfig().PrintOptions.PrintAbbreviatedHinweisen {
		out.Hinweis.Abbreviate = u.StoreObjekten().GetAbbrStore().Hinweis().Abbreviate
	}

	if u.Konfig().PrintOptions.PrintAbbreviatedShas {
		out.Sha.Abbreviate = u.StoreObjekten().GetAbbrStore().Shas().Abbreviate
	}

	return
}

func (u *Umwelt) MakeMetaIdSetWithExcludedHidden(
	cwd kennung.Matcher,
	dg gattungen.Set,
) kennung.MetaSet {
	if dg == nil {
		dg = gattungen.MakeSet(gattung.Zettel)
	}

	exc := u.MakeKennungHidden()

	i := u.MakeKennungIndex()

	return kennung.MakeMetaSet(
		cwd,
		u.MakeKennungExpanders(),
		exc,
		u.Konfig().FileExtensions,
		dg,
		u.Konfig(),
		i,
	)
}

func (u *Umwelt) MakeMetaIdSetWithoutExcludedHidden(
	cwd kennung.Matcher,
	dg gattungen.Set,
) kennung.MetaSet {
	if dg == nil {
		dg = gattungen.MakeSet(gattung.Zettel)
	}

	i := u.MakeKennungIndex()

	return kennung.MakeMetaSet(
		cwd,
		u.MakeKennungExpanders(),
		nil,
		u.Konfig().FileExtensions,
		dg,
		u.Konfig(),
		i,
	)
}

func (u *Umwelt) ApplyToOrganizeOptions(oo *organize_text.Options) {
	oo.Konfig = u.Konfig()
	oo.Expanders = u.MakeKennungExpanders()
}
