package sku

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
)

func DetermineState(c CheckedOutLike, justCheckedOut bool) {
	i := c.GetSku()
	e := c.GetSkuExternalLike().GetSku()

	if i.GetObjekteSha().IsNull() {
		c.SetState(checked_out_state.StateUntracked)
	} else if i.Metadatei.EqualsSansTai(&e.Metadatei) {
		if justCheckedOut {
			c.SetState(checked_out_state.StateJustCheckedOut)
		} else {
			c.SetState(checked_out_state.StateExistsAndSame)
		}
	} else {
		if justCheckedOut {
			c.SetState(checked_out_state.StateJustCheckedOutButDifferent)
		} else {
			c.SetState(checked_out_state.StateExistsAndDifferent)
		}
	}
}

type CheckedOut struct {
	Internal Transacted
	External External
	State    checked_out_state.State
	Error    error
}

func (c *CheckedOut) GetKasten() kennung.Kasten {
	panic(todo.Implement())
}

func (c *CheckedOut) GetSkuCheckedOutLike() CheckedOutLike {
	return c
}

func (c *CheckedOut) GetSkuExternalLike() ExternalLike {
	return &c.External
}

func (c *CheckedOut) GetSku() *Transacted {
	return &c.Internal
}

func (c *CheckedOut) GetState() checked_out_state.State {
	return c.State
}

func (c *CheckedOut) Clone() CheckedOutLike {
	panic(todo.Implement())
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

func (c *CheckedOut) SetState(v checked_out_state.State) (err error) {
	c.State = v
	return
}

func (c *CheckedOut) GetError() error {
	return c.Error
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

// func (e *CheckedOut) Remove(s schnittstellen.Standort) (err error) {
// 	// TODO check conflict state
// 	if err = e.External.FDs.Objekte.Remove(s); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	if err = e.External.FDs.Akte.Remove(s); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	e.External.FDs.Akte.Reset()
// 	e.External.FDs.Objekte.Reset()

// 	return
// }
