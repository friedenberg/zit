package umwelt

import (
	"io"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/log"
	"github.com/friedenberg/zit/src/charlie/age"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/delta/thyme"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/india/erworben"
	"github.com/friedenberg/zit/src/india/matcher"
	"github.com/friedenberg/zit/src/juliett/konfig"
	"github.com/friedenberg/zit/src/kilo/objekte_store"
	"github.com/friedenberg/zit/src/kilo/organize_text"
	"github.com/friedenberg/zit/src/mike/store_util"
	"github.com/friedenberg/zit/src/november/store_objekten"
)

type Umwelt struct {
	sonnenaufgang thyme.Time

	in  *os.File
	out *os.File
	err *os.File

	inIsTty  bool
	outIsTty bool
	errIsTty bool

	standort    standort.Standort
	erworbenCli erworben.Cli
	konfig      konfig.Compiled

	storesInitialized bool
	storeUtil         store_util.StoreUtil
	storeObjekten     *store_objekten.Store
	age               *age.Age

	matcherArchiviert matcher.Archiviert
}

func Make(kCli erworben.Cli, options Options) (u *Umwelt, err error) {
	u = &Umwelt{
		in:                os.Stdin,
		out:               os.Stdout,
		err:               os.Stderr,
		erworbenCli:       kCli,
		matcherArchiviert: matcher.MakeArchiviert(),
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

	u.sonnenaufgang = thyme.Now()

	errors.TodoP4("find a better place for this")
	{
		if u.erworbenCli.Verbose {
			errors.SetVerbose()
		} else {
			log.SetOutput(io.Discard)
		}

		if u.erworbenCli.Todo {
			errors.SetTodoOn()
		}

		standortOptions := standort.Options{
			BasePath: u.erworbenCli.BasePath,
			Debug:    u.erworbenCli.Debug,
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

	if err = u.konfig.Initialize(
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

	u.konfig.ApplyPrintOptionsKonfig(u.konfig.Akte.PrintOptions)

	// for _, rb := range u.konfig.Transacted.Objekte.Akte.Recipients {
	// 	if err = u.age.AddBech32PivYubikeyEC256(rb); err != nil {
	// 		errors.Wrap(err)
	// 		return
	// 	}
	// }

	log.Log().Printf("store version: %s", u.Konfig().GetStoreVersion())

	if u.storeUtil, err = store_util.MakeStoreUtil(
		u.Konfig(),
		u.standort,
		objekte_format.FormatForVersion(u.Konfig().GetStoreVersion()),
		u.sonnenaufgang,
	); err != nil {
		err = errors.Wrapf(err, "failed to initialize store util")
		return
	}

	if u.storeObjekten, err = store_objekten.Make(
		u.storeUtil,
	); err != nil {
		err = errors.Wrapf(err, "failed to initialize zettel meta store")
		return
	}

	errors.Log().Print("done initing checkout store")

	ptl := u.PrinterTransactedLike()

	lw := objekte_store.LogWriter{
		New:       ptl,
		Updated:   ptl,
		Unchanged: ptl,
		Archived:  ptl,
	}

	u.storeUtil.SetCheckedOutLogWriter(u.PrinterCheckedOut())
	u.storeObjekten.SetLogWriter(lw)

	u.storesInitialized = true

	return
}

func (u Umwelt) Flush() error {
	return u.age.Close()
}

func (u Umwelt) PrintMatchedArchiviertIfNecessary() {
	if !u.Konfig().PrintOptions.PrintMatchedArchiviert {
		return
	}

	c := u.GetMatcherArchiviert().Count()
	ca := u.GetMatcherArchiviert().CountArchiviert()

	if c != 0 || ca == 0 {
		return
	}

	errors.Err().Printf("%d archived objekten matched", c)
}

func (u *Umwelt) MakeKennungIndex() kennung.Index {
	return kennung.Index{}
}

func (u *Umwelt) GetMatcherArchiviert() matcher.Archiviert {
	return u.matcherArchiviert
}

func (u *Umwelt) MakeKennungExpanders() (out kennung.Abbr) {
	out.Etikett.Expand = u.StoreObjekten().GetAbbrStore().Etiketten().ExpandStringString
	out.Typ.Expand = u.StoreObjekten().GetAbbrStore().Typen().ExpandStringString
	out.Kasten.Expand = u.StoreObjekten().GetAbbrStore().Kisten().ExpandStringString
	out.Hinweis.Expand = u.StoreObjekten().GetAbbrStore().Hinweis().ExpandStringString
	out.Sha.Expand = u.StoreObjekten().GetAbbrStore().Shas().ExpandStringString

	out.Etikett.Abbreviate = u.StoreObjekten().GetAbbrStore().Etiketten().Abbreviate
	out.Typ.Abbreviate = u.StoreObjekten().GetAbbrStore().Typen().Abbreviate
	out.Kasten.Abbreviate = u.StoreObjekten().GetAbbrStore().Kisten().Abbreviate
	out.Hinweis.Abbreviate = u.StoreObjekten().GetAbbrStore().Hinweis().Abbreviate
	out.Sha.Abbreviate = u.StoreObjekten().GetAbbrStore().Shas().Abbreviate

	return
}

func (u *Umwelt) MakeMetaIdSetWithExcludedHidden(
	dg gattungen.Set,
) matcher.Query {
	if dg == nil {
		dg = gattungen.MakeSet(gattung.Zettel)
	}

	exc := u.GetMatcherArchiviert()

	i := u.MakeKennungIndex()

	return matcher.MakeQuery(
		u.Konfig(),
		u.StoreUtil().GetCwdFiles(),
		u.MakeKennungExpanders(),
		exc,
		u.Konfig().FileExtensions,
		dg,
		i,
	)
}

func (u *Umwelt) MakeQueryAll() matcher.Query {
	i := u.MakeKennungIndex()

	return matcher.MakeQueryAll(
		u.Konfig(),
		u.StoreUtil().GetCwdFiles(),
		u.MakeKennungExpanders(),
		nil,
		u.Konfig().FileExtensions,
		i,
	)
}

func (u *Umwelt) MakeMetaIdSetWithoutExcludedHidden(
	dg gattungen.Set,
) matcher.Query {
	if dg == nil {
		dg = gattungen.MakeSet(gattung.Zettel)
	}

	i := u.MakeKennungIndex()

	return matcher.MakeQuery(
		u.Konfig(),
		u.StoreUtil().GetCwdFiles(),
		u.MakeKennungExpanders(),
		nil,
		u.Konfig().FileExtensions,
		dg,
		i,
	)
}

func (u *Umwelt) ApplyToOrganizeOptions(oo *organize_text.Options) {
	oo.Konfig = u.Konfig()
	oo.Expanders = u.MakeKennungExpanders()
}