package query_spec

import (
	"io"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
)

const (
	OpOr            = ','
	OpAnd           = ' '
	OpGroupOpen     = '['
	OpGroupClose    = ']'
	OpNegation      = '^'
	OpExact         = '='
	OpNewline       = '\n'
	OpSigilLatest   = ':'
	OpSigilHistory  = '+'
	OpSigilExternal = '.'
	OpSigilHidden   = '?'
)

var mapOperators = map[rune]bool{
	OpOr:            true,
	OpAnd:           true,
	OpGroupOpen:     true,
	OpGroupClose:    true,
	OpNegation:      true,
	OpExact:         true,
	OpNewline:       true,
	OpSigilLatest:   true,
	OpSigilHistory:  true,
	OpSigilExternal: true,
	OpSigilHidden:   true,
}

// TODO make private
func IsOperator(r rune) (ok bool) {
	return isOperator(r, false)
}

func isOperator(r rune, dotAllowed bool) (ok bool) {
	if dotAllowed && r == '.' {
		return
	}

	_, ok = mapOperators[r]
	return
}

// TODO remove
func GetTokensFromReader(r io.RuneScanner) (out []string, err error) {
	var scanner TokenScanner
	scanner.Reset(r)

	for scanner.Scan() {
		out = append(out, scanner.GetToken().String())
	}

	if err = scanner.Error(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// TODO remove
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
