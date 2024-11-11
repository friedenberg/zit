package sku

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/charlie/external_state"
	"code.linenisgreat.com/zit/go/zit/src/echo/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

func InternalAndExternalEqualsWithoutTai(co SkuType) bool {
	i := co.GetSku()
	e := co.GetSkuExternal().GetSku()

	return e.Metadata.EqualsSansTai(
		&i.Metadata,
	)
}

func DetermineState(
	c SkuType,
	justCheckedOut bool,
) {
	es := c.GetSkuExternal().GetExternalState()

	if es == external_state.Recognized {
		c.SetState(checked_out_state.Recognized)
		return
	}

	i := c.GetSku()
	e := c.GetSkuExternal().GetSku()

	if i.GetObjectSha().IsNull() {
		c.SetState(checked_out_state.Untracked)
	} else if i.Metadata.EqualsSansTai(&e.Metadata) {
		if justCheckedOut {
			c.SetState(checked_out_state.JustCheckedOut)
		} else {
			c.SetState(checked_out_state.ExistsAndSame)
		}
	} else {
		c.SetState(checked_out_state.Changed)
	}
}

type CheckedOut struct {
	Internal Transacted
	External Transacted
	State    checked_out_state.State
	IsImport bool
}

func (c *CheckedOut) GetRepoId() ids.RepoId {
	return c.External.RepoId
}

func (c *CheckedOut) GetSkuExternal() *Transacted {
	return &c.External
}

func (c *CheckedOut) GetSku() *Transacted {
	return &c.Internal
}

func (c *CheckedOut) GetState() checked_out_state.State {
	return c.State
}

func (src *CheckedOut) Clone() *CheckedOut {
	dst := GetCheckedOutPool().Get()
	CheckedOutResetter.ResetWith(dst, src)
	return dst
}

func (t *CheckedOut) GetExternalObjectId() ids.ExternalObjectId {
	return t.External.GetExternalObjectId()
}

func (t *CheckedOut) GetExternalState() external_state.State {
	return t.External.GetExternalState()
}

func (a *CheckedOut) GetObjectId() *ids.ObjectId {
	return a.External.GetObjectId()
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

func (a *CheckedOut) String() string {
	return fmt.Sprintf("%s %s", &a.Internal, &a.External)
}
