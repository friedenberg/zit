package object_metadata

type QueryPath []any

func (qp *QueryPath) Reset() {
	*qp = (*qp)[:0]
}

func (qp *QueryPath) Len() int {
	return len(*qp)
}

func (qp *QueryPath) Push(q any) (err error) {
	*qp = append(*qp, q)
	return
}

func (qp *QueryPath) PushOnOk(q any, ok *bool) {
	if !*ok {
		return
	}

	qp.Push(q)
}

func (qp *QueryPath) Pop() any {
	l := qp.Len()
	q := (*qp)[l-1]
	*qp = (*qp)[0 : l-1]
	return q
}
