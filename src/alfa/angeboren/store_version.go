package angeboren

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/int_value"
)

type StoreVersion interface {
	schnittstellen.Stringer
	schnittstellen.Lessor[StoreVersion]
}

type storeVersion int_value.IntValue

func (a storeVersion) Less(b StoreVersion) bool {
	return a.String() < b.String()
}

func (a storeVersion) String() string {
	return int_value.IntValue(a).String()
}
