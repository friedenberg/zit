package id

import "testing"

func TestHeadTailFromFileName(t *testing.T) {
	input := "/the/file/head/tail.png"
	head, tail := HeadTailFromFileName(input)

	if head != "head" {
		t.Fatalf("expected head, but got %s", head)
	}

	if tail != "tail" {
		t.Fatalf("expected tail, but got %s", tail)
	}
}
