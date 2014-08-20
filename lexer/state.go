package lexer

import (
	"strings"

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
	case '/':
		return lexDivOrComment
	case '`':
		return lexRawString
	case '\n':
		l.ignore()
		insertSemicolon(l)
		// Update the index to the first token of the current line.
		l.line = len(l.tokens)
		return lexToken
	case eof:
		insertSemicolon(l)
		// Emit an EOF and terminate the lexer with a nil state function.
		l.emit(token.EOF)
		return nil
	}

	panic("not yet implemented.")
}

// lexDivOrComment lexes a division operator (/), a division assignment operator
// (/=), a line comment (//), or a general comment (/*). A slash character ('/')
// has already been consumed.
func lexDivOrComment(l *lexer) stateFn {
	r := l.next()
	switch r {
	case '=':
		// Division assignment operator (/=)
		l.emit(token.DivAssign)
		return lexToken
	case '/':
		// Line comment (//).
		return lexLineComment
	case '*':
		// General comment (/*)
		return lexGeneralComment
	default:
		// Division operator (/)
		l.backup()
		l.emit(token.Div)
		return lexToken
	}
}

// lexLineComment lexes a line comment. A line comment acts like a newline.
func lexLineComment(l *lexer) stateFn {
	insertSemicolon(l)
	for {
		switch l.next() {
		case eof:
			l.emit(token.Comment)
			// Emit an EOF and terminate the lexer with a nil state function.
			l.emit(token.EOF)
			return nil
		case '\n':
			l.emit(token.Comment)
			// Update the index to the first token of the current line.
			l.line = len(l.tokens)
			return lexToken
		}
	}
}

// lexGeneralComment lexes a general comment. A general comment containing one
// or more newlines acts like a newline, otherwise it acts like a space.
func lexGeneralComment(l *lexer) stateFn {
	hasNewline := false
	for !strings.HasSuffix(l.input[l.start:l.pos], "*/") {
		switch l.next() {
		case eof:
			return l.errorf("unexpected eof in general comment")
		case '\n':
			hasNewline = true
		}
	}
	if hasNewline {
		insertSemicolon(l)
		// Update the index to the first token of the current line.
		l.line = len(l.tokens)
	}

	l.emit(token.Comment)

	return lexToken
}

// lexRawString lexes a raw string literal (`foo`). A back quote character ('`')
// has already been consumed.
func lexRawString(l *lexer) stateFn {
	for {
		switch l.next() {
		case eof:
			return l.errorf("unexpected eof in raw string literal")
		case '`':
			l.emit(token.String)
			return lexToken
		}
	}
}

// TODO(u): Add test case for insertSemicolon; ref: go/src/pkg/go/scanner/scanner_test.go:345

// insertSemicolon inserts a semicolon if the correct conditions have been met.
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
func insertSemicolon(l *lexer) {
	insert := false
	trailingComments := false
	var pos int
	for pos = len(l.tokens) - 1; pos >= l.line; pos-- {
		last := l.tokens[pos]
		switch last.Kind {
		case token.Comment:
			// Ignore trailing comments.
			trailingComments = true
			continue
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
		break
	}

	// Insert a semicolon.
	if insert {
		tok := token.Token{
			Kind: token.Semicolon,
			Val:  ";",
		}
		l.tokens = append(l.tokens, tok)

		if trailingComments {
			// Move trailing comments to the end.
			copy(l.tokens[pos+2:], l.tokens[pos+1:])
			// Insert a semicolon before the trailing comments.
			l.tokens[pos+1] = tok
		}
	}
}
