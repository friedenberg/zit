package sku

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/charlie/external_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/thyme"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

var ErrExternalHasConflictMarker = errors.New("external has conflict marker")

type FSItem struct {
	external_state.State

	// TODO refactor this to be a string and a genre that is tied to the state
	ids.ObjectId

	Object   fd.FD
	Blob     fd.FD // TODO make set
	Conflict fd.FD

	interfaces.MutableSetLike[*fd.FD]
}

func (ef *FSItem) String() string {
	return ef.ObjectId.String()
}

func (ef *FSItem) GetExternalObjectId() ExternalObjectId {
	return ef
}

func (i *FSItem) Debug() string {
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

func (i *FSItem) GetTai() ids.Tai {
	return ids.TaiFromTime(i.LatestModTime())
}

func (i *FSItem) GetTime() thyme.Time {
	return i.LatestModTime()
}

func (i *FSItem) LatestModTime() thyme.Time {
	o, b := i.Object.ModTime(), i.Blob.ModTime()

	if o.Less(b) {
		return b
	} else {
		return o
	}
}

func (dst *FSItem) Reset() {
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

func (dst *FSItem) ResetWith(src *FSItem) {
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

func (a *FSItem) Equals(b *FSItem) (ok bool, why string) {
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

func (e *FSItem) GenerateConflictFD() (err error) {
	if err = e.Conflict.SetPath(e.ObjectId.String() + ".conflict"); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (e *FSItem) GetCheckoutModeOrError() (m checkout_mode.Mode, err error) {
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