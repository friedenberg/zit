package store_fs

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/thyme"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type FDSet struct {
	ids.ObjectId

	Object   fd.FD
	Blob     fd.FD
	Conflict fd.FD

	interfaces.MutableSetLike[*fd.FD]
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
	dst.Object.ResetWith(&src.Object)
	dst.Blob.ResetWith(&src.Blob)
	dst.Conflict.ResetWith(&src.Conflict)
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

func (e *FDSet) MakeConflictMarker() (path string) {
	path = fmt.Sprintf("%s.conflict", e.Object.GetPath())

	return
}

func (e *FDSet) conflictMarkerExists(fd *fd.FD) (ok bool) {
	if files.Exists(fmt.Sprintf("%s.conflict", fd)) {
		ok = true
	}

	return
}

func (e *FDSet) ConflictMarkerError() (err error) {
	if e.conflictMarkerExists(&e.Object) {
		err = errors.Wrapf(ErrExternalHasConflictMarker, "Object: %s", &e.Object)
		return
	}

	if e.conflictMarkerExists(&e.Blob) {
		err = errors.Wrapf(ErrExternalHasConflictMarker, "Blob: %s", &e.Blob)
		return
	}

	return
}

func (e *FDSet) GetCheckoutModeOrError() (m checkout_mode.Mode, err error) {
	switch {
	case !e.Object.IsEmpty() && !e.Blob.IsEmpty():
		m = checkout_mode.ModeMetadataAndBlob

	case !e.Blob.IsEmpty():
		m = checkout_mode.ModeBlobOnly

	case !e.Object.IsEmpty():
		m = checkout_mode.ModeMetadataOnly

	default:
		err = checkout_mode.MakeErrInvalidCheckoutMode(
			errors.Errorf("all FD's are empty"),
		)
	}

	return
}

func (e *FDSet) GetCheckoutMode() (m checkout_mode.Mode) {
	switch {
	case !e.Object.IsEmpty() && !e.Blob.IsEmpty():
		m = checkout_mode.ModeMetadataAndBlob

	case !e.Blob.IsEmpty():
		m = checkout_mode.ModeBlobOnly

	case !e.Object.IsEmpty():
		m = checkout_mode.ModeMetadataOnly
	}

	return
}
