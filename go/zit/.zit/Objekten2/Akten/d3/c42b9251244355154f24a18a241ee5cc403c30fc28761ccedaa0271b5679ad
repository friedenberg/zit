package object_metadata

import (
	"fmt"
	"strings"
)

type qpItem struct {
	Ok  bool
	Any any
}

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

func (qp *QueryPath) Push(q any) (err error) {
	*qp = append(*qp, qpItem{Any: q})
	return
}

func (qp *QueryPath) PushOnReturn(q any, ok *bool) {
	// if !*ok {
	// 	return
	// }

	qp.Push(qpItem{Ok: *ok, Any: q})
}

func (qp *QueryPath) Pop() any {
	l := qp.Len()
	q := (*qp)[l-1]
	*qp = (*qp)[0 : l-1]
	return q
}
