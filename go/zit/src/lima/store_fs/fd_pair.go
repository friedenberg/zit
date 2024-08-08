package store_fs

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/thyme"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

var ErrExternalHasConflictMarker = errors.New("external has conflict marker")

type FDPair struct {
	Object fd.FD
	Blob   fd.FD
}

func (ef *FDPair) GetTai() ids.Tai {
	return ids.TaiFromTime(ef.LatestModTime())
}

func (ef *FDPair) GetTime() thyme.Time {
	return ef.LatestModTime()
}

func (ef *FDPair) LatestModTime() thyme.Time {
	o, a := ef.Object.ModTime(), ef.Blob.ModTime()

	if o.Less(a) {
		return a
	} else {
		return o
	}
}

// func (dst *ExternalFDs) String() string {
//   panic("")
// }

func (dst *FDPair) ResetWith(src *FDPair) {
	dst.Object.ResetWith(&src.Object)
	dst.Blob.ResetWith(&src.Blob)
}

func (a *FDPair) Equals(b *FDPair) bool {
	if !a.Object.Equals(&b.Object) {
		return false
	}

	if !a.Blob.Equals(&b.Blob) {
		return false
	}

	return true
}

func (e *FDPair) MakeConflictMarker() (path string) {
	path = fmt.Sprintf("%s.conflict", e.Object.GetPath())

	return
}

func (e *FDPair) conflictMarkerExists(fd *fd.FD) (ok bool) {
	if files.Exists(fmt.Sprintf("%s.conflict", fd)) {
		ok = true
	}

	return
}

func (e *FDPair) ConflictMarkerError() (err error) {
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

func (e *FDPair) GetCheckoutModeOrError() (m checkout_mode.Mode, err error) {
	switch {
	case !e.Object.IsEmpty() && !e.Blob.IsEmpty():
		m = checkout_mode.MetadataAndBlob

	case !e.Blob.IsEmpty():
		m = checkout_mode.BlobOnly

	case !e.Object.IsEmpty():
		m = checkout_mode.MetadataOnly

	default:
		err = checkout_mode.MakeErrInvalidCheckoutMode(
			errors.Errorf("all FD's are empty"),
		)
	}

	return
}

func (e *FDPair) GetCheckoutMode() (m checkout_mode.Mode) {
	switch {
	case !e.Object.IsEmpty() && !e.Blob.IsEmpty():
		m = checkout_mode.MetadataAndBlob

	case !e.Blob.IsEmpty():
		m = checkout_mode.BlobOnly

	case !e.Object.IsEmpty():
		m = checkout_mode.MetadataOnly
	}

	return
}
