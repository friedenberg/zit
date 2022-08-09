package verzeichnisse

import "github.com/friedenberg/zit/bravo/sha"

type Row struct {
	Sha     sha.Sha `json:"Sha"`
	Key     string  `json:"Key"`
	Type    string  `json:"Type"`
	Value   string  `json:"Value"`
	Objekte []byte  `json:"Objekte"`
}
