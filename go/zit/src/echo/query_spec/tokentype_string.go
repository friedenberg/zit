// Code generated by "stringer -type=TokenType"; DO NOT EDIT.

package query_spec

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[TokenTypeIncomplete-0]
	_ = x[TokenTypeOperator-1]
	_ = x[TokenTypeIdentifier-2]
	_ = x[TokenTypeLiteral-3]
	_ = x[TokenTypeField-4]
}

const _TokenType_name = "TokenTypeIncompleteTokenTypeOperatorTokenTypeIdentifierTokenTypeLiteralTokenTypeField"

var _TokenType_index = [...]uint8{0, 19, 36, 55, 71, 85}

func (i TokenType) String() string {
	if i < 0 || i >= TokenType(len(_TokenType_index)-1) {
		return "TokenType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}