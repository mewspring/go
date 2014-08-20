// The implementation of this package is heavily inspired by Rob Pike's amazing
// talk titled "Lexical Scanning in Go" [1].
//
// [1]: https://www.youtube.com/watch?v=HxaD_trXwRE

// Package lexer implements lexical tokenization of Go source code.
package lexer

import (
	"fmt"
	"log"
	"strings"
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
	// Index to the first token of the current line.
	line int
}

// lex lexes the input by repeatedly executing the active state function until
// it returns a nil state.
func (l *lexer) lex() {
	// lexToken is the initial state function of the lexer.
	for state := lexToken; state != nil; {
		state = state(l)
	}
}

// errorf emits an error token and terminates the scan by returning a nil state
// function.
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	tok := token.Token{
		Kind: token.Error,
		Val:  fmt.Sprintf(format, args...),
	}
	l.tokens = append(l.tokens, tok)
	return nil
}

// emit emits a token of the specified token type and advances the token start
// position.
func (l *lexer) emit(kind token.Kind) {
	if kind == token.EOF {
		if l.pos < len(l.input) {
			log.Fatalf("lexer.lexer.emit: unexpected eof; pos %d < len(input) %d.\n", l.pos, len(l.input))
		}
		if l.start != l.pos {
			log.Fatalf("lexer.lexer.emit: invalid eof; pending input %q not handled.\n", l.input[l.start:])
		}
	}
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

// backup backs up one rune in the input. It can only be called once per call to
// next.
func (l *lexer) backup() {
	if l.width == 0 {
		log.Fatalln("lexer.lexer.backup: invalid width; no matching call to next.")
	}
	l.pos -= l.width
	l.width = 0
}

// accept consumes the next rune if it's from the valid set. It returns true if
// a rune was consumed and false otherwise.
func (l *lexer) accept(valid string) bool {
	r := l.next()
	if r == eof {
		return false
	}
	if strings.IndexRune(valid, r) == -1 {
		l.backup()
		return false
	}
	return true
}

// acceptRun consumes a run of runes from the valid set. It returns true if a
// rune was consumed and false otherwise.
func (l *lexer) acceptRun(valid string) bool {
	consumed := false
	for l.accept(valid) {
		consumed = true
	}
	return consumed
}

// ignore ignores any pending input read since the last token.
func (l *lexer) ignore() {
	l.start = l.pos
}

// ignoreRun ignores a run of valid runes.
func (l *lexer) ignoreRun(valid string) {
	if l.acceptRun(valid) {
		l.ignore()
	}
}
