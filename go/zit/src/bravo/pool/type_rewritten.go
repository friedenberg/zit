package pool

// type typeRewritten[FROM TO, TO interface{}] struct {
// 	inner interfaces.PoolValue[FROM]
// }

// func MakeTypeRewrittenPool[TO any, FROM interface{ TO }](p interfaces.PoolValue[FROM]) typeRewritten[FROM, TO] {
// 	return
// }

type ManualPool[T any] struct {
	FuncGet func() T
	FuncPut func(T)
}

func (ip ManualPool[T]) Get() T {
	return ip.FuncGet()
}

func (ip ManualPool[T]) Put(i T) (err error) {
	ip.FuncPut(i)
	return
}
