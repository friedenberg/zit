package query

import (
	"strings"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/bravo/log"
	"code.linenisgreat.com/zit/src/charlie/gattung"
	"code.linenisgreat.com/zit/src/echo/fd"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

type Kennung struct {
	Exact   bool
	Virtual bool
	Debug   bool
	FD      *fd.FD
	*kennung.Kennung2
}

func (k Kennung) Reduce(b *Builder) (err error) {
	if err = k.Expand(b.expanders); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// TODO support exact
func (k Kennung) ContainsMatchable(sk *sku.Transacted) bool {
	me := sk.GetMetadatei()
	switch k.GetGattung() {
	case gattung.Etikett:
		s := k.String()

		log.Log().Printf("%q in %q", k.Kennung2, iter.StringCommaSeparated(me.GetEtiketten()))

		if me.GetEtiketten().ContainsKey(s) {
			return true
		}

		if me.Verzeichnisse.GetExpandedEtiketten().ContainsKey(s) {
			return true
		}

		if me.Verzeichnisse.GetImplicitEtiketten().ContainsKey(s) {
			return true
		}

	case gattung.Typ:
		log.Log().Printf("%q -> %q", k.Kennung2, sk.GetKennung())

		if kennung.Contains(me.GetTyp(), k) {
			return true
		}

		// case kennung.ShaLike:
		// 	if Sha(kt.GetSha()).ContainsMatchable(m) {
		// 		return true
		// 	}
	}

	idl := &sk.Kennung

	if !kennung.Contains(idl, k) {
		return false
	}

	return true
}

func (k Kennung) String() string {
	var sb strings.Builder

	if k.Exact {
		sb.WriteRune('=')
	}

	if k.Virtual {
		sb.WriteRune('%')
	}

	sb.WriteString(kennung.FormattedString(k.Kennung2))

	return sb.String()
}

func (k Kennung) Each(_ schnittstellen.FuncIter[Matcher]) error {
	return nil
}

func (k Kennung) MatcherLen() int {
	return 1
}
