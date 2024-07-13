package query

type Reducer interface {
	Reduce(*Builder) error
}
