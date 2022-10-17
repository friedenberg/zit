package proto_objekte

import (
	"flag"
	"fmt"
)

type ProtoObjekte interface {
	fmt.Stringer
}

type ProtoObjektePointer interface {
	flag.Value
}
