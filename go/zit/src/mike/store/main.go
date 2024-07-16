package store

import (
	"flag"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/lua"
	"code.linenisgreat.com/zit/go/zit/src/echo/fs_home"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/mutable_config"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_id_index"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/blob_store"
	"code.linenisgreat.com/zit/go/zit/src/india/store_fs"
	"code.linenisgreat.com/zit/go/zit/src/juliett/config"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/kilo/external_store"
	"code.linenisgreat.com/zit/go/zit/src/kilo/stream_index"
	"code.linenisgreat.com/zit/go/zit/src/lima/inventory_list"
)

type Store struct {
	sunrise ids.Tai
	config  *config.Compiled
	fs_home fs_home.Home

	cwdFiles           *store_fs.Store
	externalStores     map[string]*external_store.Store
	blob_store         *blob_store.VersionedStores
	inventoryListStore inventory_list.Store
	Abbr               AbbrStore

	inventoryList          *inventory_list.InventoryList
	options                object_inventory_format.Options
	persistentObjectFormat object_inventory_format.Format
	configBlobFormat       blob_store.Format[mutable_config.Blob, *mutable_config.Blob]
	luaVMPoolBuilder       *lua.VMPoolBuilder
	tagLock                sync.Mutex

	streamIndex   *stream_index.Index
	objectIdIndex object_id_index.Index
	typeIndex     object_id_index.ObjectIdIndex[ids.Type, *ids.Type]

	protoZettel  sku.Proto
	queryBuilder *query.Builder

	checkedOutLogPrinter interfaces.FuncIter[sku.CheckedOutLike]
	Logger
}

type Logger struct {
	New, Updated, Unchanged interfaces.FuncIter[*sku.Transacted]
}

func (c *Store) Initialize(
	flags *flag.FlagSet,
	k *config.Compiled,
	st fs_home.Home,
	pmf object_inventory_format.Format,
	t ids.Tai,
	luaVMPoolBuilder *lua.VMPoolBuilder,
	qb *query.Builder,
	options object_inventory_format.Options,
) (err error) {
	c.config = k
	c.fs_home = st
	c.blob_store = blob_store.Make(st)
	c.persistentObjectFormat = pmf
	c.options = options
	c.sunrise = t
	c.luaVMPoolBuilder = luaVMPoolBuilder
	c.queryBuilder = qb

	c.typeIndex = object_id_index.MakeIndex2[ids.Type](
		c.fs_home,
		st.DirVerzeichnisse("TypenIndexV0"),
	)

	c.inventoryList = inventory_list.MakeInventoryList()

	if c.Abbr, err = newIndexAbbr(
		c.fs_home,
		st.DirVerzeichnisse("Abbr"),
	); err != nil {
		err = errors.Wrapf(err, "failed to init abbr index")
		return
	}

	if c.inventoryListStore, err = inventory_list.MakeStore(
		c.GetStandort(),
		c.GetStandort().GetLockSmith(),
		c.config.GetStoreVersion(),
		c.fs_home.ObjekteReaderWriterFactory(genres.InventoryList),
		c.fs_home,
		pmf,
		c,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if c.objectIdIndex, err = object_id_index.MakeIndex(
		c.GetKonfig(),
		c.GetStandort(),
		c.GetStandort(),
	); err != nil {
		err = errors.Wrapf(err, "failed to init zettel index")
		return
	}

	if c.streamIndex, err = stream_index.MakeIndex(
		c.GetStandort(),
		c.GetKonfig(),
		c.GetStandort().DirVerzeichnisseObjekten(),
		pmf,
		c.options,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	c.protoZettel = sku.MakeProto(
		k.GetMutableConfig().Defaults.Typ,
		k.DefaultTags,
	)

	c.configBlobFormat = blob_store.MakeBlobFormat(
		blob_store.MakeTextParserIgnoreTomlErrors[mutable_config.Blob](
			c.GetStandort(),
		),
		blob_store.ParsedBlobTomlFormatter[mutable_config.Blob, *mutable_config.Blob]{},
		c.GetStandort(),
	)

	return
}

func (s *Store) SetExternalStores(
	stores map[string]*external_store.Store,
) (err error) {
	s.externalStores = stores

	for k, es := range s.externalStores {
		es.StoreFuncs = external_store.StoreFuncs{
			FuncRealize:     s.tryRealize,
			FuncCommit:      s.tryRealizeAndOrStore,
			FuncReadSha:     s.ReadOneEnnui,
			FuncReadOneInto: s.ReadOneInto,
			FuncQuery:       s.Query,
		}

		es.Home = s.GetStandort()
		es.DirCache = s.GetStandort().DirVerzeichnisseKasten(k)

		if esfs, ok := es.StoreLike.(*store_fs.Store); ok {
			s.cwdFiles = esfs

			// TODO remove once store_fs.Store is fully ExternalStoreLike
			if err = es.Initialize(); err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	return
}

// TODO remove
func (s *Store) SetCheckedOutLogWriter(
	zelw interfaces.FuncIter[sku.CheckedOutLike],
) {
	s.checkedOutLogPrinter = zelw
}

func (s *Store) ResetIndexes() (err error) {
	if err = s.typeIndex.Reset(); err != nil {
		err = errors.Wrapf(err, "failed to reset etiketten index")
		return
	}

	if err = s.objectIdIndex.Reset(); err != nil {
		err = errors.Wrapf(err, "failed to reset index object id index")
		return
	}

	return
}

func (s *Store) SetLogWriter(lw Logger) {
	s.Logger = lw
}
