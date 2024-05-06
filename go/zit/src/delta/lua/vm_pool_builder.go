package lua

type VMPoolBuilder struct {
	proto VMPool
}

func (vpb *VMPoolBuilder) WithRequire(v LGFunction) *VMPoolBuilder {
	vpb.proto.Require = v
	return vpb
}

func (vpb *VMPoolBuilder) WithSearcher(v LGFunction) *VMPoolBuilder {
	vpb.proto.Searcher = v
	return vpb
}

func (vpb *VMPoolBuilder) Build() *VMPool {
	return &VMPool{
		Require:  vpb.proto.Require,
		Searcher: vpb.proto.Searcher,
	}
}
