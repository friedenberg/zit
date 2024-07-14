package organize_text

import (
	"bufio"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/ohio"
	"code.linenisgreat.com/zit/go/zit/src/echo/format"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type Metadatei struct {
	// metadatei.Metadatei
	ids.TagSet
	Matchers interfaces.SetLike[sku.Query]
	Comments []string
	Typ      ids.Type
}

func (m Metadatei) RemoveFromTransacted(sk *sku.Transacted) (err error) {
	mes := sk.Metadatei.GetTags().CloneMutableSetPtrLike()

	if err = m.Each(mes.Del); err != nil {
		err = errors.Wrap(err)
		return
	}

	sk.Metadatei.SetTags(mes)

	return
}

func (m Metadatei) AsMetadatei() (m1 object_metadata.Metadata) {
	m1.Type = m.Typ
	m1.SetTags(m.TagSet)
	return
}

func (m Metadatei) GetMetadataWriterTo() object_metadata.MetadataWriterTo {
	return m
}

func (m Metadatei) HasMetadataContent() bool {
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

	mes := ids.MakeTagMutableSet()

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

	m.TagSet = mes.CloneSetPtrLike()

	return
}

func (m Metadatei) WriteTo(w1 io.Writer) (n int64, err error) {
	w := format.NewLineWriter()

	for _, e := range iter.SortedStrings(m.TagSet) {
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
