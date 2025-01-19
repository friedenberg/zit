package local_working_copy

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/repo_type"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/config_immutable"
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
	"code.linenisgreat.com/zit/go/zit/src/kilo/store_abbr"
	"code.linenisgreat.com/zit/go/zit/src/lima/env_lua"
	"code.linenisgreat.com/zit/go/zit/src/lima/store_browser"
	"code.linenisgreat.com/zit/go/zit/src/lima/store_fs"
	"code.linenisgreat.com/zit/go/zit/src/lima/typed_blob_store"
	"code.linenisgreat.com/zit/go/zit/src/mike/env_box"
	"code.linenisgreat.com/zit/go/zit/src/mike/store"
	"code.linenisgreat.com/zit/go/zit/src/mike/store_config"
)

type (
	envLocal = env_local.Env
	envBox   = env_box.Env
)

type Repo struct {
	envLocal
	envBox

	sunrise ids.Tai

	envRepo     env_repo.Env
	fileEncoder store_fs.FileEncoder
	config      store_config.StoreMutable

	storeFS      *store_fs.Store
	storeAbbr    sku.AbbrStore
	dormantIndex dormant_index.Index

	storesInitialized bool
	typedBlobStore    *typed_blob_store.Store
	store             store.Store
	externalStores    map[ids.RepoId]*external_store.Store

	DormantCounter query.DormantCounter

	envLua env_lua.Env
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
		config:         store_config.Make(),
		envLocal:       repoLayout,
		envRepo:        repoLayout,
		DormantCounter: query.MakeDormantCounter(),
	}

	repo.config.Reset()

	if err := repo.initialize(options); err != nil {
		repo.CancelWithError(err)
	}

	repo.After(repo.Flush)

	return
}

func (u *Repo) GetImmutableConfig() config_immutable.Config {
	return u.GetRepoLayout().GetConfig()
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

	repo.fileEncoder = store_fs.MakeFileEncoder(repo.envRepo, repo.config)

	if err = repo.dormantIndex.Load(
		repo.envRepo,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	objectFormat := object_inventory_format.FormatForVersion(repo.envRepo.GetStoreVersion())

	boxFormatArchive := box_format.MakeBoxTransactedArchive(
		repo.GetEnv(),
		repo.GetConfig().GetCLIConfig().PrintOptions.WithPrintTai(true),
	)

	if err = repo.config.Initialize(
		repo.envRepo,
		repo.GetCLIConfig(),
	); err != nil {
		if options.GetAllowConfigReadError() {
			err = nil
		} else {
			err = errors.Wrapf(
				err,
				"CompressionType: %q",
				repo.envRepo.GetConfig().GetBlobStoreConfigImmutable().GetBlobCompression(),
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

	objectInventoryFormatOptions := object_inventory_format.Options{Tai: true}

	config := repo.GetConfig()

	if repo.storeFS, err = store_fs.Make(
		config,
		repo.PrinterFDDeleted(),
		config.GetFileExtensions(),
		repo.GetRepoLayout(),
		objectInventoryFormatOptions,
		repo.fileEncoder,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if repo.storeAbbr, err = store_abbr.NewIndexAbbr(
		config.GetCLIConfig().PrintOptions,
		repo.envRepo,
	); err != nil {
		err = errors.Wrapf(err, "failed to init abbr index")
		return
	}

	repo.envBox = env_box.Make(
		repo.envRepo,
		repo.storeFS,
		repo.storeAbbr,
	)

	repo.envLua = env_lua.Make(
		repo.envRepo,
		repo.GetStore(),
		repo.SkuFormatBoxTransactedNoColor(),
	)

	// for _, rb := range u.GetConfig().Recipients {
	// 	if err = u.age.AddBech32PivYubikeyEC256(rb); err != nil {
	// 		errors.Wrap(err)
	// 		return
	// 	}
	// }

	repo.typedBlobStore = typed_blob_store.Make(
		repo.envRepo,
		repo.envLua,
		objectFormat,
		boxFormatArchive,
	)

	if err = repo.store.Initialize(
		repo.config,
		repo.envRepo,
		objectFormat,
		repo.sunrise,
		repo.envLua,
		repo.makeQueryBuilder().
			WithDefaultGenres(ids.MakeGenre(genres.TrueGenre()...)),
		objectInventoryFormatOptions,
		boxFormatArchive,
		repo.typedBlobStore,
		&repo.dormantIndex,
		repo.storeAbbr,
	); err != nil {
		err = errors.Wrapf(err, "failed to initialize store util")
		return
	}

	ui.Log().Printf(
		"store version: %s",
		repo.GetConfig().GetImmutableConfig().GetStoreVersion(),
	)

	repo.externalStores = map[ids.RepoId]*external_store.Store{
		{}: {
			StoreLike: repo.storeFS,
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
