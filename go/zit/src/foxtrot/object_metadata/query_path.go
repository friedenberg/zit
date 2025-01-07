package object_metadata

import (
	"fmt"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
)

type qpItem struct {
	Ok  bool
	Any any
}

// TODO remove this entirely in favor of ad hoc debugging of queries
type QueryPath []qpItem

func (qp *QueryPath) String() string {
	var sb strings.Builder

	for _, i := range *qp {
		fmt.Fprintf(&sb, "%t: %s", i.Ok, i.Any)
	}

	return sb.String()
}

func (qp *QueryPath) Reset() {
	*qp = (*qp)[:0]
}

func (qp *QueryPath) Len() int {
	return len(*qp)
}

func (qp *QueryPath) PushOnReturn(q any, ok *bool) {
	if !ui.IsVerbose() {
		return
	}

	si, _ := errors.MakeStackInfo(1)
	ui.Log().Print("QueryPath", *ok, si.FileNameLine())
}

func (qp *QueryPath) Pop() any {
	l := qp.Len()
	q := (*qp)[l-1]
	*qp = (*qp)[0 : l-1]
	return q
}
