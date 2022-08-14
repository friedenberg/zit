package url

import (
	"net/url"
)

func ParseURL(u string) (ur *url.URL, err error) {
	ur, err = url.Parse(u)

	if err != nil {
		return
	}

	if ur.Scheme == "" || ur.Scheme == "http" {
		ur.Scheme = "https"
		return ParseURL(ur.String())
	}

	return
}
