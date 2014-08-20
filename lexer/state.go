package lexer

import (
	"log"

	"github.com/mewlang/go/token"
)

// A stateFn represents the state of the lexer as a function that returns a
// state function.
type stateFn func(l *lexer) stateFn

// whitespace specifies the white space characters (except newline), which
// include spaces (U+0020), horizontal tabs (U+0009), and carriage returns
// (U+000D).
const whitespace = " \t\r"

// lexToken lexes a token of the Go programming language. It is the initial
// state function of the lexer.
func lexToken(l *lexer) stateFn {
	// Ignore white space characters (except newline).
	l.ignoreRun(whitespace)

	r := l.next()
	switch r {
	case '\n', eof:
		return lexAutoSemicolon
	}

	panic("not yet implemented.")
}

// lexAutoSemicolon inserts a semicolon if the correct conditions have been met.
//
// When the input is broken into tokens, a semicolon is automatically inserted
// into the token stream at the end of a non-blank line if the line's final
// token is
//    * an identifier
//    * an integer, floating-point, imaginary, rune, or string literal
//    * one of the keywords break, continue, fallthrough, or return
//    * one of the operators and delimiters ++, --, ), ], or }
//
// ref: http://golang.org/ref/spec#Semicolons
func lexAutoSemicolon(l *lexer) stateFn {
	atEOF := false
	switch l.pos - l.start {
	case 0:
		// End of file has been reached.
		atEOF = true
	case 1:
		l.ignore()
		// A newline character ('\n') has been consumed.
	default:
		log.Fatalf("lexer.lexAutoSemicolon: expected eof or newline, got %q.\n", l.input[l.start:])
	}

	// When the input is broken into tokens, a semicolon is automatically
	// inserted into the token stream at the end of a non-blank line if the
	// line's final token is
	insert := false
	if len(l.tokens) > l.line {
		last := l.tokens[len(l.tokens)-1]
		switch last.Kind {
		case token.Ident:
			// * an identifier
			insert = true
		case token.Int, token.Float, token.Imag, token.Rune, token.String:
			// * an integer, floating-point, imaginary, rune, or string literal
			insert = true
		case token.Break, token.Continue, token.Fallthrough, token.Return:
			// * one of the keywords break, continue, fallthrough, or return
			insert = true
		case token.Inc, token.Dec, token.Rparen, token.Rbrack, token.Rbrace:
			// * one of the operators and delimiters ++, --, ), ], or }
			insert = true
		}

		// Insert semicolon.
		if insert {
			tok := token.Token{
				Kind: token.Semicolon,
				Val:  ";",
			}
			l.tokens = append(l.tokens, tok)
		}
	}

	if atEOF {
		// Emit an EOF and terminate the lexer with a nil state function.
		l.emit(token.EOF)
		return nil
	}

	// Update the index to the first token of the current line.
	l.line = len(l.tokens)

	return lexToken
}
