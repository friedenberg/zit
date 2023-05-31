package store_util

import (
	"github.com/friedenberg/zit/src/alfa/errors"
)

func (c *common) Flush() (err error) {
	if err = c.typenIndex.Flush(c); err != nil {
		err = errors.Wrapf(err, "failed to flush typen index")
		return
	}

	if err = c.kennungIndex.Flush(); err != nil {
		err = errors.Wrapf(err, "failed to flush kennung index")
		return
	}

	if err = c.Abbr.Flush(); err != nil {
		err = errors.Wrapf(err, "failed to flush abbr index")
		return
	}

	return
}
