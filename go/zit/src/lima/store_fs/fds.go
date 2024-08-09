package store_fs

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/charlie/external_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/thyme"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type FDSet struct {
	external_state.State

	ids.ObjectId

	Object   fd.FD
	Blob     fd.FD
	Conflict fd.FD

	interfaces.MutableSetLike[*fd.FD]
}

func (ef *FDSet) String() string {
	return ef.ObjectId.String()
}

func (ef *FDSet) GetObjectId() *ids.ObjectId {
	return &ef.ObjectId
}

func (ef *FDSet) GetExternalObjectId() *ids.ObjectId {
	return &ef.ObjectId
}

func (ef *FDSet) Debug() string {
	return fmt.Sprintf(
		"ObjectId: %q, Object: %q, Blob: %q, Conflict: %q, All: %q",
		&ef.ObjectId,
		&ef.Object,
		&ef.Blob,
		&ef.Conflict,
		ef.MutableSetLike,
	)
}

func (ef *FDSet) GetTai() ids.Tai {
	return ids.TaiFromTime(ef.LatestModTime())
}

func (ef *FDSet) GetTime() thyme.Time {
	return ef.LatestModTime()
}

func (ef *FDSet) LatestModTime() thyme.Time {
	o, a := ef.Object.ModTime(), ef.Blob.ModTime()

	if o.Less(a) {
		return a
	} else {
		return o
	}
}

func (dst *FDSet) ResetWith(src *FDSet) {
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

func (a *FDSet) Equals(b *FDSet) bool {
	if !a.Object.Equals(&b.Object) {
		return false
	}

	if !a.Blob.Equals(&b.Blob) {
		return false
	}

	return true
}

func (e *FDSet) GenerateConflictFD() (err error) {
	if err = e.Conflict.SetPath(e.Object.String() + ".conflict"); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (e *FDSet) GetCheckoutModeOrError() (m checkout_mode.Mode, err error) {
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
