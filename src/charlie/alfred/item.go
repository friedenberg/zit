package alfred

import (
	"encoding/json"
	"strings"
)

type Item struct {
	Title        string   `json:"title,omitempty"`
	Arg          string   `json:"arg,omitempty"`
	Subtitle     string   `json:"subtitle,omitempty"`
	Match        string   `json:"match,omitempty"`
	Icon         ItemIcon `json:"icon,omitempty"`
	Uid          string   `json:"uid,omitempty"`
	ItemType     string   `json:"type,omitempty"`
	QuicklookUrl string   `json:"quicklookurl,omitempty"`
	Text         ItemText `json:"text,omitempty"`
	// Valid        bool     `json:"valid,omitempty"`
}

type ItemText struct {
	Copy string `json:"copy,omitempty"`
}

type ItemIcon struct {
	Type string `json:"type,omitempty"`
	Path string `json:"path,omitempty"`
}

func GenerateItemsJson(i []Item) (j string, err error) {
	sb := &strings.Builder{}

	for idx, v := range i {
		if idx > 0 {
			sb.WriteString(",")
		}

		var alfredItemJson []byte
		alfredItemJson, err = json.Marshal(v)

		if err != nil {
			return
		}

		sb.WriteString(string(alfredItemJson))
	}

	j = sb.String()
	return
}
