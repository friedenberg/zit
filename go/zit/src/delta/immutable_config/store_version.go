package immutable_config

import (
	"io"
	"os"
	"strconv"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/values"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
)

func ReadStoreVersionFromFile(p string) (v storeVersion, err error) {
	var b []byte

	var f *os.File

	if f, err = files.Open(p); err != nil {
		if errors.IsNotExist(err) {
			v = storeVersion(6)
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

	var i uint64

	if i, err = strconv.ParseUint(string(b), 10, 16); err != nil {
		err = errors.Wrap(err)
		return
	}

	v = storeVersion(i)

	return
}

type storeVersion values.Int

func (a storeVersion) Less(b interfaces.StoreVersion) bool {
	return a.String() < b.String()
}

func (a storeVersion) String() string {
	return values.Int(a).String()
}

func (a storeVersion) GetInt() int {
	return values.Int(a).Int()
}
