package box

//go:generate stringer -type=TokenType
type TokenType int

const (
	TokenTypeIncomplete = TokenType(iota)
	TokenTypeOperator   // " =,.:+?^[]"
	TokenTypeIdentifier // ["one", "uno", "tag", "one", "type", "/browser/bookmark-1", "@sha"...]
	TokenTypeLiteral    // ["\"some text\"", "\"some text \\\" with escape\""]
)

//go:generate stringer -type=SeqType
type SeqType int

const (
	SeqTypeUnknown    = SeqType(iota)
	SeqTypeIdentifier // one/uno, tag, !type, /browser/bookmark-1, @abcd
	SeqTypeField      // url=blah blob=hello contents="wow" contents="wow with\" quote"
)
