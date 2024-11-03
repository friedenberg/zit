package builtin_types

import "code.linenisgreat.com/zit/go/zit/src/echo/ids"

const (
	TagTypeTomlV0        = "toml-tag-v0"
	TagTypeTomlV1        = "toml-tag-v1"
	TagTypeLuaV1         = "lua-tag-v1"
	TagTypeLuaV2         = "lua-tag-v2"
	TagTypeLatestDefault = TagTypeTomlV1

	TypeTypeTomlV0        = "toml-type-v0"
	TypeTypeTomlV1        = "toml-type-v1"
	TypeTypeLatestDefault = TypeTypeTomlV1

	ConfigTypeTomlV0        = "toml-config-v0"
	ConfigTypeTomlV1        = "toml-config-v1"
	ConfigTypeLatestDefault = ConfigTypeTomlV1
)

var (
	allSlice []string
	allMap   map[string]struct{}
)

func init() {
	allMap = map[string]struct{}{
		TagTypeTomlV0:    {},
		TagTypeTomlV1:    {},
		TagTypeLuaV1:     {},
		TagTypeLuaV2:     {},
		TypeTypeTomlV0:   {},
		TypeTypeTomlV1:   {},
		ConfigTypeTomlV0: {},
		ConfigTypeTomlV1: {},
	}

	for k := range allMap {
		allSlice = append(allSlice, k)
	}
}

func IsBuiltin(tipe ids.Type) bool {
	_, ok := allMap[tipe.String()]
	return ok
}
