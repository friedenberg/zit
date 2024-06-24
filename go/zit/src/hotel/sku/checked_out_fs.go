package sku

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
)

type CheckedOutFS struct {
	Internal Transacted
	External ExternalFS
	State    checked_out_state.State
	IsImport bool
	Error    error
}

func (c *CheckedOutFS) GetSkuCheckedOutLike() CheckedOutLike {
	return c
}

func (c *CheckedOutFS) GetSkuExternalLike() ExternalLike {
	return &c.External
}

func (c *CheckedOutFS) GetSku() *Transacted {
	return &c.Internal
}

func (c *CheckedOutFS) GetState() checked_out_state.State {
	return c.State
}

func (c *CheckedOutFS) SetState(v checked_out_state.State) (err error) {
	c.State = v
	return
}

func (c *CheckedOutFS) GetError() error {
	return c.Error
}

func (c *CheckedOutFS) SetError(err error) {
	if err == nil {
		return
	}

	c.State = checked_out_state.StateError
	c.Error = err
}

func (c *CheckedOutFS) InternalAndExternalEqualsSansTai() bool {
	return c.External.Metadatei.EqualsSansTai(
		&c.Internal.Metadatei,
	)
}

func (c *CheckedOutFS) DetermineState(justCheckedOut bool) {
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

func (a *CheckedOutFS) String() string {
	return fmt.Sprintf("%s %s", &a.Internal, &a.External)
}

func (e *CheckedOutFS) Remove(s schnittstellen.Standort) (err error) {
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
