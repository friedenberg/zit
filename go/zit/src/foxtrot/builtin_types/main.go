package builtin_types

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

const (
	TagTypeTomlV0 = "!toml-tag-v0"
	TagTypeTomlV1 = "!toml-tag-v1"
	TagTypeLuaV1  = "!lua-tag-v1"
	TagTypeLuaV2  = "!lua-tag-v2"

	TypeTypeTomlV0 = "!toml-type-v0"
	TypeTypeTomlV1 = "!toml-type-v1"

	ConfigTypeTomlV0 = "!toml-config-v0"
	ConfigTypeTomlV1 = "!toml-config-v1"

	InventoryListTypeV0 = "!inventory_list-v0"
	InventoryListTypeV1 = "!inventory_list-v1"

	RepoTypeLocalRelay = "!toml-repo-relay-local-v0"

	ZettelIdListTypeV0 = "!zettel_id_list-v0"
)

var (
	allSlice []BuiltinType
	allMap   map[ids.Type]BuiltinType
	defaults map[genres.Genre]BuiltinType
)

type BuiltinType struct {
	ids.Type
	genres.Genre
	Default bool
}

func init() {
	allMap = make(map[ids.Type]BuiltinType)
	defaults = make(map[genres.Genre]BuiltinType)

	register(TagTypeTomlV0, genres.Tag, false)
	register(TagTypeTomlV1, genres.Tag, true)
	register(TagTypeLuaV1, genres.Tag, false)
	register(TagTypeLuaV2, genres.Tag, false)

	register(TypeTypeTomlV0, genres.Type, false)
	register(TypeTypeTomlV1, genres.Type, true)

	register(ConfigTypeTomlV0, genres.Config, false)
	register(ConfigTypeTomlV1, genres.Config, true)

	register(InventoryListTypeV0, genres.InventoryList, false)
	register(InventoryListTypeV1, genres.InventoryList, true)

	register(ZettelIdListTypeV0, genres.None, false)

	register(RepoTypeLocalRelay, genres.Repo, false)
}

func register(tipeString string, g genres.Genre, isDefault bool) {
	registerBuiltinType(
		BuiltinType{
			Type:    ids.MustType(tipeString),
			Genre:   g,
			Default: isDefault,
		},
	)
}

func registerBuiltinType(bt BuiltinType) {
	if _, exists := allMap[bt.Type]; exists {
		panic(fmt.Sprintf("builtin type registered more than once: %s", bt.Type))
	}

	if _, exists := defaults[bt.Genre]; exists && bt.Default {
		panic(fmt.Sprintf("builtin default type registered more than once: %s", bt.Type))
	}

	allMap[bt.Type] = bt
	allSlice = append(allSlice, bt)

	if bt.Default {
		defaults[bt.Genre] = bt
	}
}

func IsBuiltin(tipe ids.Type) bool {
	_, ok := allMap[tipe]
	return ok
}

func Get(t ids.Type) (BuiltinType, bool) {
	bt, ok := allMap[t]
	return bt, ok
}

func GetOrPanic(idString string) BuiltinType {
	t := ids.MustType(idString)
	bt, ok := Get(t)

	if !ok {
		panic(fmt.Sprintf("no builtin type found for %q", t))
	}

	return bt
}

func Default(genre genres.Genre) (ids.Type, bool) {
	bt, ok := defaults[genre]
	return bt.Type, ok
}

func DefaultOrPanic(genre genres.Genre) ids.Type {
	t, ok := Default(genre)

	if !ok {
		panic(fmt.Sprintf("default missing for genre %q", genre))
	}

	return t
}
