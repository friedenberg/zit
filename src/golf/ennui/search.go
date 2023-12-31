package ennui

import (
	"bytes"
	"fmt"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/charlie/sha"
)

func (e *ennui) seekToFirstBinarySearch(shMet *sha.Sha) (err error) {
	if e.f == nil {
		err = collections.ErrNotFound("fd nil: " + shMet.String())
		return
	}

	var low, mid, hi int64
	shMid := &sha.Sha{}

	var rowCount int64

	if rowCount, err = e.GetRowCount(); err != nil {
		err = errors.Wrap(err)
		return
	}

	hi = rowCount
	loops := 0

	for low <= hi {
		loops++
		mid = (hi + low) / 2

		// var loc int64

		if _, err = e.f.Seek(mid*RowSize, io.SeekStart); err != nil {
			err = errors.Wrap(err)
			return
		}

		if _, err = shMid.ReadFrom(e.f); err != nil {
			err = errors.Wrap(err)
			return
		}

		cmp := bytes.Compare(shMet.GetShaBytes(), shMid.GetShaBytes())

		switch cmp {
		case -1:
			if low == hi-1 {
				low = hi
			} else {
				hi = mid - 1
			}

		case 0:
			// found
			if _, err = e.f.Seek(mid*RowSize, io.SeekStart); err != nil {
				err = errors.Wrap(err)
				return
			}

			return

		case 1:
			low = mid + 1

		default:
			panic("not possible")
		}
	}

	err = collections.ErrNotFound(fmt.Sprintf("%d: %s", loops, shMet.String()))

	return
}

func (e *ennui) seekToFirstLinearSearch(shMet *sha.Sha) (err error) {
	if e.f == nil {
		err = collections.ErrNotFound("fd nil: " + shMet.String())
		return
	}

	var rowCount int64
	shMid := &sha.Sha{}

	if rowCount, err = e.GetRowCount(); err != nil {
		err = errors.Wrap(err)
		return
	}

	for loc := int64(0); loc <= rowCount; loc++ {
		// var loc int64

		if _, err = e.f.Seek(loc*RowSize, io.SeekStart); err != nil {
			err = errors.Wrap(err)
			return
		}

		if _, err = shMid.ReadFrom(e.f); err != nil {
			err = errors.Wrap(err)
			return
		}

		if bytes.Equal(shMet.GetShaBytes(), shMid.GetShaBytes()) {

			if _, err = e.f.Seek(loc*RowSize, io.SeekStart); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}
	}

	err = collections.ErrNotFound(shMet.String())

	return
}
