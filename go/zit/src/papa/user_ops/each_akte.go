package user_ops

import (
	"fmt"
	"os/exec"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
	"github.com/google/shlex"
)

// TODO move to store_fs
type EachBlob struct {
	*env.Local
	Utility string
}

func (c EachBlob) Run(
	zsc sku.SkuTypeSet,
) (err error) {
	if zsc.Len() == 0 {
		return
	}

	var blob_store []string

	for col := range zsc.All() {
		var fds *sku.FSItem

		if fds, err = c.GetStore().GetStoreFS().ReadFSItemFromExternal(
			col.GetSkuExternal(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		blob_store = append(blob_store, fds.Blob.GetPath())
	}

	v := fmt.Sprintf("running utility: %q", c.Utility)

	if err = c.PrinterHeader()(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	var args []string

	if args, err = shlex.Split(c.Utility); err != nil {
		err = errors.Wrap(err)
		return
	}

	args = append(args, blob_store...)

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = c.Out()
	cmd.Stdin = c.In()
	cmd.Stderr = c.Err()

	if err = cmd.Start(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
