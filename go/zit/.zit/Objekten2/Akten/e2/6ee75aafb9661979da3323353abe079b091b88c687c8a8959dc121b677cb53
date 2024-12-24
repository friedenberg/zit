package object_metadata

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/flag_policy"
	"code.linenisgreat.com/zit/go/zit/src/bravo/expansion"
	flag2 "code.linenisgreat.com/zit/go/zit/src/bravo/flags"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/descriptions"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type MetadataWriterTo interface {
	io.WriterTo
	HasMetadataContent() bool
}

type Field = string_format_writer.Field

type Metadata struct {
	// StoreVersion values.Int
	Description descriptions.Description
	Tags        ids.TagMutableSet // public for gob, but should be private
	Type        ids.Type

	Shas
	Tai ids.Tai

	Comments []string
	Cache    Cache
	Fields   []Field
}

func (m *Metadata) GetMetadata() *Metadata {
	return m
}

func (m *Metadata) Sha() *sha.Sha {
	return &m.SelfMetadataObjectIdParent
}

func (m *Metadata) Mutter() *sha.Sha {
	return &m.ParentMetadataObjectIdParent
}

func (m *Metadata) AddToFlagSet(f *flag.FlagSet) {
	f.Var(
		&m.Description,
		"description",
		"the description to use for created or updated Zettels",
	)

	// TODO add support for tag_paths
	fes := flag2.Make(
		flag_policy.FlagPolicyAppend,
		func() string {
			return m.Cache.TagPaths.String()
		},
		func(v string) (err error) {
			vs := strings.Split(v, ",")

			for _, v := range vs {
				if err = m.AddTagString(v); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			return
		},
		func() {
			m.ResetTags()
		},
	)

	f.Var(
		fes,
		"tags",
		"the tags to use for created or updated object",
	)

	f.Func(
		"type",
		"the type for the created or updated object",
		func(v string) (err error) {
			return m.Type.Set(v)
		},
	)
}

func (z *Metadata) UserInputIsEmpty() bool {
	if !z.Description.IsEmpty() {
		return false
	}

	if z.Tags != nil && z.Tags.Len() > 0 {
		return false
	}

	if !ids.IsEmpty(z.Type) {
		return false
	}

	return true
}

func (z *Metadata) IsEmpty() bool {
	if !z.Blob.IsNull() {
		return false
	}

	if !z.UserInputIsEmpty() {
		return false
	}

	if !z.Tai.IsZero() {
		return false
	}

	return true
}

// TODO fix issue with GetTags being nil sometimes
func (m *Metadata) GetTags() ids.TagSet {
	if m.Tags == nil {
		m.Tags = ids.MakeTagMutableSet()
	}

	return m.Tags
}

func (m *Metadata) ResetTags() {
	if m.Tags == nil {
		m.Tags = ids.MakeTagMutableSet()
	}

	m.Tags.Reset()
	m.Cache.TagPaths.Reset()
}

func (z *Metadata) AddTagString(es string) (err error) {
	if es == "" {
		return
	}

	var e ids.Tag

	if err = e.Set(es); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = z.AddTagPtr(&e); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (m *Metadata) AddTagPtr(e *ids.Tag) (err error) {
	if e == nil || e.String() == "" {
		return
	}

	if m.Tags == nil {
		m.Tags = ids.MakeTagMutableSet()
	}

	ids.AddNormalizedTag(m.Tags, e)
	cs := catgut.MakeFromString(e.String())
	m.Cache.TagPaths.AddTag(cs)

	return
}

func (m *Metadata) AddTagPtrFast(e *ids.Tag) (err error) {
	if m.Tags == nil {
		m.Tags = ids.MakeTagMutableSet()
	}

	if err = m.Tags.Add(*e); err != nil {
		err = errors.Wrap(err)
		return
	}

	cs := catgut.MakeFromString(e.String())

	if err = m.Cache.TagPaths.AddTag(cs); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (m *Metadata) SetTags(e ids.TagSet) {
	if m.Tags == nil {
		m.Tags = ids.MakeTagMutableSet()
	}

	m.Tags.Reset()

	if e == nil {
		return
	}

	if e.Len() == 1 && e.Any().String() == "" {
		panic("empty tag set")
	}

	errors.PanicIfError(e.EachPtr(m.AddTagPtr))
}

func (z *Metadata) GetType() ids.Type {
	return z.Type
}

func (z *Metadata) GetTypePtr() *ids.Type {
	return &z.Type
}

func (z *Metadata) GetTai() ids.Tai {
	return z.Tai
}

// TODO-P2 remove
func (b *Metadata) EqualsSansTai(a *Metadata) bool {
	return EqualerSansTai.Equals(a, b)
}

// TODO-P2 remove
func (pz *Metadata) Equals(z1 *Metadata) bool {
	return Equaler.Equals(pz, z1)
}

func (a *Metadata) Subtract(
	b *Metadata,
) {
	if a.Type.String() == b.Type.String() {
		a.Type = ids.Type{}
	}

	if a.Tags == nil {
		return
	}

	// ui.Debug().Print("before", b.Tags, a.Tags)

	for e := range b.Tags.AllPtr() {
		// ui.Debug().Print(e)
		a.Tags.DelPtr(e)
	}

	// ui.Debug().Print("after", b.Tags, a.Tags)
}

func (mp *Metadata) AddComment(f string, vals ...interface{}) {
	mp.Comments = append(mp.Comments, fmt.Sprintf(f, vals...))
}

func (selbst *Metadata) SetMutter(mg Getter) (err error) {
	mutter := mg.GetMetadata()

	if err = selbst.Mutter().SetShaLike(
		mutter.Sha(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = selbst.ParentMetadataObjectIdParent.SetShaLike(
		&mutter.SelfMetadataObjectIdParent,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (m *Metadata) GenerateExpandedTags() {
	m.Cache.SetExpandedTags(ids.ExpandMany(
		m.GetTags(),
		expansion.ExpanderRight,
	))
}
