// Package lexer implements lexical tokenization of Go source code.
package lexer

import (
	"github.com/mewlang/go/token"
)

// Parse lexes the input string into a slice of tokens. While breaking the input
// into tokens, the next token is the longest sequence of characters that form a
// valid token.
func Parse(input string) (tokens []token.Token) {
	panic("not yet implemented.")
}
