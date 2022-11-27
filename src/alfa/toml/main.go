package toml

import "github.com/pelletier/go-toml/v2"

var (
	Unmarshal  = toml.Unmarshal
	NewDecoder = toml.NewDecoder
	NewEncoder = toml.NewEncoder
)

type (
	StrictMissingError = toml.StrictMissingError
)
