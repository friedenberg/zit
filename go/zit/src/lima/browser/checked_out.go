package browser

import (
	"fmt"

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
	return *(ids.MustRepoId("browser"))
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

	c.State = checked_out_state.StateError
	c.Error = err
}

func (a *CheckedOut) String() string {
	return fmt.Sprintf("%s %s", &a.Internal, &a.External)
}

var CheckedOutResetter checkedOutResetter

type checkedOutResetter struct{}

func (checkedOutResetter) Reset(a *CheckedOut) {
	a.State = checked_out_state.StateUnknown
	a.IsImport = false
	a.Error = nil

	sku.TransactedResetter.Reset(&a.Internal)
	sku.TransactedResetter.Reset(&a.External.Transacted)
	// TODO reset item
}

func (checkedOutResetter) ResetWith(a *CheckedOut, b *CheckedOut) {
	a.State = b.State
	a.IsImport = b.IsImport
	a.Error = b.Error

	sku.TransactedResetter.ResetWith(&a.Internal, &b.Internal)
	sku.TransactedResetter.ResetWith(&a.External.Transacted, &b.External.Transacted)
	// TODO reset item
	// a.External.FDs.ResetWith(&b.External.FDs)
}
