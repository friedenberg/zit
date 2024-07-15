package sku_fmt

import (
	"io"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/fs_home"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type Json struct {
	BlobString  string   `json:"blob-string"`
	BlobSha     string   `json:"blob-sha"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	ObjectId    string   `json:"object-id"`
	Sha         string   `json:"sha"`
	Type        string   `json:"type"`
	Tai         string   `json:"tai"`
}

func (j *Json) FromStringAndMetadatei(
	k string,
	m *object_metadata.Metadata,
	s fs_home.Home,
) (err error) {
	var r sha.ReadCloser

	if r, err = s.BlobReader(&m.Blob); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, r)

	var out strings.Builder

	if _, err = io.Copy(&out, r); err != nil {
		err = errors.Wrap(err)
		return
	}

	j.BlobString = out.String()
	j.BlobSha = m.Blob.String()
	j.Description = m.Description.String()
	j.Tags = iter.Strings(m.GetTags())
	j.ObjectId = k
	j.Sha = m.SelfMetadataWithoutTai.String()
	j.Tai = m.Tai.String()
	j.Type = m.Type.String()

	return
}

func (j *Json) FromTransacted(
	sk *sku.Transacted,
	s fs_home.Home,
) (err error) {
	return j.FromStringAndMetadatei(sk.ObjectId.String(), sk.GetMetadata(), s)
}

func (j *Json) ToTransacted(sk *sku.Transacted, s fs_home.Home) (err error) {
	var w sha.WriteCloser

	if w, err = s.BlobWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, w)

	if _, err = io.Copy(w, strings.NewReader(j.BlobString)); err != nil {
		err = errors.Wrap(err)
		return
	}

	// TODO-P1 support states of akte vs akte sha
	sk.SetBlobSha(w.GetShaLike())

	// if err = sk.Metadatei.Tai.Set(j.Tai); err != nil {
	// 	err = errors.Wrap(err)
	// 	return
	// }

	if err = sk.ObjectId.Set(j.ObjectId); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = sk.Metadata.Type.Set(j.Type); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = sk.Metadata.Description.Set(j.Description); err != nil {
		err = errors.Wrap(err)
		return
	}

	var es ids.TagSet

	if es, err = ids.MakeTagSetStrings(j.Tags...); err != nil {
		err = errors.Wrap(err)
		return
	}

	sk.Metadata.SetTags(es)
	sk.Metadata.GenerateExpandedTags()

	return
}
