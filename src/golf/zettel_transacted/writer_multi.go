package zettel_transacted

import "github.com/friedenberg/zit/src/alfa/errors"

type multiWriter struct {
	writers []Writer
	pool    *Pool
}

func MakeWriterMulti(pool *Pool, ws ...Writer) Writer {
	return &multiWriter{
		pool:    pool,
		writers: ws,
	}
}

func (w multiWriter) WriteZettelTransacted(z *Zettel) (err error) {
	if w.pool != nil {
		defer w.pool.Put(z)
	}

	for _, w := range w.writers {
		if err = w.WriteZettelTransacted(z); err != nil {
			if errors.IsEOF(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}
	}

	return
}
