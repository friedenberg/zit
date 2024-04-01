package kennung

import (
	"encoding/gob"
	"sync"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/charlie/collections_ptr"
	"code.linenisgreat.com/zit/src/delta/gattung"
)

var (
	registerOnce        sync.Once
	registryLock        *sync.Mutex
	registryGattung     map[gattung.Gattung]Kennung
	registryQueryPrefix map[string]Kennung
)

func once() {
	registryLock = &sync.Mutex{}
	registryGattung = make(map[gattung.Gattung]Kennung)
	registryQueryPrefix = make(map[string]Kennung)
}

func register[T Kennung, TPtr interface {
	schnittstellen.StringSetterPtr[T]
	Kennung
}](id T,
) {
	gob.Register(&id)
	gob.Register(collections_ptr.MakeMutableValueSet[T, TPtr](nil))
	gob.Register(collections_ptr.MakeValueSet[T, TPtr](nil))
	registerOnce.Do(once)

	registryLock.Lock()
	defer registryLock.Unlock()

	ok := false
	var id1 Kennung
	g := gattung.Must(id.GetGattung())

	if id1, ok = registryGattung[g]; ok {
		panic(
			errors.Errorf(
				"gattung %s has two registrations: %s (old), %s (new)",
				g,
				id1,
				id,
			),
		)
	}

	registryGattung[g] = id

	if idQueryPrefix, ok := Kennung(id).(QueryPrefixer); ok {
		p := idQueryPrefix.GetQueryPrefix()

		if id1, ok := registryQueryPrefix[p]; ok {
			panic(
				errors.Errorf(
					"prefix '%s' has two registrations: %s (old), %s (new)",
					p,
					id1,
					id,
				),
			)
		}

		registryQueryPrefix[p] = id
	}
}
