package umwelt

import (
	"flag"
	"io"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/age"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/delta/lua"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/echo/standort"
	"code.linenisgreat.com/zit/go/zit/src/echo/thyme"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/erworben"
	"code.linenisgreat.com/zit/go/zit/src/golf/objekte_format"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt"
	"code.linenisgreat.com/zit/go/zit/src/juliett/chrome"
	"code.linenisgreat.com/zit/go/zit/src/juliett/konfig"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/kilo/organize_text"
	"code.linenisgreat.com/zit/go/zit/src/mike/store"
)

type Umwelt struct {
	sonnenaufgang thyme.Time

	in  *os.File
	out *os.File
	err *os.File

	flags *flag.FlagSet

	inIsTty  bool
	outIsTty bool
	errIsTty bool

	standort    standort.Standort
	erworbenCli erworben.Cli
	konfig      konfig.Compiled
	schlummernd query.Schlummernd

	storesInitialized bool
	store             store.Store
	virtualStores     map[string]*query.VirtualStoreInitable
	age               *age.Age

	matcherArchiviert query.Archiviert

	luaSkuFormat *sku_fmt.Organize
}

func Make(
	flags *flag.FlagSet,
	kCli erworben.Cli,
	options Options,
) (u *Umwelt, err error) {
	u = &Umwelt{
		in:                os.Stdin,
		out:               os.Stdout,
		err:               os.Stderr,
		flags:             flags,
		erworbenCli:       kCli,
		matcherArchiviert: query.MakeArchiviert(),
		virtualStores:     make(map[string]*query.VirtualStoreInitable),
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

// TODO investigate removing unnecessary resets like from organize
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
		if u.erworbenCli.Verbose && !u.erworbenCli.Quiet {
			ui.SetVerbose()
		} else {
			ui.SetOutput(io.Discard)
		}

		if u.erworbenCli.Todo {
			errors.SetTodoOn()
		}

		standortOptions := standort.Options{
			BasePath: u.erworbenCli.BasePath,
			Debug:    u.erworbenCli.Debug,
			DryRun:   u.erworbenCli.DryRun,
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

	if err = u.schlummernd.Load(
		u.standort,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.konfig.Initialize(
		u.standort,
		u.erworbenCli,
		&u.schlummernd,
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

	ui.Log().Printf("store version: %s", u.GetKonfig().GetStoreVersion())

	if u.konfig.ChrestEnabled {
		u.virtualStores["%chrome"] = &query.VirtualStoreInitable{
			VirtualStore: chrome.MakeChrome(u.GetKonfig(), u.Standort()),
		}
	}

	if err = u.store.Initialize(
		u.flags,
		u.GetKonfig(),
		u.standort,
		objekte_format.FormatForVersion(u.GetKonfig().GetStoreVersion()),
		u.sonnenaufgang,
		// u.virtualStores,
		(&lua.VMPoolBuilder{}).WithSearcher(u.LuaSearcher),
		// u.PrinterFDDeleted(),
		u.makeQueryBuilder().
			WithDefaultGattungen(kennung.MakeGattung(gattung.TrueGattung()...)),
	); err != nil {
		err = errors.Wrapf(err, "failed to initialize store util")
		return
	}

	ui.Log().Print("done initing checkout store")

	ptl := u.PrinterTransactedLike()

	lw := store.Logger{
		New:     ptl,
		Updated: ptl,
		Unchanged: func(sk *sku.Transacted) (err error) {
			if !u.konfig.PrintOptions.PrintUnchanged {
				return
			}

			return ptl(sk)
		},
	}

	u.store.SetCheckedOutLogWriter(u.PrinterCheckedOutLike())
	u.store.SetLogWriter(lw)

	u.storesInitialized = true

	u.luaSkuFormat = u.SkuFmtOrganize()

	return
}

func (u *Umwelt) Flush() error {
	return u.age.Close()
}

func (u *Umwelt) PrintMatchedArchiviertIfNecessary() {
	if !u.GetKonfig().PrintOptions.PrintMatchedArchiviert {
		return
	}

	c := u.GetMatcherArchiviert().Count()
	ca := u.GetMatcherArchiviert().CountArchiviert()

	if c != 0 || ca == 0 {
		return
	}

	ui.Err().Printf("%d archived objekten matched", c)
}

func (u *Umwelt) MakeKennungIndex() kennung.Index {
	return kennung.Index{}
}

func (u *Umwelt) GetMatcherArchiviert() query.Archiviert {
	return u.matcherArchiviert
}

func (u *Umwelt) GetChrestStore() *query.VirtualStoreInitable {
	return u.virtualStores["%chrome"]
}

func (u *Umwelt) makeQueryBuilder() *query.Builder {
	return query.MakeBuilder(
		u.Standort(),
		u.GetStore().GetAkten(),
		u.GetStore().GetVerzeichnisse(),
		u.GetChrestStore(),
		(&lua.VMPoolBuilder{}).WithSearcher(u.LuaSearcher),
	)
}

func (u *Umwelt) MakeQueryBuilderExcludingHidden(
	dg kennung.Gattung,
) *query.Builder {
	if dg.IsEmpty() {
		dg = kennung.MakeGattung(gattung.Zettel)
	}

	return u.makeQueryBuilder().
		WithDefaultGattungen(dg).
		WithVirtualEtiketten(u.konfig.Filters).
		WithCwd(u.GetStore().GetCwdFiles()).
		WithFileExtensionGetter(u.GetKonfig().FileExtensions).
		WithHidden(u.GetMatcherArchiviert()).
		WithExpanders(u.GetStore().GetAbbrStore().GetAbbr())
}

func (u *Umwelt) MakeQueryBuilder(
	dg kennung.Gattung,
) *query.Builder {
	if dg.IsEmpty() {
		dg = kennung.MakeGattung(gattung.Zettel)
	}

	return u.makeQueryBuilder().
		WithDefaultGattungen(dg).
		WithVirtualEtiketten(u.konfig.Filters).
		WithCwd(u.GetStore().GetCwdFiles()).
		WithFileExtensionGetter(u.GetKonfig().FileExtensions).
		WithExpanders(u.GetStore().GetAbbrStore().GetAbbr())
}

func (u *Umwelt) ApplyToOrganizeOptions(oo *organize_text.Options) {
	oo.Konfig = u.GetKonfig()
	oo.Abbr = u.GetStore().GetAbbrStore().GetAbbr()
}

func (u *Umwelt) GetVirtualStores() map[string]*query.VirtualStoreInitable {
	return u.virtualStores
}
