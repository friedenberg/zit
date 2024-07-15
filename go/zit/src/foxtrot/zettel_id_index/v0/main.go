package zettel_id_index

import (
	"bufio"
	"encoding/gob"
	"io"
	"math/rand"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/coordinates"
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/object_id_provider"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type encodedIds struct {
	AvailableIds map[int]bool
}

type oldIndex struct {
	su interfaces.CacheIOFactory

	lock sync.Mutex
	path string

	encodedIds

	oldZettelIdStore *object_id_provider.Provider

	didRead    bool
	hasChanges bool

	nonRandomSelection bool
}

func MakeIndex(
	k interfaces.Config,
	s interfaces.Directory,
	su interfaces.CacheIOFactory,
) (i *oldIndex, err error) {
	i = &oldIndex{
		path:               s.FileVerzeichnisseKennung(),
		nonRandomSelection: k.UsePredictableHinweisen(),
		su:                 su,
		encodedIds: encodedIds{
			AvailableIds: make(map[int]bool),
		},
	}

	if i.oldZettelIdStore, err = object_id_provider.New(s); err != nil {
		if errors.IsNotExist(err) {
			errors.TodoP4("determine which layer handles no-create kasten")
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (i *oldIndex) Flush() (err error) {
	i.lock.Lock()
	defer i.lock.Unlock()

	if !i.hasChanges {
		ui.Log().Print("no changes")
		return
	}

	var w1 io.WriteCloser

	if w1, err = i.su.WriteCloserCache(i.path); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, w1)

	w := bufio.NewWriter(w1)

	defer errors.DeferredFlusher(&err, w)

	enc := gob.NewEncoder(w)

	if err = enc.Encode(i.encodedIds); err != nil {
		err = errors.Wrapf(err, "failed to write encoded kennung")
		return
	}

	return
}

func (i *oldIndex) readIfNecessary() (err error) {
	i.lock.Lock()
	defer i.lock.Unlock()

	if i.didRead {
		return
	}

	ui.Log().Print("reading")

	i.didRead = true

	var r1 io.ReadCloser

	if r1, err = i.su.ReadCloserCache(i.path); err != nil {
		if errors.IsNotExist(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	defer r1.Close()

	r := bufio.NewReader(r1)

	dec := gob.NewDecoder(r)

	if err = dec.Decode(&i.encodedIds); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *oldIndex) Reset() (err error) {
	i.lock.Lock()
	defer i.lock.Unlock()

	lMax := i.oldZettelIdStore.Left().Len() - 1
	rMax := i.oldZettelIdStore.Right().Len() - 1

	if lMax == 0 {
		err = errors.Errorf("left kennung are empty")
		return
	}

	if rMax == 0 {
		err = errors.Errorf("right kennung are empty")
		return
	}

	i.AvailableIds = make(map[int]bool, lMax*rMax)

	for l := 0; l <= lMax; l++ {
		for r := 0; r <= rMax; r++ {
			k := &coordinates.ZettelIdCoordinate{
				Left:  coordinates.Int(l),
				Right: coordinates.Int(r),
			}

			ui.Log().Print(k)

			n := int(k.Id())
			i.AvailableIds[n] = true
		}
	}

	i.hasChanges = true

	return
}

func (i *oldIndex) AddHinweis(k1 ids.IdLike) (err error) {
	if !k1.GetGenre().EqualsGenre(genres.Zettel) {
		err = genres.MakeErrUnsupportedGenre(k1)
		return
	}

	var h ids.ZettelId

	if err = h.Set(k1.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = i.readIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var left, right int

	if left, err = i.oldZettelIdStore.Left().ZettelId(h.GetHead()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if right, err = i.oldZettelIdStore.Right().ZettelId(h.GetTail()); err != nil {
		err = errors.Wrap(err)
		return
	}

	k := coordinates.ZettelIdCoordinate{
		Left:  coordinates.Int(left),
		Right: coordinates.Int(right),
	}

	n := k.Id()
	ui.Log().Printf("deleting %d, %s", n, h)

	i.lock.Lock()
	defer i.lock.Unlock()

	delete(i.AvailableIds, int(n))

	i.hasChanges = true

	return
}

func (i *oldIndex) CreateHinweis() (h *ids.ZettelId, err error) {
	if err = i.readIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if len(i.AvailableIds) == 0 {
		err = errors.Errorf("no available kennung")
		return
	}

	if len(i.AvailableIds) == 0 {
		err = errors.Wrap(object_id_provider.ErrZettelIdsExhausted{})
		return
	}

	ri := 0

	if len(i.AvailableIds) > 1 {
		ri = rand.Intn(len(i.AvailableIds) - 1)
	}

	m := 0
	j := 0

	for n := range i.AvailableIds {
		if i.nonRandomSelection {
			if m == 0 {
				m = n
				continue
			}

			if n > m {
				continue
			}

			m = n
		} else {
			j++
			m = n

			if j == ri {
				break
			}
		}
	}

	delete(i.AvailableIds, int(m))

	i.hasChanges = true

	return i.makeZettelIdButDontStore(m)
}

func (i *oldIndex) makeZettelIdButDontStore(
	j int,
) (h *ids.ZettelId, err error) {
	k := &coordinates.ZettelIdCoordinate{}
	k.SetInt(coordinates.Int(j))

	h, err = ids.MakeZettelIdFromProvidersAndCoordinates(
		k.Id(),
		i.oldZettelIdStore.Left(),
		i.oldZettelIdStore.Right(),
	)
	if err != nil {
		err = errors.Wrapf(err, "trying to make hinweis for %s", k)
		return
	}

	return
}

func (i *oldIndex) PeekHinweisen(m int) (hs []*ids.ZettelId, err error) {
	if err = i.readIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if m > len(i.AvailableIds) || m == 0 {
		m = len(i.AvailableIds)
	}

	hs = make([]*ids.ZettelId, 0, m)
	j := 0

	for n := range i.AvailableIds {
		k := &coordinates.ZettelIdCoordinate{}
		k.SetInt(coordinates.Int(n))

		var h *ids.ZettelId

		if h, err = i.makeZettelIdButDontStore(n); err != nil {
			err = errors.Wrap(err)
			return
		}

		hs = append(hs, h)

		j++

		if j == m {
			break
		}
	}

	return
}
