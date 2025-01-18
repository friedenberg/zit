package local_working_copy

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/alfa/repo_type"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_local"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_repo"
	"code.linenisgreat.com/zit/go/zit/src/hotel/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/box_format"
	"code.linenisgreat.com/zit/go/zit/src/kilo/dormant_index"
	"code.linenisgreat.com/zit/go/zit/src/kilo/external_store"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/blob_store"
	"code.linenisgreat.com/zit/go/zit/src/lima/store_browser"
	"code.linenisgreat.com/zit/go/zit/src/lima/store_fs"
	"code.linenisgreat.com/zit/go/zit/src/mike/env_config"
	"code.linenisgreat.com/zit/go/zit/src/mike/store"
)

type Repo struct {
	env_local.Env

	sunrise ids.Tai

	layout       env_repo.Env
	fileEncoder  store_fs.FileEncoder
	config       env_config.EnvMutable
	dormantIndex dormant_index.Index

	storesInitialized bool
	blobStore         *blob_store.VersionedStores
	store             store.Store
	externalStores    map[ids.RepoId]*external_store.Store

	DormantCounter query.DormantCounter

	luaSkuFormat *box_format.BoxTransacted
}

func Make(
	env env_local.Env,
	options Options,
) *Repo {
	layoutOptions := env_repo.Options{
		BasePath: env.GetCLIConfig().BasePath,
	}

	var repoLayout env_repo.Env

	{
		var err error

		if repoLayout, err = env_repo.Make(
			env,
			layoutOptions,
		); err != nil {
			env.CancelWithError(err)
		}
	}

	return MakeWithLayout(options, repoLayout)
}

func MakeWithLayout(
	options Options,
	repoLayout env_repo.Env,
) (repo *Repo) {
	repo = &Repo{
		config:         env_config.Make(),
		Env:            repoLayout,
		layout:         repoLayout,
		DormantCounter: query.MakeDormantCounter(),
	}

	repo.config.Reset()

	if err := repo.initialize(options); err != nil {
		repo.CancelWithError(err)
	}

	repo.After(repo.Flush)

	return
}

func (u *Repo) GetRepoType() repo_type.Type {
	return u.GetRepoLayout().GetConfig().GetRepoType()
}

func (u *Repo) GetStoreVersion() interfaces.StoreVersion {
	return u.GetRepoLayout().GetConfig().GetStoreVersion()
}

// TODO investigate removing unnecessary resets like from organize
func (u *Repo) Reset() (err error) {
	return u.initialize(OptionsEmpty)
}

func (repo *Repo) initialize(
	options Options,
) (err error) {
	if err = repo.Flush(); err != nil {
		err = errors.Wrap(err)
		return
	}

	// ui.Debug().Print(repo.layout.GetConfig().GetBlobStoreImmutableConfig().GetCompressionType())
	repo.sunrise = ids.NowTai()

	repo.fileEncoder = store_fs.MakeFileEncoder(repo.layout, repo.config)

	if err = repo.dormantIndex.Load(
		repo.layout,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	objectFormat := object_inventory_format.FormatForVersion(repo.layout.GetStoreVersion())
	boxFormatArchive := box_format.MakeBoxTransactedArchive(
		repo.GetEnv(),
		repo.GetConfig().GetCLIConfig().PrintOptions.WithPrintTai(true),
	)

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
			err = errors.Wrapf(
				err,
				"CompressionType: %q",
				repo.layout.GetConfig().GetBlobStoreImmutableConfig().GetCompressionType(),
			)
			return
		}
	}

	if repo.GetConfig().GetRepoType() != repo_type.TypeWorkingCopy {
		err = repo_type.ErrUnsupportedRepoType{
			Expected: repo_type.TypeWorkingCopy,
			Actual:   repo.GetConfig().GetImmutableConfig().GetRepoType(),
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
		repo.config,
		repo.layout,
		objectFormat,
		repo.sunrise,
		repo.MakeLuaVMPoolBuilder(),
		repo.makeQueryBuilder().
			WithDefaultGenres(ids.MakeGenre(genres.TrueGenre()...)),
		ofo,
		boxFormatArchive,
		repo.blobStore,
		&repo.dormantIndex,
	); err != nil {
		err = errors.Wrapf(err, "failed to initialize store util")
		return
	}

	ui.Log().Printf(
		"store version: %s",
		repo.GetConfig().GetImmutableConfig().GetStoreVersion(),
	)

	var sfs *store_fs.Store

	config := repo.GetConfig()

	if sfs, err = store_fs.Make(
		config,
		repo.PrinterFDDeleted(),
		config.GetFileExtensions(),
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
				config,
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
			if !repo.config.GetCLIConfig().PrintOptions.PrintUnchanged {
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
	if !u.GetConfig().GetCLIConfig().PrintOptions.PrintMatchedDormant {
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
