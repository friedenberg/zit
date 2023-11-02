package pool

var MapPool mapPool

func init() {
}

type mapPool struct {
	elements map[string]string
}
