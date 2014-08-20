package lexer

import (
	"strconv"
	"strings"

	"github.com/mewlang/go/token"
)

const (
	// whitespace specifies the white space characters (except newline), which
	// include spaces (U+0020), horizontal tabs (U+0009), and carriage returns
	// (U+000D).
	whitespace = " \t\r"
	// decimal specifies the decimal digit characters.
	decimal = "0123456789"
	// hex specifies the hexadecimal digit characters.
	hex = "0123456789ABCDEFabcdef"
)

// A stateFn represents the state of the lexer as a function that returns a
// state function.
type stateFn func(l *lexer) stateFn

// lexToken lexes a token of the Go programming language. It is the initial
// state function of the lexer.
func lexToken(l *lexer) stateFn {
	// Ignore white space characters (except newline).
	l.ignoreRun(whitespace)

	r := l.next()
	switch r {
	case eof:
		insertSemicolon(l)
		// Emit an EOF and terminate the lexer with a nil state function.
		l.emit(token.EOF)
		return nil
	case '\n':
		l.ignore()
		insertSemicolon(l)
		// Update the index to the first token of the current line.
		l.line = len(l.tokens)
		return lexToken
	case '/':
		return lexDivOrComment
	case '!':
		return lexNot
	case '<':
		return lexLessArrowOrShl
	case '>':
		return lexGreaterOrShr
	case '&':
		return lexAndOrClear
	case '|':
		return lexOr
	case '*':
		return lexMul
	case '%':
		return lexMod
	case '+':
		return lexAddOrInc
	case '-':
		return lexSubOrDec
	case '.', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		l.backup()
		return lexDotOrNumber
	case '\'':
		return lexRune
	case '"':
		return lexString
	case '`':
		return lexRawString
	}

	panic("not yet implemented.")
}

// lexDivOrComment lexes a division operator (/), a division assignment operator
// (/=), a line comment (//), or a general comment (/*). A slash character (/)
// has already been consumed.
func lexDivOrComment(l *lexer) stateFn {
	r := l.next()
	switch r {
	case '=':
		// Division assignment operator (/=).
		l.emit(token.DivAssign)
		return lexToken
	case '/':
		// Line comment (//).
		return lexLineComment
	case '*':
		// General comment (/*).
		return lexGeneralComment
	default:
		// Division operator (/).
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

// lexNot lexes a logical not operator (!), or a not equal comparison operator
// (!=). An exclamation mark character (!) has already been consumed.
func lexNot(l *lexer) stateFn {
	r := l.next()
	switch r {
	case '=':
		// Not equal comparison operator (!=).
		l.emit(token.Neq)
	default:
		// Logical not operator (!).
		l.backup()
		l.emit(token.Not)
	}
	return lexToken
}

// lexLessArrowOrShl lexes a less than comparison operator (<), a less than or
// equal comparison operator (<=), a left shift operator (<<), a left shift
// assignment operator (<<=), or a channel communication operator (<-). A
// less-than sign character (<) has already been consumed.
func lexLessArrowOrShl(l *lexer) stateFn {
	r := l.next()
	switch r {
	case '-':
		// Channel communication operator (<-).
		l.emit(token.Arrow)
	case '<':
		if l.accept("=") {
			// Left shift assignment operator (<<=).
			l.emit(token.ShlAssign)
		} else {
			// Left shift operator (<<).
			l.emit(token.Shl)
		}
	case '=':
		// Less than or equal comparison operator (<=).
		l.emit(token.Lte)
	default:
		// Less than comparison operator (<).
		l.backup()
		l.emit(token.Lt)
	}
	return lexToken
}

// lexGreaterOrShr lexes a greater than comparison operator (>), a greater than
// or equal comparison operator (>=), a right shift operator (>>), or a right
// shift assignment operator (>>=). A greater-than sign character (>) has
// already been consumed.
func lexGreaterOrShr(l *lexer) stateFn {
	r := l.next()
	switch r {
	case '>':
		if l.accept("=") {
			// Right shift assignment operator (>>=).
			l.emit(token.ShrAssign)
		} else {
			// Right shift operator (>>).
			l.emit(token.Shr)
		}
	case '=':
		// Greater than or equal comparison operator (>=).
		l.emit(token.Gte)
	default:
		// Greater than comparison operator (>).
		l.backup()
		l.emit(token.Gt)
	}
	return lexToken
}

// lexAndOrClear lexes a bitwise AND operator (&), a bitwise AND assignment
// operator (&=), a bit clear operator (&^), a bit clear assignment operator
// (&^=), or a logical AND operator (&&). An ampersand character (&) has already
// been consumed.
func lexAndOrClear(l *lexer) stateFn {
	r := l.next()
	switch r {
	case '^':
		if l.accept("=") {
			// Bit clear assignment operator (&^=).
			l.emit(token.ClearAssign)
		} else {
			// Bit clear operator (&^).
			l.emit(token.Clear)
		}
	case '&':
		// Logical AND operator (&&).
		l.emit(token.Land)
	case '=':
		// Bitwise AND assignment operator (&=).
		l.emit(token.AndAssign)
	default:
		// Bitwise AND operator (&).
		l.backup()
		l.emit(token.And)
	}
	return lexToken
}

// lexOr lexes a bitwise OR operator (|), a bitwise OR assignment operator (|=),
// or a logical OR operator (||). A vertical bar character (|) has already been
// consumed.
func lexOr(l *lexer) stateFn {
	r := l.next()
	switch r {
	case '|':
		// Logical OR operator (||).
		l.emit(token.Lor)
	case '=':
		// Bitwise OR assignment operator (|=).
		l.emit(token.OrAssign)
	default:
		// Bitwise OR operator (|).
		l.backup()
		l.emit(token.Or)
	}
	return lexToken
}

// lexMul lexes a multiplication operator (*), or a multiplication assignment
// operator (*=). An asterisk character (*) has already been consumed.
func lexMul(l *lexer) stateFn {
	r := l.next()
	switch r {
	case '=':
		// Multiplication assignment operator (*=).
		l.emit(token.MulAssign)
	default:
		// Multiplication operator (*). The semantical analysis will determine if
		// the token is part of a pointer dereference expression.
		l.backup()
		l.emit(token.Mul)
	}
	return lexToken
}

// lexMod lexes a modulo operator (%), or a modulo assignment operator (%=). A
// percent sign character (%) has already been consumed.
func lexMod(l *lexer) stateFn {
	r := l.next()
	switch r {
	case '=':
		// Modulo assignment operator (%=).
		l.emit(token.ModAssign)
	default:
		// Modulo operator (%).
		l.backup()
		l.emit(token.Mod)
	}
	return lexToken
}

// lexAddOrInc lexes an addition operator (+), an addition assignment operator
// (+=), or an increment statement operator (++). A plus character (+) has
// already been consumed.
func lexAddOrInc(l *lexer) stateFn {
	r := l.next()
	switch r {
	case '=':
		// Addition assignment operator (+=).
		l.emit(token.AddAssign)
	case '+':
		// Increment statement operator (++).
		l.emit(token.Inc)
	default:
		// Addition operator (+). The semantical analysis will determine if the
		// token is part of a positive number as an unary operator.
		l.backup()
		l.emit(token.Add)
	}
	return lexToken
}

// lexSubOrDec lexes a subtraction operator (-), a subtraction assignment
// operator (-=), or a decrement statement operator (--). A minus character (-)
// has already been consumed.
func lexSubOrDec(l *lexer) stateFn {
	r := l.next()
	switch r {
	case '=':
		// Subtraction assignment operator (-=).
		l.emit(token.SubAssign)
	case '-':
		// Decrement statement operator (--).
		l.emit(token.Dec)
	default:
		// Subtraction operator (-). The semantical analysis will determine if the
		// token is part of a negative number as an unary operator.
		l.backup()
		l.emit(token.Sub)
	}
	return lexToken
}

// lexDotOrNumber lexes a dot delimiter (.), an ellipsis delimiter (...), or a
// number (123, 0x7B, 0173, .123, 123.45, 1e-15, 2i).
func lexDotOrNumber(l *lexer) stateFn {
	// Integer part.
	var kind token.Kind
	if l.accept("0") {
		kind = token.Int
		// Early return for hexadecimal constant.
		if l.accept("xX") {
			if !l.acceptRun(hex) {
				return l.errorf("malformed hexadecimal constant")
			}
			l.emit(token.Int)
			return lexToken
		}
	}
	if l.acceptRun(decimal) {
		kind = token.Int
	}

	// Decimal point.
	if l.accept(".") {
		if kind == token.Int {
			kind = token.Float
		} else {
			kind = token.Dot
		}
	}

	// Fraction part.
	if l.acceptRun(decimal) {
		kind = token.Float
	}

	// Early return for dot or ellipsis delimiter.
	if kind == token.Dot {
		if strings.HasPrefix(l.input[l.pos:], "..") {
			l.pos += 2
			l.width = 0
			kind = token.Ellipsis
		}
		l.emit(kind)
		return lexToken
	}

	// Exponent part.
	if l.accept("eE") {
		l.accept("+-")
		if !l.acceptRun(decimal) {
			return l.errorf("malformed exponent of floating-point constant")
		}
	}

	// Imaginary.
	if l.accept("i") {
		kind = token.Imag
	}

	l.emit(kind)
	return lexToken
}

// lexRune lexes a rune literal ('a'). A single quote character (') has already
// been consumed.
func lexRune(l *lexer) stateFn {
	switch l.next() {
	case eof:
		return l.errorf("unexpected eof in rune literal")
	case '\n':
		return l.errorf("unexpected newline in rune literal")
	case '\\':
		// Consume backslash escape sequence.
		l.backup()
		_, multibyte, tail, err := strconv.UnquoteChar(l.input[l.pos:], '\'')
		if err != nil {
			return l.errorf("invalid escape sequence in interpreted string literal; %v", err)
		}
		if multibyte {
			delta := len(l.input[l.pos:]) - len(tail)
			l.pos += delta
			l.width = 0
		} else {
			l.pos++
			l.width = 0
		}
	}
	if !l.accept("'") {
		return l.errorf("missing ' in rune literal")
	}
	l.emit(token.Rune)
	return lexToken
}

// lexString lexes an interpreted string literal ("foo"). A double quote
// character (") has already been consumed.
func lexString(l *lexer) stateFn {
	for {
		switch l.next() {
		case eof:
			return l.errorf("unexpected eof in interpreted string literal")
		case '\n':
			return l.errorf("unexpected newline in interpreted string literal")
		case '\\':
			// Consume backslash escape sequence.
			l.backup()
			_, multibyte, tail, err := strconv.UnquoteChar(l.input[l.pos:], '"')
			if err != nil {
				return l.errorf("invalid escape sequence in interpreted string literal; %v", err)
			}
			if multibyte {
				delta := len(l.input[l.pos:]) - len(tail)
				l.pos += delta
				l.width = 0
			} else {
				l.pos++
				l.width = 0
			}
		case '"':
			l.emit(token.String)
			return lexToken
		}
	}
}

// lexRawString lexes a raw string literal (`foo`). A back quote character (`)
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
