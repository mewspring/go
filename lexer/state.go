package lexer

import "log"

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
	switch l.pos - l.start {
	case 0:
		// End of file has been reached.
	case 1:
		// A newline character ('\n') has been consumed.
	default:
		log.Fatalf("lexer.lexAutoSemicolon: expected eof or newline; pending input %q not handled.\n", l.input[l.start:])
	}

	// Ignore blank lines.

	panic("not yet implemented.")
}
