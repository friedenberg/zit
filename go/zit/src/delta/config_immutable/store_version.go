package config_immutable

import (
	"io"
	"os"
	"strconv"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/values"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
)

var CurrentStoreVersion = StoreVersion(values.Int(9))

func MakeStoreVersion(sv interfaces.StoreVersion) StoreVersion {
	return StoreVersion(values.Int(sv.GetInt()))
}

type StoreVersion values.Int

func (a StoreVersion) Less(b interfaces.StoreVersion) bool {
	return a.String() < b.String()
}

func (a StoreVersion) String() string {
	return values.Int(a).String()
}

func (a StoreVersion) GetInt() int {
	return values.Int(a).Int()
}

func (v *StoreVersion) Set(p string) (err error) {
	var i uint64

	if i, err = strconv.ParseUint(p, 10, 16); err != nil {
		err = errors.Wrap(err)
		return
	}

	*v = StoreVersion(i)

	if CurrentStoreVersion.Less(v) {
		err = errors.Wrap(ErrFutureStoreVersion{StoreVersion: v})
		return
	}

	return
}

func (v *StoreVersion) ReadFromFile(
	p string,
) (err error) {
	if err = v.ReadFromFileOrVersion(p, CurrentStoreVersion); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (v *StoreVersion) ReadFromFileOrVersion(
	p string,
	alternative StoreVersion,
) (err error) {
	var b []byte

	var f *os.File

	if f, err = files.Open(p); err != nil {
		if errors.IsNotExist(err) {
			*v = alternative
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	if b, err = io.ReadAll(f); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = v.Set(string(b)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
