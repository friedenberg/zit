package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/hotel/objekte_store"
)

type GattungStore interface{}

type reindexer interface {
	// updateExternal(objekte.External) error
	reindexOne(sku.DataIdentity) (schnittstellen.Stored, error)
}

type commonStore[
	OBJEKTE schnittstellen.Objekte[OBJEKTE],
	OBJEKTEPtr schnittstellen.ObjektePtr[OBJEKTE],
	KENNUNG schnittstellen.Id[KENNUNG],
	KENNUNGPtr schnittstellen.IdPtr[KENNUNG],
	VERZEICHNISSE any,
	VERZEICHNISSEPtr schnittstellen.VerzeichnissePtr[VERZEICHNISSE, OBJEKTE],
] interface {
	reindexer
	GattungStore

	objekte_store.TransactedLogger[*objekte.Transacted[
		OBJEKTE,
		OBJEKTEPtr,
		KENNUNG,
		KENNUNGPtr,
		VERZEICHNISSE,
		VERZEICHNISSEPtr,
	]]

	objekte_store.Querier[
		KENNUNGPtr,
		*objekte.Transacted[
			OBJEKTE,
			OBJEKTEPtr,
			KENNUNG,
			KENNUNGPtr,
			VERZEICHNISSE,
			VERZEICHNISSEPtr,
		],
	]

	objekte_store.AkteTextSaver[
		OBJEKTE,
		OBJEKTEPtr,
	]

	objekte_store.CreateOrUpdater[
		OBJEKTEPtr,
		KENNUNGPtr,
		*objekte.Transacted[
			OBJEKTE,
			OBJEKTEPtr,
			KENNUNG,
			KENNUNGPtr,
			VERZEICHNISSE,
			VERZEICHNISSEPtr,
		],
		*objekte.CheckedOut[
			OBJEKTE,
			OBJEKTEPtr,
			KENNUNG,
			KENNUNGPtr,
			VERZEICHNISSE,
			VERZEICHNISSEPtr,
		],
	]

	objekte_store.TransactedInflator[
		OBJEKTE,
		OBJEKTEPtr,
		KENNUNG,
		KENNUNGPtr,
		VERZEICHNISSE,
		VERZEICHNISSEPtr,
	]

	objekte_store.Inheritor[*objekte.Transacted[
		OBJEKTE,
		OBJEKTEPtr,
		KENNUNG,
		KENNUNGPtr,
		VERZEICHNISSE,
		VERZEICHNISSEPtr,
	]]
}
