package builtin_types

import "code.linenisgreat.com/zit/go/zit/src/echo/ids"

const (
	TypeTypeV0     = "toml-type-v0"
	TypeTypeV1     = "toml-type-v1"
	TypeTypeLatest = TypeTypeV1

	ConfigTypeV0     = "toml-config-v0"
	ConfigTypeV1     = "toml-config-v1"
	ConfigTypeLatest = ConfigTypeV1
)

var (
	allSlice []string
	allMap   map[string]struct{}
)

func init() {
	allMap = map[string]struct{}{
		TypeTypeV0:   {},
		TypeTypeV1:   {},
		ConfigTypeV0: {},
		ConfigTypeV1: {},
	}

	for k := range allMap {
		allSlice = append(allSlice, k)
	}
}

func IsBuiltin(tipe ids.Type) bool {
	_, ok := allMap[tipe.String()]
	return ok
}
