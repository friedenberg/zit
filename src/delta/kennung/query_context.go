package kennung

// type QueryContext[T QueryKennung[T], TPtr QueryKennungPtr[T]] struct {
// 	Expander func(string) (string, error)
// 	Include  schnittstellen.Set[T]
// 	Exclude  schnittstellen.Set[T]
// }

// func (qc QueryContext[T, TPtr]) MakeQuerySet() QuerySet[T, TPtr] {
// }

type Expanders struct {
	Sha, Etikett, Hinweis, Typ, Kasten func(string) (string, error)
}
