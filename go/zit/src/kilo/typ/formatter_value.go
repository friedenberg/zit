package typ

import (
	"io"
	"strings"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/log"
	"code.linenisgreat.com/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/src/delta/lua"
	"code.linenisgreat.com/zit/src/delta/typ_akte"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/juliett/objekte"
)

type FormatterValue struct {
	string
}

func (f FormatterValue) String() string {
	return f.string
}

func (f *FormatterValue) Set(v string) (err error) {
	v1 := strings.TrimSpace(strings.ToLower(v))
	switch v1 {
	case "action-names", "hooks.on_pre_commit":
		f.string = v1

	default:
		err = objekte.MakeErrUnsupportedFormatterValue(v1, gattung.Typ)
		return
	}

	return
}

func (f *FormatterValue) FuncFormatter(
	out io.Writer,
	af schnittstellen.AkteIOFactory,
	agp schnittstellen.AkteGetterPutter[*typ_akte.V0],
) schnittstellen.FuncIter[*sku.Transacted] {
	switch f.string {
	case "action-names":
		f := typ_akte.MakeFormatterActionNames()

		return func(o *sku.Transacted) (err error) {
			var akte *typ_akte.V0

			if akte, err = agp.GetAkte(o.GetAkteSha()); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer agp.PutAkte(akte)

			if _, err = f.Format(out, akte); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "hooks.on_pre_commit":
		return func(o *sku.Transacted) (err error) {
			var akte *typ_akte.V0

			if akte, err = agp.GetAkte(o.GetAkteSha()); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer agp.PutAkte(akte)

			script, ok := akte.Hooks.(string)

			if !ok || script == "" {
				return
			}

			var vp lua.VMPool

			if err = vp.Set(script); err != nil {
				err = errors.Wrap(err)
				return
			}

			vm := vp.Get()
			defer vp.Put(vm)

			var tt *lua.LTable

			if tt, err = vm.GetTopTableOrError(); err != nil {
				err = errors.Wrap(err)
				return
			}

			f := vm.GetField(tt, "on_pre_commit")

			log.Out().Print(f.String())

			return
		}

	default:
		return func(_ *sku.Transacted) (err error) {
			err = objekte.MakeErrUnsupportedFormatterValue(
				f.string,
				gattung.Typ,
			)
			return
		}
	}
}
