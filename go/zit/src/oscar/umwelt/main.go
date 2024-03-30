package umwelt

import (
	"io"
	"os"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/bravo/files"
	"code.linenisgreat.com/zit/src/bravo/log"
	"code.linenisgreat.com/zit/src/charlie/age"
	"code.linenisgreat.com/zit/src/charlie/gattung"
	"code.linenisgreat.com/zit/src/delta/standort"
	"code.linenisgreat.com/zit/src/delta/thyme"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/golf/objekte_format"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/india/erworben"
	"code.linenisgreat.com/zit/src/india/query"
	"code.linenisgreat.com/zit/src/juliett/chrome"
	"code.linenisgreat.com/zit/src/juliett/konfig"
	"code.linenisgreat.com/zit/src/kilo/objekte_store"
	"code.linenisgreat.com/zit/src/lima/organize_text"
	"code.linenisgreat.com/zit/src/mike/store"
)

type VirtualStore = query.VirtualStoreInitable

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
	store             store.Store
	virtualStores     map[string]*VirtualStore
	age               *age.Age

	matcherArchiviert query.Archiviert
}

func Make(kCli erworben.Cli, options Options) (u *Umwelt, err error) {
	u = &Umwelt{
		in:                os.Stdin,
		out:               os.Stdout,
		err:               os.Stderr,
		erworbenCli:       kCli,
		matcherArchiviert: query.MakeArchiviert(),
		virtualStores:     make(map[string]*VirtualStore),
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

	if err = u.store.Initialize(
		u.Konfig(),
		u.standort,
		objekte_format.FormatForVersion(u.Konfig().GetStoreVersion()),
		u.sonnenaufgang,
	); err != nil {
		err = errors.Wrapf(err, "failed to initialize store util")
		return
	}

	errors.Log().Print("done initing checkout store")

	ptl := u.PrinterTransactedLike()

	lw := objekte_store.LogWriter{
		New:       ptl,
		Updated:   ptl,
		Unchanged: func(sk *sku.Transacted) (err error) {
      if !u.konfig.PrintOptions.PrintUnchanged {
        return
      }

      return ptl(sk)
    },
		Archived:  ptl,
	}

	u.store.SetCheckedOutLogWriter(u.PrinterCheckedOut())
	u.store.SetLogWriter(lw)

	u.storesInitialized = true

	u.virtualStores["%chrome"] = &query.VirtualStoreInitable{
		VirtualStore: chrome.MakeChrome(u.Standort()),
	}

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

func (u *Umwelt) GetMatcherArchiviert() query.Archiviert {
	return u.matcherArchiviert
}

func (u *Umwelt) MakeKennungExpanders() (out kennung.Abbr) {
	out.Etikett.Expand = u.GetStore().GetAbbrStore().Etiketten().ExpandStringString
	out.Typ.Expand = u.GetStore().GetAbbrStore().Typen().ExpandStringString
	out.Kasten.Expand = u.GetStore().GetAbbrStore().Kisten().ExpandStringString
	out.Hinweis.Expand = u.GetStore().GetAbbrStore().Hinweis().ExpandStringString
	out.Sha.Expand = u.GetStore().GetAbbrStore().Shas().ExpandStringString

	out.Etikett.Abbreviate = u.GetStore().GetAbbrStore().Etiketten().Abbreviate
	out.Typ.Abbreviate = u.GetStore().GetAbbrStore().Typen().Abbreviate
	out.Kasten.Abbreviate = u.GetStore().GetAbbrStore().Kisten().Abbreviate
	out.Hinweis.Abbreviate = u.GetStore().GetAbbrStore().Hinweis().Abbreviate
	out.Sha.Abbreviate = u.GetStore().GetAbbrStore().Shas().Abbreviate

	return
}

func (u *Umwelt) MakeMetaIdSetWithExcludedHidden(
	dg kennung.Gattung,
) *query.Builder {
	if dg.IsEmpty() {
		dg = kennung.MakeGattung(gattung.Zettel)
	}

	return query.MakeBuilder(u.Standort(), u.virtualStores["%chrome"]).
		WithDefaultGattungen(dg).
		WithCwd(u.GetStore().GetCwdFiles()).
		WithFileExtensionGetter(u.Konfig().FileExtensions).
		WithHidden(u.GetMatcherArchiviert()).
		WithExpanders(u.MakeKennungExpanders())
}

func (u *Umwelt) MakeMetaIdSetWithoutExcludedHidden(
	dg kennung.Gattung,
) *query.Builder {
	if dg.IsEmpty() {
		dg = kennung.MakeGattung(gattung.Zettel)
	}

	return query.MakeBuilder(u.Standort(), u.virtualStores["%chrome"]).
		WithDefaultGattungen(dg).
		WithCwd(u.GetStore().GetCwdFiles()).
		WithFileExtensionGetter(u.Konfig().FileExtensions).
		WithExpanders(u.MakeKennungExpanders())
}

func (u *Umwelt) ApplyToOrganizeOptions(oo *organize_text.Options) {
	oo.Konfig = u.Konfig()
	oo.Expanders = u.MakeKennungExpanders()
}
