package store_browser

import "fmt"

type browserId struct {
	Browser string `json:"browser"`
	Pid     string `json:"pid"`
}

// TODO actually populate this correctly
func (bi browserId) String() string {
	if bi.Pid == "" {
		return "browser"
	} else {
		return fmt.Sprintf("browser/%s", bi.Pid)
	}
}
