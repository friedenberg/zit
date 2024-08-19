package query_spec

/* Token types:
- id: [-!/a-z_.]+
- literal: ".*"

*/

//go:generate stringer -type=TokenType
type TokenType int

const (
	TokenTypeIncomplete = TokenType(iota)
	TokenTypeOperator   // " =,.:+?^[]"
	TokenTypeIdentifier // ["one/uno", "tag-one", "!type", "/browser/bookmark-1"...]
	TokenTypeLiteral    // ["\"some text\"", "\"some text \\\" with escape\""]
	TokenTypeField      // ["field=\"some text\"", "url=\"some text \\\" with escape\""]
)
