package zittish

import (
	"bufio"
	"io"
	"strings"
	"unicode/utf8"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
)

const (
	OpOr           = ','
	OpAnd          = ' '
	OpGroupOpen    = '['
	OpGroupClose   = ']'
	OpNegation     = '^'
	OpExact        = '='
	OpNewline      = '\n'
	OpSigilSchwanz = ':'
	OpSigilHistory = '+'
	OpSigilCwd     = '.'
	OpSigilHidden  = '?'
)

var mapMatcherOperators = map[rune]bool{
	OpOr:           true,
	OpAnd:          true,
	OpGroupOpen:    true,
	OpGroupClose:   true,
	OpNegation:     true,
	OpExact:        true,
	OpNewline:      true,
	OpSigilSchwanz: true,
	OpSigilHistory: true,
	OpSigilCwd:     true,
	OpSigilHidden:  true,
}

func IsMatcherOperator(r rune) (ok bool) {
	_, ok = mapMatcherOperators[r]
	return
}

func NextToken(
	rr io.RuneScanner,
	token *catgut.String,
) (err error) {
	afterFirst := false

	for {
		var r rune
		r, _, err = rr.ReadRune()

		if err != nil {
			if token.Len() > 0 && err == io.EOF {
				err = nil
			}

			return
		}

		wasSplitRune := IsMatcherOperator(r)

		switch {
		case wasSplitRune && !afterFirst:
			token.WriteRune(r)
			return

		case !wasSplitRune:
			token.WriteRune(r)
			afterFirst = true
			continue

		default:
			if err = rr.UnreadRune(); err != nil {
				err = errors.Wrapf(err, "%c", r)
			}

			return
		}
	}
}

func SplitMatcher(
	data []byte,
	atEOF bool,
) (advance int, token []byte, err error) {
	for width, i := 0, 0; i < len(data); i += width {
		var r rune

		r, width = utf8.DecodeRune(data[i:])

		wasSplitRune := IsMatcherOperator(r)

		switch {
		case !wasSplitRune:
			continue

		case wasSplitRune && i == 0:
			return width, data[:width], nil

		default:
			return i, data[:i], nil
		}
	}

	// If we're at EOF, we have a final, non-empty, non-terminated word.  Return
	// it.
	if atEOF && len(data) > 0 {
		return len(data), data[0:], nil
	}

	return 0, nil, nil
}

func GetTokensFromReader(r io.Reader) (out []string, err error) {
	scanner := bufio.NewScanner(r)

	scanner.Split(SplitMatcher)

	for scanner.Scan() {
		out = append(out, scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func GetTokensFromStrings(vs ...string) (out []string, err error) {
	for i, v := range vs {
		if i > 0 {
			out = append(out, " ")
		}

		var more []string

		if more, err = GetTokensFromReader(strings.NewReader(v)); err != nil {
			err = errors.Wrap(err)
			return
		}

		out = append(out, more...)
	}

	return
}
