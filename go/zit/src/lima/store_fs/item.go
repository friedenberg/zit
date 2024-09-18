package store_fs

import (
	"fmt"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/charlie/external_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/delta/thyme"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

var ErrExternalHasConflictMarker = errors.New("external has conflict marker")

type Item struct {
	external_state.State

	// TODO refactor this to be a string and a genre that is tied to the state
	ids.ObjectId

	Object   fd.FD
	Blob     fd.FD // TODO make set
	Conflict fd.FD

	interfaces.MutableSetLike[*fd.FD]
}

func (ef *Item) String() string {
	return ef.ObjectId.String()
}

func (ef *Item) GetExternalObjectId() sku.ExternalObjectId {
	return ef
}

func (i *Item) Debug() string {
	return fmt.Sprintf(
		"State: %q, Genre: %q, ObjectId: %q, Object: %q, Blob: %q, Conflict: %q, All: %q",
		i.State,
		i.GetGenre(),
		&i.ObjectId,
		&i.Object,
		&i.Blob,
		&i.Conflict,
		i.MutableSetLike,
	)
}

func (i *Item) GetTai() ids.Tai {
	return ids.TaiFromTime(i.LatestModTime())
}

func (i *Item) GetTime() thyme.Time {
	return i.LatestModTime()
}

func (i *Item) LatestModTime() thyme.Time {
	o, b := i.Object.ModTime(), i.Blob.ModTime()

	if o.Less(b) {
		return b
	} else {
		return o
	}
}

func (dst *Item) Reset() {
	dst.State = 0
	dst.ObjectId.Reset()
	dst.Object.Reset()
	dst.Blob.Reset()
	dst.Conflict.Reset()

	if dst.MutableSetLike == nil {
		dst.MutableSetLike = collections_value.MakeMutableValueSet[*fd.FD](nil)
	} else {
		dst.MutableSetLike.Reset()
	}
}

func (dst *Item) ResetWith(src *Item) {
	if dst == src {
		return
	}

	dst.State = src.State
	dst.ObjectId.ResetWith(&src.ObjectId)
	dst.Object.ResetWith(&src.Object)
	dst.Blob.ResetWith(&src.Blob)
	dst.Conflict.ResetWith(&src.Conflict)

	if dst.MutableSetLike == nil {
		dst.MutableSetLike = collections_value.MakeMutableValueSet[*fd.FD](nil)
	}

	dst.MutableSetLike.Reset()

	if src.MutableSetLike != nil {
		src.MutableSetLike.Each(dst.MutableSetLike.Add)
	}

	// TODO consider if this approach actually works
	if !dst.Object.IsEmpty() {
		dst.Add(&dst.Object)
	}

	if !dst.Blob.IsEmpty() {
		dst.Add(&dst.Blob)
	}

	if !dst.Conflict.IsEmpty() {
		dst.Add(&dst.Conflict)
	}
}

func (a *Item) Equals(b *Item) (ok bool, why string) {
	if ok, why = a.Object.Equals2(&b.Object); !ok {
		return false, fmt.Sprintf("Object.%s", why)
	}

	if ok, why = a.Blob.Equals2(&b.Blob); !ok {
		return false, fmt.Sprintf("Blob.%s", why)
	}

	if ok, why = a.Conflict.Equals2(&b.Conflict); !ok {
		return false, fmt.Sprintf("Conflict.%s", why)
	}

	if !iter.SetEquals(a.MutableSetLike, b.MutableSetLike) {
		return false, "set"
	}

	return
}

func (e *Item) GenerateConflictFD() (err error) {
	if err = e.Conflict.SetPath(e.ObjectId.String() + ".conflict"); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (e *Item) GetCheckoutModeOrError() (m checkout_mode.Mode, err error) {
	switch {
	case !e.Object.IsEmpty() && !e.Blob.IsEmpty():
		m = checkout_mode.MetadataAndBlob

	case !e.Blob.IsEmpty():
		m = checkout_mode.BlobOnly

	case !e.Object.IsEmpty():
		m = checkout_mode.MetadataOnly

	default:
		if e.State == external_state.Recognized {
			m = checkout_mode.BlobRecognized
			return
		}

		err = checkout_mode.MakeErrInvalidCheckoutMode(
			errors.Errorf("all FD's are empty: %s", e.Debug()),
		)
	}

	return
}

// TODO replace with fields
func (s *Store) ReadFromExternal(el sku.ExternalLike) (i *Item, err error) {
	i = &Item{} // TODO use pool or use dir_items?
	i.Reset()

	e := el.(*External)

	// TODO handle sort order
	for _, f := range e.Transacted.Fields {
		var fdee *fd.FD
		switch strings.ToLower(f.Key) {
		case "object":
			fdee = &i.Object

		case "blob":
			fdee = &i.Blob

		case "conflict":
			fdee = &i.Conflict

		default:
			err = errors.Errorf("unexpected field: %#v", f)
			return
		}

		// if we've already set one of object, blob, or conflict, don't set it again
		// and instead add a new FD to the item
		if !fdee.IsEmpty() {
			fdee = &fd.FD{}
		}

		if err = fdee.SetIgnoreNotExists(f.Value); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = i.Add(fdee); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	i.State = e.State
	i.ObjectId.ResetWith(&e.ExternalObjectId)
	// i.ObjectId.ResetWith(&e.ExternalObjectId)

	if ok, why := i.Equals(&e.item); !ok {
		err = errors.BadRequestf(
			"expected:\n%s\nactual:\n%s\nwhy: %q\nfields: %#v",
			e.item.Debug(),
			i.Debug(),
			why,
			e.Transacted.Fields,
		)

		return
	}

	return
}

// TODO replace with fields
func (s *Store) WriteToExternal(i *Item, el sku.ExternalLike) (err error) {
	e := el.(*External)
	e.Transacted.Fields = e.Transacted.Fields[:0]
	k := &i.ObjectId

	e.ExternalObjectId.ResetWith(k)

	if e.ExternalObjectId.String() != k.String() {
		err = errors.Errorf("expected %q but got %q", k, &e.ExternalObjectId)
	}

	// if k.String() != "" {
	// 	if k.GetGenre() != genres.Blob {
	// 		if err = e.Transacted.ObjectId.Set(k.String()); err != nil {
	// 			err = nil
	// 			// err = errors.Wrap(err)
	// 			// return
	// 		}
	// 	}
	// }
	// e.Transacted.ObjectId.SetGenre(k.GetGenre())
	e.item.ResetWith(i) // TODO remove

	m := &e.Transacted.Metadata
	m.Tai = i.GetTai()

	fdees := iter.SortedValues(i.MutableSetLike)

	for _, f := range fdees {
		field := sku.Field{
			Value:     f.GetPath(),
			ColorType: string_format_writer.ColorTypeId,
		}

		switch {
		case i.Object.Equals(f):
			field.Key = "object"

		case i.Conflict.Equals(f):
			field.Key = "conflict"

		case i.Blob.Equals(f):
			fallthrough

		default:
			field.Key = "blob"
		}

		e.Transacted.Fields = append(e.Transacted.Fields, field)
	}

	return
}
