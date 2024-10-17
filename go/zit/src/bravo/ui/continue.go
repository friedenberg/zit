package ui

import "fmt"

func Continue(header string, err error) (shouldContinue bool) {
	Err().Printf("%s:", header)
	Err().Print(err)
	Err().Printf("ignore and continue? (y/*)")

	var answer rune
	var n int

	if n, err = fmt.Scanf("%c", &answer); err != nil {
		shouldContinue = false
		Err().Printf("failed to read answer: %s", err)
		return
	}

	if n != 1 {
		shouldContinue = false
		Err().Printf("failed to read at exactly 1 answer: %s", err)
		return
	}

	if answer == 'y' || answer == 'Y' {
		shouldContinue = true
	}

	return
}
