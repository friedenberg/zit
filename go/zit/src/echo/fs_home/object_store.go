package fs_home

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/id"
	"code.linenisgreat.com/zit/go/zit/src/delta/age"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
)

type ObjectStore struct {
	basePath         string
	age              *age.Age
	immutable_config immutable_config.Config
	interfaces.DirectoryPaths
	MoverFactory
}

func (s ObjectStore) objectReader(
	g interfaces.GenreGetter,
	sh sha.ShaLike,
) (rc sha.ReadCloser, err error) {
	var p string

	if p, err = s.DirObjectGenre(
		g,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	o := FileReadOptions{
		Age:             s.age,
		Path:            id.Path(sh.GetShaLike(), p),
		CompressionType: s.immutable_config.CompressionType,
	}

	if rc, err = NewFileReader(o); err != nil {
		err = errors.Wrapf(err, "Gattung: %s", g.GetGenre())
		err = errors.Wrapf(err, "Sha: %s", sh.GetShaLike())
		return
	}

	return
}

func (s ObjectStore) objectWriter(
	g interfaces.GenreGetter,
) (wc sha.WriteCloser, err error) {
	var p string

	if p, err = s.DirObjectGenre(
		g,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	o := MoveOptions{
		Age:                      s.age,
		FinalPath:                p,
		GenerateFinalPathFromSha: true,
		LockFile:                 s.immutable_config.LockInternalFiles,
		CompressionType:          s.immutable_config.CompressionType,
	}

	if wc, err = s.NewMover(o); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// func (s Home) ImportBlobIfNecessary(
//   blobSha interfaces.ShaGetter,
// ) (err error) {
// 	if s.HasBlob(blobSha) {
// 		return
// 	}

// 	p := id.Path(blobSha, c.Blobs)

// 	o := fs_home.FileReadOptions{
// 		Age:             ag,
// 		Path:            p,
// 		CompressionType: c.CompressionType,
// 	}

// 	var rc sha.ReadCloser

// 	if rc, err = fs_home.NewFileReader(o); err != nil {
// 		if errors.IsNotExist(err) {
// 			co.SetError(errors.Errorf("blob missing: %q", p))
// 			err = coErrPrinter(co)
// 		} else {
// 			err = errors.Wrapf(err, "Path: %q", p)
// 		}

// 		return
// 	}

// 	defer errors.DeferredCloser(&err, rc)

// 	var aw sha.WriteCloser

// 	if aw, err = u.GetFSHome().BlobWriter(); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	defer errors.DeferredCloser(&err, aw)

// 	var n int64

// 	if n, err = io.Copy(aw, rc); err != nil {
// 		co.SetError(errors.New("blob copy failed"))
// 		err = coErrPrinter(co)
// 		return
// 	}

// 	shaRc := rc.GetShaLike()

// 	if !shaRc.EqualsSha(blobSha) {
// 		co.SetError(errors.New("blob sha mismatch"))
// 		err = coErrPrinter(co)
// 		ui.TodoRecoverable(
// 			"sku blob mismatch: sku had %s while blob store had %s",
// 			co.Internal.GetBlobSha(),
// 			shaRc,
// 		)
// 	}

// 	// TODO switch to Err and fix test
// 	ui.Out().Printf("copied Blob %s (%d bytes)", blobSha, n)

// 	return
// }
