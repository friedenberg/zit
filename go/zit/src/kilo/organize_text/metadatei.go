package organize_text

import (
	"bufio"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/ohio"
	"code.linenisgreat.com/zit/go/zit/src/echo/format"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type Metadatei struct {
	// metadatei.Metadatei
	kennung.EtikettSet
	Matchers interfaces.SetLike[sku.Query]
	Comments []string
	Typ      kennung.Typ
}

func (m Metadatei) RemoveFromTransacted(sk *sku.Transacted) (err error) {
	mes := sk.Metadatei.GetEtiketten().CloneMutableSetPtrLike()

	if err = m.Each(mes.Del); err != nil {
		err = errors.Wrap(err)
		return
	}

	sk.Metadatei.SetEtiketten(mes)

	return
}

func (m Metadatei) AsMetadatei() (m1 metadatei.Metadatei) {
	m1.Typ = m.Typ
	m1.SetEtiketten(m.EtikettSet)
	return
}

func (m Metadatei) HasMetadateiContent() bool {
	if m.Len() > 0 {
		return true
	}

	tString := m.Typ.String()

	if tString != "" {
		return true
	}

	return false
}

func (m *Metadatei) ReadFrom(r1 io.Reader) (n int64, err error) {
	r := bufio.NewReader(r1)

	mes := kennung.MakeEtikettMutableSet()

	if n, err = format.ReadLines(
		r,
		ohio.MakeLineReaderRepeat(
			ohio.MakeLineReaderKeyValues(
				map[string]interfaces.FuncSetString{
					"%": func(v string) (err error) {
						m.Comments = append(m.Comments, v)
						return
					},
					"-": iter.MakeFuncSetString(mes),
					"!": m.Typ.Set,
				},
			),
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	m.EtikettSet = mes.CloneSetPtrLike()

	return
}

func (m Metadatei) WriteTo(w1 io.Writer) (n int64, err error) {
	w := format.NewLineWriter()

	for _, e := range iter.SortedStrings(m.EtikettSet) {
		w.WriteFormat("- %s", e)
	}

	tString := m.Typ.String()

	if tString != "" {
		w.WriteFormat("! %s", tString)
	}

	if m.Matchers != nil {
		for _, c := range iter.SortedStrings(m.Matchers) {
			w.WriteFormat("%% Matcher:%s", c)
		}
	}

	for _, c := range m.Comments {
		w.WriteFormat("%% %s", c)
	}

	return w.WriteTo(w1)
}

func (m Metadatei) GetOptionComments(
	f optionCommentFactory,
) (ocs []Option, err error) {
	em := errors.MakeMulti()

	for _, c := range m.Comments {
		var oc Option

		oc, err = f.Make(c)

		if err == nil {
			ocs = append(ocs, oc)
		} else {
			em.Add(err)
		}
	}

	if em.Len() > 0 {
		err = em
	}

	return
}
