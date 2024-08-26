package store_browser

import (
	"fmt"
)

type browserItemId struct {
	BrowserId browserId `json:"browser-id"`
	Id        string    `json:"id"`
	Type      string    `json:"type"`
}

func (bi browserItemId) String() string {
	return fmt.Sprintf("/%s/%s-%s", bi.BrowserId, bi.Type, bi.Id)
}
