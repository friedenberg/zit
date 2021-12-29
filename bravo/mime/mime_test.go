package mime

// import (
// 	"testing"
// )

// func TestExtensionTypeNoPrefix(t *testing.T) {
// 	mt := `markdown`
// 	var m Mime
// 	var err error

// 	if err = m.Set(mt); err != nil {
// 		t.Fatalf("%s", err)
// 	}

// 	if m.String() != Markdown.String() {
// 		t.Fatalf("expected %s, but got %s", Markdown.String(), m.String())
// 	}
// }

// func TestExtensionTypeYesPrefix(t *testing.T) {
// 	mt := `.markdown`
// 	var m Mime
// 	var err error

// 	if err = m.Set(mt); err != nil {
// 		t.Fatalf("%s", err)
// 	}

// 	if m.String() != Markdown.String() {
// 		t.Fatalf("expected %s, but got %s", Markdown.String(), m.String())
// 	}
// }

// func TestExtensionTypeYesPrefixAndFile(t *testing.T) {
// 	mt := `testfile.markdown`
// 	var m Mime
// 	var err error

// 	if err = m.Set(mt); err != nil {
// 		t.Fatalf("%s", err)
// 	}

// 	if m.String() != Markdown.String() {
// 		t.Fatalf("expected %s, but got %s", Markdown.String(), m.String())
// 	}
// }

// func TestPreferredExtension(t *testing.T) {
// 	m := Markdown

// 	expected := `.markdown`

// 	if m.PreferredExtension() != expected {
// 		t.Fatalf("expected %s, but got %s", expected, m.PreferredExtension())
// 	}
// }
