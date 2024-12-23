package repo_local

import (
	"io"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/bravo/xdg"
	"code.linenisgreat.com/zit/go/zit/src/delta/age"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/box_format"
	"code.linenisgreat.com/zit/go/zit/src/india/dormant_index"
	"code.linenisgreat.com/zit/go/zit/src/juliett/blob_store"
	"code.linenisgreat.com/zit/go/zit/src/juliett/config"
	"code.linenisgreat.com/zit/go/zit/src/kilo/external_store"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
	"code.linenisgreat.com/zit/go/zit/src/lima/store_browser"
	"code.linenisgreat.com/zit/go/zit/src/lima/store_fs"
	"code.linenisgreat.com/zit/go/zit/src/mike/store"
)

type Repo struct {
	*env.Env

	sunrise ids.Tai

	dirLayout    dir_layout.DirLayout
	fileEncoder  store_fs.FileEncoder
	config       config.Compiled
	dormantIndex dormant_index.Index

	storesInitialized bool
	blobStore         *blob_store.VersionedStores
	store             store.Store
	age               *age.Age
	externalStores    map[ids.RepoId]*external_store.Store

	DormantCounter query.DormantCounter

	luaSkuFormat *box_format.BoxTransacted
}

func MakeFromConfigAndXDGDotenvPath(
	context errors.Context,
	config *config.Compiled,
	xdgDotenvPath string,
) (local *Repo, err error) {
	dotenv := xdg.Dotenv{
		XDG: &xdg.XDG{},
	}

	var f *os.File

	if f, err = os.Open(xdgDotenvPath); err != nil {
		err = errors.Wrap(err)
		return
	}

	if _, err = dotenv.ReadFrom(f); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = f.Close(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var primitiveFSHome dir_layout.Primitive

	if primitiveFSHome, err = dir_layout.MakePrimitiveWithXDG(
		config.Debug,
		*dotenv.XDG,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	env := env.Make(
		context,
		nil,
		config.Cli(),
		primitiveFSHome,
	)

	if local, err = Make(
		env,
		OptionsEmpty,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func Make(
	env *env.Env,
	options Options,
) (u *Repo, err error) {
	u = &Repo{
		Env:            env,
		DormantCounter: query.MakeDormantCounter(),
	}

	u.config.Reset()

	err = u.Initialize(options)

	return
}

func (u *Repo) GetRepo() repo.Repo {
	return u
}

// TODO investigate removing unnecessary resets like from organize
func (u *Repo) Reset() (err error) {
	return u.Initialize(OptionsEmpty)
}

func (u *Repo) Initialize(options Options) (err error) {
	if err = u.Flush(); err != nil {
		err = errors.Wrap(err)
		return
	}

	u.sunrise = ids.NowTai()

	ui.TodoP4("find a better place for this")
	{
		if u.GetCLIConfig().Verbose && !u.GetCLIConfig().Quiet {
			ui.SetVerbose(true)
		} else {
			ui.SetOutput(io.Discard)
		}

		if u.GetCLIConfig().Todo {
			ui.SetTodoOn()
		}

		standortOptions := dir_layout.Options{
			BasePath: u.GetCLIConfig().BasePath,
		}

		if u.dirLayout, err = dir_layout.Make(
			standortOptions,
			u.GetDirLayoutPrimitive(),
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
		u.GetCLIConfig(),
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

	// for _, rb := range u.GetConfig().Recipients {
	// 	if err = u.age.AddBech32PivYubikeyEC256(rb); err != nil {
	// 		errors.Wrap(err)
	// 		return
	// 	}
	// }
	ofo := object_inventory_format.Options{Tai: true}

	if err = u.store.Initialize(
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

	lw := store.UIDelegate{
		TransactedNew:     ptl,
		TransactedUpdated: ptl,
		TransactedUnchanged: func(sk *sku.Transacted) (err error) {
			if !u.config.PrintOptions.PrintUnchanged {
				return
			}

			return ptl(sk)
		},
		CheckedOutCheckedOut: u.PrinterCheckedOut(
			box_format.CheckedOutHeaderState{},
		),
	}

	u.store.SetUIDelegate(lw)

	u.storesInitialized = true

	u.luaSkuFormat = u.SkuFormatBoxTransactedNoColor()

	return
}

func (u *Repo) Flush() (err error) {
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

func (u *Repo) PrintMatchedArchiviertIfNecessary() {
	if !u.GetConfig().PrintOptions.PrintMatchedDormant {
		return
	}

	c := u.GetMatcherDormant().Count()
	ca := u.GetMatcherDormant().CountArchiviert()

	if c != 0 || ca == 0 {
		return
	}

	ui.Err().Printf("%d archived objects matched", c)
}

func (u *Repo) MakeObjectIdIndex() ids.Index {
	return ids.Index{}
}

func (u *Repo) GetMatcherDormant() query.DormantCounter {
	return u.DormantCounter
}

func (u *Repo) GetExternalStoreForQuery(
	repoId ids.RepoId,
) (sku.ExternalStoreForQuery, bool) {
	e, ok := u.externalStores[repoId]
	return e, ok
}
