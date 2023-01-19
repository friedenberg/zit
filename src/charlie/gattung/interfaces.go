package gattung

import (
	"github.com/friedenberg/zit/src/schnittstellen"
)

type Keyer[T any, T1 schnittstellen.Ptr[T]] interface {
	Key(T1) string
}

//   _____                               _           _
//  |_   _| __ __ _ _ __  ___  __ _  ___| |_ ___  __| |
//    | || '__/ _` | '_ \/ __|/ _` |/ __| __/ _ \/ _` |
//    | || | | (_| | | | \__ \ (_| | (__| ||  __/ (_| |
//    |_||_|  \__,_|_| |_|___/\__,_|\___|\__\___|\__,_|
//
