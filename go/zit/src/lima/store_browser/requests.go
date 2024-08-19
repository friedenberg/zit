package store_browser

type requestUrlsPut struct {
	Deleted []string              `json:"deleted"`
	Added   []createOneTabRequest `json:"added"`
}

type createOneTabRequest struct {
	Url string `json:"url"`
	// WindowId string `json:"windowId"`
}
