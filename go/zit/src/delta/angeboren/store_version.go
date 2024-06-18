package angeboren

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/values"
)

type storeVersion values.Int

func (a storeVersion) Less(b schnittstellen.StoreVersion) bool {
	return a.String() < b.String()
}

func (a storeVersion) String() string {
	return values.Int(a).String()
}

func (a storeVersion) GetInt() int {
	return values.Int(a).Int()
}
