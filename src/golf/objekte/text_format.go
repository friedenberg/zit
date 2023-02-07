package objekte

// import (
// 	"io"

// 	"github.com/friedenberg/zit/src/alfa/errors"
// 	"github.com/friedenberg/zit/src/alfa/schnittstellen"
// 	"github.com/friedenberg/zit/src/alfa/toml"
// 	"github.com/friedenberg/zit/src/bravo/sha"
// )

// type TextFormat[
// 	T schnittstellen.Objekte[T],
// 	T1 schnittstellen.ObjektePtr[T],
// ] struct {
// 	arf              schnittstellen.AkteIOFactory
// 	ignoreTomlErrors bool
// }

// func MakeFormatText[
// 	T schnittstellen.Objekte[T],
// 	T1 schnittstellen.ObjektePtr[T],
// ](arf schnittstellen.AkteIOFactory) *TextFormat[T, T1] {
// 	return &TextFormat[T, T1]{
// 		arf: arf,
// 	}
// }

// // TODO-P4 remove
// func MakeFormatTextIgnoreTomlErrors[
// 	T schnittstellen.Objekte[T],
// 	T1 schnittstellen.ObjektePtr[T],
// ](arf schnittstellen.AkteIOFactory) *TextFormat[T, T1] {
// 	return &TextFormat[T, T1]{
// 		arf:              arf,
// 		ignoreTomlErrors: true,
// 	}
// }

// func (f TextFormat[T, T1]) Parse(r io.Reader, t T1) (n int64, err error) {
// 	return f.ReadFormat(r, t)
// }

// func (f TextFormat[T, T1]) ReadFormat(r io.Reader, t T1) (n int64, err error) {
// 	var aw sha.WriteCloser

// 	if aw, err = f.arf.AkteWriter(); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	defer errors.Deferred(&err, aw.Close)

// 	pr, pw := io.Pipe()
// 	td := toml.NewDecoder(pr)

// 	chDone := make(chan error)

// 	go func(pr *io.PipeReader) {
// 		var err error
// 		defer func() {
// 			chDone <- err
// 			close(chDone)
// 		}()

// 		defer func() {
// 			if r := recover(); r != nil {
// 				if f.ignoreTomlErrors {
// 					err = nil
// 				} else {
// 					err = toml.MakeError(errors.Errorf("panicked during toml decoding: %s", r))
// 					pr.CloseWithError(errors.Wrap(err))
// 				}
// 			}
// 		}()

// 		if err = td.Decode(&t.Akte); err != nil {
// 			switch {
// 			case !errors.IsEOF(err) && !f.ignoreTomlErrors:
// 				err = errors.Wrap(toml.MakeError(err))
// 				pr.CloseWithError(err)

// 			case !errors.IsEOF(err) && f.ignoreTomlErrors:
// 				err = nil
// 			}
// 		}
// 	}(pr)

// 	mw := io.MultiWriter(aw, pw)

// 	if n, err = io.Copy(mw, r); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	if err = pw.Close(); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	if err = <-chDone; err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	t.Sha = sha.Make(aw.Sha())

// 	return
// }

// func (f TextFormat[T, T1]) Format(w io.Writer, t T1) (n int64, err error) {
// 	return f.WriteFormat(w, t)
// }

// func (f TextFormat[T, T1]) WriteFormat(w io.Writer, t T1) (n int64, err error) {
// 	var ar sha.ReadCloser

// 	if ar, err = f.arf.AkteReader(t.Sha); err != nil {
// 		if errors.IsNotExist(err) {
// 			err = nil
// 		} else {
// 			err = errors.Wrap(err)
// 		}

// 		return
// 	}

// 	defer errors.Deferred(&err, ar.Close)

// 	if n, err = io.Copy(w, ar); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	return
// }
