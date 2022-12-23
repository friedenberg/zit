package objekte

// type flusher[T gattung.Element, T1 gattung.ElementPtr[T]] struct {
// 	af               gattung.AkteIOFactory
// 	fwc              FuncWriteCloser
// 	objekteFormatter Formatter2
// 	akteFormatter    gattung.Formatter[T, T1]
// }

// func MakeFlusher[T gattung.Element, T1 gattung.ElementPtr[T]](
// 	af gattung.AkteIOFactory,
// 	akteFormatter gattung.Formatter[T, T1],
// 	fwc FuncWriteCloser,
// ) *flusher[T, T1] {
// 	return &flusher[T, T1]{
// 		af:               af,
// 		fwc:              fwc,
// 		objekteFormatter: *MakeFormatter2(),
// 		akteFormatter:    akteFormatter,
// 	}
// }

// func (h *flusher[T, T1]) Flush(
// 	to gattung.StoredPtr,
// 	a *T,
// ) (err error) {
// 	if h.akteFormatter != nil {
// 		var w sha.WriteCloser

// 		if w, err = h.af.AkteWriter(); err != nil {
// 			err = errors.Wrap(err)
// 			return
// 		}

// 		defer errors.Deferred(&err, w.Close)

// 		if _, err = h.akteFormatter.WriteFormat(w, a); err != nil {
// 			err = errors.Wrap(err)
// 			return
// 		}
// 	}

// 	{
// 		var w sha.WriteCLoser

// 		if w, err = h.fwc(to.ObjekteSha()); err != nil {
// 			err = errors.Wrap(err)
// 			return
// 		}

// 		defer errors.Deferred(&err, r.Close)

// 		if _, err = h.objekteFormatter.ReadFormat(r, to); err != nil {
// 			err = errors.Wrap(err)
// 			return
// 		}
// 	}

// 	return
// }
