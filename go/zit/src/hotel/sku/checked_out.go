package sku

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

func InternalAndExternalEqualsWithoutTai(col CheckedOutLike) bool {
	i := col.GetSku()
	e := col.GetSkuExternalLike().GetSku()

	return e.Metadata.EqualsSansTai(
		&i.Metadata,
	)
}

func DetermineState(c CheckedOutLike, justCheckedOut bool) {
	i := c.GetSku()
	e := c.GetSkuExternalLike().GetSku()

	if i.GetObjectSha().IsNull() {
		c.SetState(checked_out_state.StateUntracked)
	} else if i.Metadata.EqualsSansTai(&e.Metadata) {
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

func (c *CheckedOut) GetRepoId() ids.RepoId {
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
	return c.External.Metadata.EqualsSansTai(
		&c.Internal.Metadata,
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
	if c.Internal.GetObjectSha().IsNull() {
		c.State = checked_out_state.StateUntracked
	} else if c.Internal.Metadata.EqualsSansTai(&c.External.Metadata) {
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
