package repo_local

import (
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/repo_type"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/age"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/xdg"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/echo/repo_layout"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/config_mutable_cli"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_inventory_format"
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

	layout       repo_layout.Layout
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
	context *errors.Context,
	config config_mutable_cli.Config,
	xdgDotenvPath string,
	options env.Options,
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

	dirLayout := dir_layout.MakeWithXDG(
		context,
		config.Debug,
		*dotenv.XDG,
	)

	env := env.Make(
		context,
		config,
		dirLayout,
		options,
	)

	local = Make(
		env,
		OptionsEmpty,
	)

	return
}

func Make(
	env *env.Env,
	options Options,
) (repo *Repo) {
	repo = &Repo{
		Env:            env,
		DormantCounter: query.MakeDormantCounter(),
	}

	repo.config.Reset()

	if err := repo.initialize(options); err != nil {
		env.CancelWithError(err)
	}

	repo.After(repo.Flush)

	return
}

func (u *Repo) GetRepo() repo.Repo {
	return u
}

// TODO investigate removing unnecessary resets like from organize
func (u *Repo) Reset() (err error) {
	return u.initialize(OptionsEmpty)
}

func (repo *Repo) initialize(options Options) (err error) {
	if err = repo.Flush(); err != nil {
		err = errors.Wrap(err)
		return
	}

	repo.sunrise = ids.NowTai()

	ui.TodoP4("find a better place for this")
	{
		layoutOptions := repo_layout.Options{
			BasePath: repo.GetCLIConfig().BasePath,
		}

		if repo.layout, err = repo_layout.Make(
			repo.Env,
			layoutOptions,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	repo.fileEncoder = store_fs.MakeFileEncoder(repo.layout, &repo.config)

	if err = repo.dormantIndex.Load(
		repo.layout,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	objectFormat := object_inventory_format.FormatForVersion(repo.layout.GetStoreVersion())
	boxFormatArchive := repo.MakeBoxArchive(true)

	repo.blobStore = blob_store.Make(
		repo.layout,
		repo.MakeLuaVMPoolBuilder(),
		objectFormat,
		boxFormatArchive,
	)

	if err = repo.config.Initialize(
		repo.layout,
		repo.GetCLIConfig(),
		&repo.dormantIndex,
		repo.blobStore,
	); err != nil {
		if options.GetAllowConfigReadError() {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	if repo.GetConfig().GetRepoType() != repo_type.TypeReadWrite {
		err = repo_type.ErrUnsupportedRepoType{
			Expected: repo_type.TypeReadWrite,
			Actual:   repo.GetConfig().GetRepoType(),
		}

		return
	}

	// for _, rb := range u.GetConfig().Recipients {
	// 	if err = u.age.AddBech32PivYubikeyEC256(rb); err != nil {
	// 		errors.Wrap(err)
	// 		return
	// 	}
	// }
	ofo := object_inventory_format.Options{Tai: true}

	if err = repo.store.Initialize(
		repo.GetConfig(),
		repo.layout,
		objectFormat,
		repo.sunrise,
		repo.MakeLuaVMPoolBuilder(),
		repo.makeQueryBuilder().
			WithDefaultGenres(ids.MakeGenre(genres.TrueGenre()...)),
		ofo,
		boxFormatArchive,
		repo.blobStore,
	); err != nil {
		err = errors.Wrapf(err, "failed to initialize store util")
		return
	}

	ui.Log().Printf("store version: %s", repo.GetConfig().GetStoreVersion())

	var sfs *store_fs.Store

	k := repo.GetConfig()

	if sfs, err = store_fs.Make(
		k,
		repo.PrinterFDDeleted(),
		k.GetFileExtensions(),
		repo.GetRepoLayout(),
		ofo,
		repo.fileEncoder,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	repo.externalStores = map[ids.RepoId]*external_store.Store{
		{}: {
			StoreLike: sfs,
		},
		*(ids.MustRepoId("browser")): {
			StoreLike: store_browser.Make(
				k,
				repo.GetRepoLayout(),
				repo.PrinterTransactedDeleted(),
			),
		},
	}

	if err = repo.store.SetExternalStores(
		repo.externalStores,
	); err != nil {
		err = errors.Wrapf(err, "failed to set external stores")
		return
	}

	ui.Log().Print("done initing checkout store")

	ptl := repo.PrinterTransacted()

	lw := store.UIDelegate{
		TransactedNew:     ptl,
		TransactedUpdated: ptl,
		TransactedUnchanged: func(sk *sku.Transacted) (err error) {
			if !repo.config.PrintOptions.PrintUnchanged {
				return
			}

			return ptl(sk)
		},
		CheckedOutCheckedOut: repo.PrinterCheckedOut(
			box_format.CheckedOutHeaderState{},
		),
	}

	repo.store.SetUIDelegate(lw)

	repo.storesInitialized = true

	repo.luaSkuFormat = repo.SkuFormatBoxTransactedNoColor()

	return
}

func (u *Repo) Flush() (err error) {
	wg := errors.MakeWaitGroupParallel()

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

func (u *Repo) PrintMatchedDormantIfNecessary() {
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
