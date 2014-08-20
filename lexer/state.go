package lexer

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/mewlang/go/token"
)

const (
	// whitespace specifies the white space characters (except newline), which
	// include spaces (U+0020), horizontal tabs (U+0009), and carriage returns
	// (U+000D).
	whitespace = " \t\r"
	// decimal specifies the decimal digit characters.
	decimal = "0123456789"
	// octal specifies the octal digit characters.
	octal = "01234567"
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
	case '=':
		return lexEqOrAssign
	case ':':
		return lexColonOrDeclAssign
	case '*':
		return lexMul
	case '%':
		return lexMod
	case '^':
		return lexXor
	case '+':
		return lexAddOrInc
	case '-':
		return lexSubOrDec
	case '(':
		l.emit(token.Lparen)
		return lexToken
	case '[':
		l.emit(token.Lbrack)
		return lexToken
	case '{':
		l.emit(token.Lbrace)
		return lexToken
	case ')':
		l.emit(token.Rparen)
		return lexToken
	case ']':
		l.emit(token.Rbrack)
		return lexToken
	case '}':
		l.emit(token.Rbrace)
		return lexToken
	case ',':
		l.emit(token.Comma)
		return lexToken
	case ';':
		l.emit(token.Semicolon)
		return lexToken
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

	// Check if r is a Unicode letter or an underscore character.
	if isLetter(r) {
		return lexKeywordOrIdent
	}

	return l.errorf("syntax error; unexpected %q", r)
}

// isLetter returns true if r is a Unicode letter or an underscore, and false
// otherwise.
func isLetter(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

// lexDivOrComment lexes a division operator (/), a division assignment operator
// (/=), a line comment (//), or a general comment (/*). A slash character (/)
// has already been consumed.
func lexDivOrComment(l *lexer) stateFn {
	switch l.next() {
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
	switch l.next() {
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
	switch l.next() {
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
	switch l.next() {
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
	switch l.next() {
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
	switch l.next() {
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

// lexEqOrAssign lexes an equal comparison operator (==), or an assignment
// operator (=). An equal sign character (=) has already been consumed.
func lexEqOrAssign(l *lexer) stateFn {
	switch l.next() {
	case '=':
		// Equal comparison operator (==).
		l.emit(token.Eq)
	default:
		// Assignment operator (=).
		l.backup()
		l.emit(token.Assign)
	}
	return lexToken
}

// lexColonOrDeclAssign lexes a colon delimiter (:), or a declare and initialize
// operator (:=). A colon character (:) has already been consumed.
func lexColonOrDeclAssign(l *lexer) stateFn {
	switch l.next() {
	case '=':
		// Declare and initialize operator (:=).
		l.emit(token.DeclAssign)
	default:
		// Colon delimiter (:).
		l.backup()
		l.emit(token.Colon)
	}
	return lexToken
}

// lexMul lexes a multiplication operator (*), or a multiplication assignment
// operator (*=). An asterisk character (*) has already been consumed.
func lexMul(l *lexer) stateFn {
	switch l.next() {
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
	switch l.next() {
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

// lexXor lexes a bitwise XOR operator (^), or a bitwise XOR assignment operator
// (^=). A caret character (^) has already been consumed.
func lexXor(l *lexer) stateFn {
	switch l.next() {
	case '=':
		// Bitwise XOR assignment operator (^=).
		l.emit(token.XorAssign)
	default:
		// Bitwise XOR operator (^). The semantical analysis will determine if the
		// token is part of a bitwise complement expression.
		l.backup()
		l.emit(token.Xor)
	}
	return lexToken
}

// lexAddOrInc lexes an addition operator (+), an addition assignment operator
// (+=), or an increment statement operator (++). A plus character (+) has
// already been consumed.
func lexAddOrInc(l *lexer) stateFn {
	switch l.next() {
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
	switch l.next() {
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
		kind = token.Float

		// Optional sign.
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
		err := consumeEscape(l, '\'')
		if err != nil {
			return l.errorf("invalid escape sequence in interpreted string literal; %v", err)
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
			err := consumeEscape(l, '"')
			if err != nil {
				return l.errorf("invalid escape sequence in interpreted string literal; %v", err)
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

// keywords specifies the reserved keywords of the Go programming language.
var keywords = map[string]token.Kind{
	"break":       token.Break,
	"case":        token.Case,
	"chan":        token.Chan,
	"const":       token.Const,
	"continue":    token.Continue,
	"default":     token.Default,
	"defer":       token.Defer,
	"else":        token.Else,
	"fallthrough": token.Fallthrough,
	"for":         token.For,
	"func":        token.Func,
	"go":          token.Go,
	"goto":        token.Goto,
	"if":          token.If,
	"import":      token.Import,
	"interface":   token.Interface,
	"map":         token.Map,
	"package":     token.Package,
	"range":       token.Range,
	"return":      token.Return,
	"select":      token.Select,
	"struct":      token.Struct,
	"switch":      token.Switch,
	"type":        token.Type,
	"var":         token.Var,
}

// lexKeywordOrIdent lexes a keyword, or an identifier. A Unicode letter or an
// underscore character (_) has already been consumed.
func lexKeywordOrIdent(l *lexer) stateFn {
	for {
		r := l.next()
		if !isLetter(r) && !unicode.IsDigit(r) {
			l.backup()
			break
		}
	}
	s := l.input[l.start:l.pos]
	if kind, ok := keywords[s]; ok {
		l.emit(kind)
	} else {
		l.emit(token.Ident)
	}
	return lexToken
}

// consumeEscape consumes an escape sequence. A valid single-character escape
// sequence is specified by valid. Single quotes are only valid within rune
// literals and double quotes are only valid within string literals. A backslash
// character (\) has already been consumed.
//
// Several backslash escapes allow arbitrary values to be encoded as ASCII text.
// There are four ways to represent the integer value as a numeric constant: \x
// followed by exactly two hexadecimal digits; \u followed by exactly four
// hexadecimal digits; \U followed by exactly eight hexadecimal digits, and a
// plain backslash \ followed by exactly three octal digits. In each case the
// value of the literal is the value represented by the digits in the
// corresponding base.
//
// Although these representations all result in an integer, they have different
// valid ranges. Octal escapes must represent a value between 0 and 255
// inclusive. Hexadecimal escapes satisfy this condition by construction. The
// escapes \u and \U represent Unicode code points so within them some values
// are illegal, in particular those above 0x10FFFF and surrogate halves.
//
// After a backslash, certain single-character escapes represent special values:
//    \a   U+0007 alert or bell
//    \b   U+0008 backspace
//    \f   U+000C form feed
//    \n   U+000A line feed or newline
//    \r   U+000D carriage return
//    \t   U+0009 horizontal tab
//    \v   U+000b vertical tab
//    \\   U+005c backslash
//    \'   U+0027 single quote  (valid escape only within rune literals)
//    \"   U+0022 double quote  (valid escape only within string literals)
//
// All other sequences starting with a backslash are illegal inside rune and
// string literals.
//
// ref: http://golang.org/ref/spec#Rune_literals
func consumeEscape(l *lexer, valid rune) error {
	r := l.next()
	switch r {
	case '0', '1', '2', '3':
		// Octal escape.
		if !l.accept(octal) || !l.accept(octal) {
			return fmt.Errorf("non-octal character %q in octal escape", l.next())
		}
		s := l.input[l.pos-3 : l.pos]
		_, err := strconv.ParseUint(s, 8, 8)
		if err != nil {
			return fmt.Errorf("invalid octal escape; %v", err)
		}
	case 'x':
		// Hexadecimal escape.
		if !l.accept(hex) || !l.accept(hex) {
			return fmt.Errorf("non-hex character %q in hex escape", l.next())
		}
	case 'u', 'U':
		// Unicode escape.
		n := 4
		if r == 'U' {
			n = 8
		}
		for i := 0; i < n; i++ {
			if !l.accept(hex) {
				return fmt.Errorf("non-hex character %q in Unicode escape", l.next())
			}
		}
		s := l.input[l.pos-n : l.pos]
		x, err := strconv.ParseUint(s, 16, 32)
		if err != nil {
			return fmt.Errorf("invalid Unicode escape; %v", err)
		}
		r := rune(x)
		if !utf8.ValidRune(r) {
			return fmt.Errorf("invalid rune %q in Unicode escape", r)
		}
	case 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\', valid:
		// Single-character escape.
	default:
		return fmt.Errorf("unknown escape sequence: %q", r)
	}
	return nil
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
