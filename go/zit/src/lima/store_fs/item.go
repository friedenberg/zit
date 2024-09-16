package store_fs

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/charlie/external_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/delta/thyme"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

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

	if src.MutableSetLike != nil {
		src.MutableSetLike.Each(dst.MutableSetLike.Add)
	}
}

func (a *Item) Equals(b *Item) bool {
	if !a.Object.Equals(&b.Object) {
		return false
	}

	if !a.Blob.Equals(&b.Blob) {
		return false
	}

	return true
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
func (i *Item) ReadFromExternal(el sku.ExternalLike) (err error) {
	e := el.(*External)
	i.ResetWith(&e.item)
	return
}

// TODO replace with fields
func (i *Item) WriteToExternal(e *External) (err error) {
	k := &i.ObjectId

	if err = e.ExternalObjectId.SetRaw(k.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if k.String() != "" {
		if k.GetGenre() != genres.Blob {
			if err = e.Transacted.ObjectId.Set(k.String()); err != nil {
				err = nil
				// err = errors.Wrap(err)
				// return
			}
		}
	}
	// e.Transacted.ObjectId.SetGenre(k.GetGenre())
	e.item.ResetWith(i) // TODO remove

	m := &e.Transacted.Metadata
	m.Tai = i.GetTai()

	if i.MutableSetLike != nil {
		if err = i.MutableSetLike.Each(
			func(f *fd.FD) (err error) {
				field := sku.Field{
					Value:     f.GetPath(),
					ColorType: string_format_writer.ColorTypeId,
				}

				switch {
				case f == &i.Object:
					field.Key = "object"

				case f == &i.Conflict:
					field.Key = "conflict"

				case f == &i.Blob:
					fallthrough

				default:
					field.Key = "blob"
				}

				e.Transacted.Fields = append(e.Transacted.Fields, field)

				return
			},
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
