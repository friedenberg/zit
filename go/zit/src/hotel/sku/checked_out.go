package sku

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/charlie/external_state"
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

func DetermineState(
	c CheckedOutLike,
	justCheckedOut bool,
) {
	es := c.GetSkuExternalLike().GetExternalState()

	if es == external_state.Recognized {
		c.SetState(checked_out_state.Recognized)
		return
	}

	i := c.GetSku()
	e := c.GetSkuExternalLike().GetSku()

	if i.GetObjectSha().IsNull() {
		c.SetState(checked_out_state.Untracked)
	} else if i.Metadata.EqualsSansTai(&e.Metadata) {
		if justCheckedOut {
			c.SetState(checked_out_state.JustCheckedOut)
		} else {
			c.SetState(checked_out_state.ExistsAndSame)
		}
	} else {
		if justCheckedOut {
			c.SetState(checked_out_state.JustCheckedOutButDifferent)
		} else {
			c.SetState(checked_out_state.ExistsAndDifferent)
		}
	}
}

type CheckedOut struct {
	Internal Transacted
	External ExternalLike
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
	return c.External
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

	c.State = checked_out_state.Error
	c.Error = err
}

func (c *CheckedOut) InternalAndExternalEqualsSansTai() bool {
	return c.External.GetSku().GetMetadata().EqualsSansTai(
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

func (a *CheckedOut) String() string {
	return fmt.Sprintf("%s %s", &a.Internal, a.External)
}
