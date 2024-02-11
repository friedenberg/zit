package script_config

import (
	"testing"

	"code.linenisgreat.com/zit/src/alfa/toml"
	"code.linenisgreat.com/zit/src/bravo/test_logz"
)

func TestToml(t1 *testing.T) {
	t := test_logz.T{T: t1}

	strToml := `
description = "wow"
file-extension = "pdf"
uti = "com.adobe.pdf"
script = """
cat
"""
  `

	var sut WithOutputFormat

	{
		err := toml.Unmarshal([]byte(strToml), &sut)
		t.AssertNoError(err)
	}

	if sut.Description != "wow" {
		t.Errorf("expected Description 'wow' but got %q", sut.Description)
	}

	if sut.FileExtension != "pdf" {
		t.Errorf("expected FileExtension 'pdf' but got %q", sut.FileExtension)
	}

	if sut.UTI != "com.adobe.pdf" {
		t.Errorf("expected UTI 'com.adobe.pdf' but got %q", sut.UTI)
	}

	if sut.Script != "cat\n" {
		t.Errorf("expected Script 'cat\\n' but got %q", sut.Script)
	}
}
