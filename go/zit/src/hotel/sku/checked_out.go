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
	internal Transacted
	External Transacted
	State    checked_out_state.State
	IsImport bool
}

func (c *CheckedOut) GetRepoId() ids.RepoId {
	return c.GetSkuExternal().RepoId
}

func (c *CheckedOut) GetSkuExternal() *Transacted {
	return &c.External
}

func (c *CheckedOut) GetSku() *Transacted {
	return &c.internal
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
	return t.GetSkuExternal().GetExternalObjectId()
}

func (t *CheckedOut) GetExternalState() external_state.State {
	return t.GetSkuExternal().GetExternalState()
}

func (a *CheckedOut) GetObjectId() *ids.ObjectId {
	return a.GetSkuExternal().GetObjectId()
}

func (c *CheckedOut) InternalAndExternalEqualsSansTai() bool {
	return c.GetSkuExternal().GetSku().GetMetadata().EqualsSansTai(
		&c.GetSku().Metadata,
	)
}

func (c *CheckedOut) SetState(v checked_out_state.State) (err error) {
	c.State = v
	return
}

func (a *CheckedOut) String() string {
	return fmt.Sprintf("%s %s", a.GetSku(), a.GetSkuExternal())
}
