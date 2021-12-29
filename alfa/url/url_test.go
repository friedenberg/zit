package url

import "testing"

func assertParseUrlString(t *testing.T, in, out string) {
	t.Helper()
	u, err := ParseURL(in)

	if err != nil {
		t.Errorf("failed to parse url: %s", err)
	}

	actual := u.String()

	if actual != out {
		t.Errorf("expected hostname '%s', but got '%s'", out, actual)
	}
}

func TestParseURLString(t *testing.T) {
	assertParseUrlString(t, "https://www.google.com", "https://www.google.com")
}

func TestParseURLString2(t *testing.T) {
	assertParseUrlString(t, "http://www.google.com", "https://www.google.com")
}
