// The implementation of this package is heavily inspired by Rob Pike's amazing
// talk titled "Lexical Scanning in Go" [1].
//
// [1]: https://www.youtube.com/watch?v=HxaD_trXwRE

// Package lexer implements lexical tokenization of Go source code.
package lexer

import (
	"log"
	"unicode/utf8"

	"github.com/mewlang/go/token"
)

// Parse lexes the input string into a slice of tokens. While breaking the input
// into tokens, the next token is the longest sequence of characters that form a
// valid token.
func Parse(input string) (tokens []token.Token) {
	l := &lexer{
		input: input,
		// TODO(u): Fix cap; estimate the average token size by lexing the source
		// code of the standard library.
		tokens: make([]token.Token, 0, len(input)/3),
	}

	// Tokenize the input.
	l.lex()

	return l.tokens
}

// A lexer lexes an input string into a slice of tokens.
type lexer struct {
	// The input string.
	input string
	// Start position of the current token.
	start int
	// Current position in the input.
	pos int
	// Width in byte of the last rune read with next.
	width int
	// A slice of scanned tokens.
	tokens []token.Token
}

// lex lexes the input by repeatedly executing the active state function until
// it returns a nil state.
func (l *lexer) lex() {
	// lexToken is the initial state function of the lexer.
	for state := lexToken; state != nil; {
		state = state(l)
	}
}

// emit emits a token of the provided token type and advances the token start
// position.
func (l *lexer) emit(kind token.Kind) {
	tok := token.Token{
		Kind: kind,
		Val:  l.input[l.start:l.pos],
	}
	l.tokens = append(l.tokens, tok)
	l.start = l.pos
}

// eof is the rune returned by next when no more input is available.
const eof = -1

// next consumes and returns the next rune of the input.
func (l *lexer) next() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

// peek returns but does not consume the next rune of the input.
func (l *lexer) peek() (r rune) {
	r = l.next()
	if r == eof {
		return eof
	}
	l.backup()
	return r
}

// backup backs up one rune in the input. It can only be called once per call to
// next.
func (l *lexer) backup() {
	if l.width == 0 {
		// TODO(u): Handle eof elsewhere so we never hit this case.
		log.Fatalln("lexer.lexer.backup: invalid width; no matching call to next.")
	}
	l.pos -= l.width
	l.width = 0
}
