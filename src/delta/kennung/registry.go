package kennung

import (
	"encoding/gob"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
)

var (
	registerOnce        sync.Once
	registryLock        *sync.Mutex
	registryGattung     map[gattung.Gattung]IdLike
	registryQueryPrefix map[string]IdLike
)

func once() {
	registryLock = &sync.Mutex{}
	registryGattung = make(map[gattung.Gattung]IdLike)
	registryQueryPrefix = make(map[string]IdLike)
}

func register(id IdLike) {
	gob.Register(id)
	registerOnce.Do(once)

	registryLock.Lock()
	defer registryLock.Unlock()

	ok := false
	var id1 IdLike
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

	if idQueryPrefix, ok := id.(QueryPrefixer); ok {
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
