package store_fs

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type CheckedOut struct {
	Internal sku.Transacted
	External External
	State    checked_out_state.State
	IsImport bool
	Error    error
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

func (c *CheckedOut) SetError(err error) {
	if err == nil {
		return
	}

	c.State = checked_out_state.StateError
	c.Error = err
}

func (c *CheckedOut) InternalAndExternalEqualsSansTai() bool {
	return c.External.Metadatei.EqualsSansTai(
		&c.Internal.Metadatei,
	)
}

func (c *CheckedOut) DetermineState(justCheckedOut bool) {
	if c.Internal.GetObjekteSha().IsNull() {
		c.State = checked_out_state.StateUntracked
	} else if c.Internal.Metadatei.EqualsSansTai(&c.External.Metadatei) {
		if justCheckedOut {
			c.State = checked_out_state.StateJustCheckedOut
		} else {
			c.State = checked_out_state.StateExistsAndSame
		}
	} else {
		if justCheckedOut {
			c.State = checked_out_state.StateJustCheckedOutButDifferent
		} else {
			c.State = checked_out_state.StateExistsAndDifferent
		}
	}
}

func (a *CheckedOut) String() string {
	return fmt.Sprintf("%s %s", &a.Internal, &a.External)
}

func (e *CheckedOut) Remove(s schnittstellen.Standort) (err error) {
	// TODO check conflict state
	if err = e.External.FDs.Objekte.Remove(s); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = e.External.FDs.Akte.Remove(s); err != nil {
		err = errors.Wrap(err)
		return
	}

	e.External.FDs.Akte.Reset()
	e.External.FDs.Objekte.Reset()

	return
}

func ToSliceFilesZettelen(
	s CheckedOutSet,
) (out []string, err error) {
	return iter.DerivedValues(
		s,
		func(z *CheckedOut) (e string, err error) {
			e = z.External.GetObjekteFD().GetPath()

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
	a.State = checked_out_state.StateUnknown
	a.IsImport = false
	a.Error = nil

	sku.TransactedResetter.Reset(&a.Internal)
	sku.TransactedResetter.Reset(&a.External.Transacted)
	a.External.FDs.Objekte.Reset()
	a.External.FDs.Akte.Reset()
}

func (checkedOutResetter) ResetWith(a *CheckedOut, b *CheckedOut) {
	a.State = b.State
	a.IsImport = b.IsImport
	a.Error = b.Error

	sku.TransactedResetter.ResetWith(&a.Internal, &b.Internal)
	sku.TransactedResetter.ResetWith(&a.External.Transacted, &b.External.Transacted)
	a.External.FDs.ResetWith(&b.External.FDs)
}
