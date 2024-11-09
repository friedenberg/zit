package env

import (
	"flag"
	"io"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/age"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/mutable_config_blobs"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/box_format"
	"code.linenisgreat.com/zit/go/zit/src/india/dormant_index"
	"code.linenisgreat.com/zit/go/zit/src/juliett/blob_store"
	"code.linenisgreat.com/zit/go/zit/src/juliett/config"
	"code.linenisgreat.com/zit/go/zit/src/kilo/external_store"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/store_browser"
	"code.linenisgreat.com/zit/go/zit/src/lima/store_fs"
	"code.linenisgreat.com/zit/go/zit/src/mike/store"
)

type Env struct {
	sunrise ids.Tai

	in  *os.File
	out *os.File
	err *os.File

	flags *flag.FlagSet

	inIsTty  bool
	outIsTty bool
	errIsTty bool

	primitiveFSHome dir_layout.Primitive
	dirLayout       dir_layout.DirLayout
	cliConfig       mutable_config_blobs.Cli
	fileEncoder     store_fs.FileEncoder
	config          config.Compiled
	dormantIndex    dormant_index.Index

	storesInitialized bool
	blobStore         *blob_store.VersionedStores
	store             store.Store
	age               *age.Age
	externalStores    map[ids.RepoId]*external_store.Store

	DormantCounter query.DormantCounter

	luaSkuFormat *box_format.BoxTransacted
}

func Make(
	flags *flag.FlagSet,
	kCli mutable_config_blobs.Cli,
	options Options,
	primitiveFSHome dir_layout.Primitive,
) (u *Env, err error) {
	u = &Env{
		in:              os.Stdin,
		out:             os.Stdout,
		err:             os.Stderr,
		flags:           flags,
		cliConfig:       kCli,
		DormantCounter:  query.MakeDormantCounter(),
		primitiveFSHome: primitiveFSHome,
	}

	u.config.Reset()

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
func (u *Env) Reset() (err error) {
	return u.Initialize(OptionsEmpty)
}

func (u *Env) Initialize(options Options) (err error) {
	if err = u.Flush(); err != nil {
		err = errors.Wrap(err)
		return
	}

	u.sunrise = ids.NowTai()

	ui.TodoP4("find a better place for this")
	{
		if u.cliConfig.Verbose && !u.cliConfig.Quiet {
			ui.SetVerbose()
		} else {
			ui.SetOutput(io.Discard)
		}

		if u.cliConfig.Todo {
			ui.SetTodoOn()
		}

		standortOptions := dir_layout.Options{
			BasePath: u.cliConfig.BasePath,
		}

		if u.dirLayout, err = dir_layout.Make(
			standortOptions,
			u.primitiveFSHome,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	u.fileEncoder = store_fs.MakeFileEncoder(u.dirLayout, &u.config)

	if err = u.dormantIndex.Load(
		u.dirLayout,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	objectFormat := object_inventory_format.FormatForVersion(u.GetConfig().GetStoreVersion())
	boxFormatArchive := u.MakeBoxArchive(true)

	u.blobStore = blob_store.Make(
		u.dirLayout,
		u.MakeLuaVMPoolBuilder(),
		objectFormat,
		boxFormatArchive,
	)

	if err = u.config.Initialize(
		u.dirLayout,
		u.cliConfig,
		&u.dormantIndex,
		u.blobStore,
	); err != nil {
		if options.GetAllowConfigReadError() {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	u.config.ApplyPrintOptionsConfig(u.config.PrintOptions)

	// for _, rb := range u.GetConfig().Recipients {
	// 	if err = u.age.AddBech32PivYubikeyEC256(rb); err != nil {
	// 		errors.Wrap(err)
	// 		return
	// 	}
	// }
	ofo := object_inventory_format.Options{Tai: true}

	if err = u.store.Initialize(
		u.flags,
		u.GetConfig(),
		u.dirLayout,
		objectFormat,
		u.sunrise,
		u.MakeLuaVMPoolBuilder(),
		u.makeQueryBuilder().
			WithDefaultGenres(ids.MakeGenre(genres.TrueGenre()...)),
		ofo,
		boxFormatArchive,
		u.blobStore,
	); err != nil {
		err = errors.Wrapf(err, "failed to initialize store util")
		return
	}

	ui.Log().Printf("store version: %s", u.GetConfig().GetStoreVersion())

	var sfs *store_fs.Store

	k := u.GetConfig()

	if sfs, err = store_fs.Make(
		k,
		u.PrinterFDDeleted(),
		k.GetFileExtensions(),
		u.GetDirectoryLayout(),
		ofo,
		u.fileEncoder,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	u.externalStores = map[ids.RepoId]*external_store.Store{
		{}: {
			StoreLike: sfs,
		},
		*(ids.MustRepoId("browser")): {
			StoreLike: store_browser.Make(
				k,
				u.GetDirectoryLayout(),
				u.PrinterTransactedDeleted(),
			),
		},
	}

	if err = u.store.SetExternalStores(
		u.externalStores,
	); err != nil {
		err = errors.Wrapf(err, "failed to set external stores")
		return
	}

	ui.Log().Print("done initing checkout store")

	ptl := u.PrinterTransacted()

	lw := store.Logger{
		New:     ptl,
		Updated: ptl,
		Unchanged: func(sk *sku.Transacted) (err error) {
			if !u.config.PrintOptions.PrintUnchanged {
				return
			}

			return ptl(sk)
		},
	}

	u.store.SetCheckedOutLogWriter(u.PrinterCheckedOut())
	u.store.SetLogWriter(lw)

	u.storesInitialized = true

	u.luaSkuFormat = u.SkuFormatBoxNoColor()

	return
}

func (u *Env) Flush() (err error) {
	wg := quiter.MakeErrorWaitGroupParallel()

	wg.Do(u.age.Close)

	for k, vs := range u.externalStores {
		ui.Log().Printf("will flush virtual store: %s", k)
		wg.Do(vs.Flush)
	}

	if err = wg.GetError(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (u *Env) PrintMatchedArchiviertIfNecessary() {
	if !u.GetConfig().PrintOptions.PrintMatchedDormant {
		return
	}

	c := u.GetMatcherArchiviert().Count()
	ca := u.GetMatcherArchiviert().CountArchiviert()

	if c != 0 || ca == 0 {
		return
	}

	ui.Err().Printf("%d archived objects matched", c)
}

func (u *Env) MakeObjectIdIndex() ids.Index {
	return ids.Index{}
}

func (u *Env) GetMatcherArchiviert() query.DormantCounter {
	return u.DormantCounter
}

func (u *Env) GetExternalStoreForQuery(
	repoId ids.RepoId,
) (sku.ExternalStoreForQuery, bool) {
	e, ok := u.externalStores[repoId]
	return e, ok
}
