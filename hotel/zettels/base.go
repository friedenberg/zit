package zettels

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
	"strings"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/alfa/logz"
	"github.com/friedenberg/zit/bravo/id"
	"github.com/friedenberg/zit/bravo/sha"
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/echo/sharded_store"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	"github.com/friedenberg/zit/golf/stored_zettel_formats"
)

func (zs zettels) storeBaseZettel(z stored_zettel.Stored) (shaz sha.Sha, err error) {
	sb := &strings.Builder{}
	enc := base64.NewEncoder(base64.StdEncoding, sb)
	hash := sha256.New()

	w := io.MultiWriter(enc, hash)

	f := stored_zettel_formats.Objekte{}

	if _, err = f.WriteTo(z, w); err != nil {
		err = errors.Errorf("%s: %s", zs.basePath, err)
		return
	}

	shaz = sha.FromHash(hash)

	var s sharded_store.Shard

	if s, err = zs.store.Shard(shaz.Head()); err != nil {
		err = errors.Error(err)
		return
	}

	if err = enc.Close(); err != nil {
		err = errors.Error(err)
		return
	}

	s.Set(shaz.String(), sb.String())

	for _, e := range z.Zettel.Etiketten {
		zs.etiketten.Add(e)
	}

	return
}

func (zs zettels) update(zettel stored_zettel.Stored) (err error) {
	sb := &strings.Builder{}
	w := base64.NewEncoder(base64.StdEncoding, sb)

	f := stored_zettel_formats.Objekte{}

	if _, err = f.WriteTo(zettel, w); err != nil {
		err = errors.Errorf("%s: %s", zs.basePath, err)
		return
	}

	var s sharded_store.Shard

	if s, err = zs.store.Shard(zettel.Sha.Head()); err != nil {
		err = errors.Error(err)
		return
	}

	if err = w.Close(); err != nil {
		err = errors.Error(err)
		return
	}

	s.Set(zettel.Sha.String(), sb.String())

	return
}

func (zs zettels) updateMutterIfNecessary(mutter, kinder sha.Sha) (err error) {
	if mutter.Equals(kinder) {
		err = errors.Errorf("updating mutter and kinder to same sha: %s", mutter)
		return
	}

	logz.Printf("setting mutter '%s' to kinder '%s'", mutter, kinder)
	if mutter.IsNull() {
		return
	}

	var named stored_zettel.Named

	if named, err = zs.Read(mutter); err != nil {
		err = errors.Error(err)
		return
	}

	named.Kinder = kinder

	if err = zs.update(named.Stored); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (zs zettels) readStoredZettel(sha sha.Sha) (stored stored_zettel.Stored, err error) {
	stored.Sha = sha

	var s sharded_store.Shard

	if s, err = zs.store.Shard(stored.Sha.Head()); err != nil {
		err = errors.Error(err)
		return
	}

	var cr string
	var ok bool

	if cr, ok = s.Read(stored.Sha.String()); !ok {
		err = ErrShaNotFound{Sha: stored.Sha}
		return
	}

	sr := strings.NewReader(cr)
	dec := base64.NewDecoder(base64.StdEncoding, sr)

	f := stored_zettel_formats.Objekte{}

	if _, err = f.ReadFrom(&stored, dec); err != nil {
		err = errors.Errorf("%s: %s", zs.basePath, err)
		return
	}

	return
}

func (zs zettels) readNamedZettel(id id.Id) (named stored_zettel.Named, err error) {
	var ok bool

	if named.Sha, ok = id.(sha.Sha); ok {
		named.Hinweis, err = zs.hinweisen.ReadSha(named.Sha)

		if zs.Konfig().AllowMissingHinweis {
			err = nil
		}

		if err != nil {
			err = errors.Error(err)
			return
		}
	} else {
		if named.Hinweis, ok = id.(hinweis.Hinweis); !ok {
			err = errors.Errorf("unsupported id: '%s'", id)
			return
		}

		if named.Sha, err = zs.hinweisen.Read(named.Hinweis); err != nil {
			err = errors.Error(err)
			return
		}
	}

	if named.Sha.String() == "" {
		panic("empty sha")
	}

	if named.Stored, err = zs.readStoredZettel(named.Sha); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (zs zettels) handleFormatError() {
}
