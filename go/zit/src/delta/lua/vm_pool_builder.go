package lua

import (
	"io"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
)

type VMPoolBuilder struct {
	proto        VMPool
	scriptReader io.Reader
	apply        schnittstellen.FuncIter[*VM]
}

func (vpb *VMPoolBuilder) Clone() *VMPoolBuilder {
	clone := *vpb
	return &clone
}

func (vpb *VMPoolBuilder) WithRequire(v LGFunction) *VMPoolBuilder {
	vpb.proto.Require = v
	return vpb
}

func (vpb *VMPoolBuilder) WithSearcher(v LGFunction) *VMPoolBuilder {
	vpb.proto.Searcher = v
	return vpb
}

func (sp *VMPoolBuilder) WithScript(
	script string,
) *VMPoolBuilder {
	sp.scriptReader = strings.NewReader(script)
	return sp
}

func (sp *VMPoolBuilder) WithReader(
	r io.Reader,
) *VMPoolBuilder {
	sp.scriptReader = r
	return sp
}

func (sp *VMPoolBuilder) WithApply(
	apply schnittstellen.FuncIter[*VM],
) *VMPoolBuilder {
	sp.apply = apply
	return sp
}

func (vpb *VMPoolBuilder) Build() (vm *VMPool, err error) {
	vm = &VMPool{
		Require:  vpb.proto.Require,
		Searcher: vpb.proto.Searcher,
	}

	if vpb.scriptReader == nil {
		err = errors.Errorf("no script or reader set")
		return
	}

	if err = vm.SetReader(vpb.scriptReader, vpb.apply); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func MakeVMPoolWithZitSearcher(
	script string,
	searcher LGFunction,
	apply schnittstellen.FuncIter[*VM],
) (ml *VMPool, err error) {
	b := (&VMPoolBuilder{}).WithSearcher(searcher).WithScript(script).WithApply(apply)

	if ml, err = b.Build(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func MakeVMPoolWithZitRequire(
	script string,
	require LGFunction,
	apply schnittstellen.FuncIter[*VM],
) (ml *VMPool, err error) {
	b := (&VMPoolBuilder{}).WithRequire(require).WithScript(script).WithApply(apply)

	if ml, err = b.Build(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
