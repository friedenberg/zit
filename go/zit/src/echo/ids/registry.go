package ids

import (
	"encoding/gob"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_ptr"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
)

var (
	registerOnce   sync.Once
	registryLock   *sync.Mutex
	registryGenres map[genres.Genre]interfaces.ObjectId
)

func once() {
	registryLock = &sync.Mutex{}
	registryGenres = make(map[genres.Genre]interfaces.ObjectId)
}

func register[T interfaces.ObjectId, TPtr interface {
	interfaces.StringSetterPtr[T]
	interfaces.ObjectId
}](id T,
) {
	gob.Register(&id)
	gob.Register(collections_ptr.MakeMutableValueSet[T, TPtr](nil))
	gob.Register(collections_ptr.MakeValueSet[T, TPtr](nil))
	registerOnce.Do(once)

	registryLock.Lock()
	defer registryLock.Unlock()

	ok := false
	var id1 interfaces.ObjectId
	g := genres.Must(id.GetGenre())

	if id1, ok = registryGenres[g]; ok {
		panic(
			errors.ErrorWithStackf(
				"genre %s has two registrations: %s (old), %s (new)",
				g,
				id1,
				id,
			),
		)
	}

	registryGenres[g] = id
}
