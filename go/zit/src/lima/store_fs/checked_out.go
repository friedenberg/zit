package store_fs

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type CheckedOut struct {
	Internal sku.Transacted
	External External
	State    checked_out_state.State
	IsImport bool
	Error    error
}

func (c *CheckedOut) GetRepoId() ids.RepoId {
	return *(ids.MustRepoId(""))
}

func (c *CheckedOut) GetSkuCheckedOutLike() sku.CheckedOutLike {
	return c
}

func (c *CheckedOut) GetSkuExternalLike() sku.ExternalLike {
	return &c.External
}

func (c *CheckedOut) GetSku() *sku.Transacted {
	return &c.Internal
}

func (c *CheckedOut) GetState() checked_out_state.State {
	return c.State
}

func (c *CheckedOut) SetState(v checked_out_state.State) (err error) {
	c.State = v
	return
}

func (c *CheckedOut) GetError() error {
	return c.Error
}

func (a *CheckedOut) Clone() sku.CheckedOutLike {
	b := GetCheckedOutPool().Get()
	CheckedOutResetter.ResetWith(b, a)
	return b
}

func (c *CheckedOut) SetError(err error) {
	if err == nil {
		return
	}

	c.State = checked_out_state.Error
	c.Error = err
}

func (a *CheckedOut) String() string {
	return fmt.Sprintf("%s %s", &a.Internal, &a.External)
}

func (e *CheckedOut) Remove(s interfaces.Directory) (err error) {
	// TODO check conflict state
	if err = e.External.FDs.Object.Remove(s); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = e.External.FDs.Blob.Remove(s); err != nil {
		err = errors.Wrap(err)
		return
	}

	e.External.FDs.Blob.Reset()
	e.External.FDs.Object.Reset()

	return
}

func ToSliceFilesZettelen(
	s sku.CheckedOutLikeSet,
) (out []string, err error) {
	return iter.DerivedValues(
		s,
		func(col sku.CheckedOutLike) (e string, err error) {
			z := col.(*CheckedOut)
			e = z.External.GetObjectFD().GetPath()

			if e == "" {
				err = collections.MakeErrStopIteration()
				return
			}

			return
		},
	)
}

func ToSliceFilesBlobs(
	s sku.CheckedOutLikeSet,
) (out []string, err error) {
	return iter.DerivedValues(
		s,
		func(col sku.CheckedOutLike) (e string, err error) {
			z := col.(*CheckedOut)
			e = z.External.GetBlobFD().GetPath()

			if e == "" {
				err = collections.MakeErrStopIteration()
				return
			}

			return
		},
	)
}

var CheckedOutResetter checkedOutResetter

type checkedOutResetter struct{}

func (checkedOutResetter) Reset(a *CheckedOut) {
	a.State = checked_out_state.Unknown
	a.IsImport = false
	a.Error = nil

	sku.TransactedResetter.Reset(&a.Internal)
	sku.TransactedResetter.Reset(&a.External.Transacted)
	a.External.FDs.Object.Reset()
	a.External.FDs.Blob.Reset()
}

func (checkedOutResetter) ResetWith(a *CheckedOut, b *CheckedOut) {
	a.State = b.State
	a.IsImport = b.IsImport
	a.Error = b.Error

	sku.TransactedResetter.ResetWith(&a.Internal, &b.Internal)
	sku.TransactedResetter.ResetWith(&a.External.Transacted, &b.External.Transacted)
	a.External.FDs.ResetWith(&b.External.FDs)
}
