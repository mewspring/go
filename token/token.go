// Package token defines constants representing the lexical tokens of the Go
// programming language.
package token

import "fmt"

// A Token represents a lexical token of the Go programming language.
type Token struct {
	// The token type.
	Kind
	// The string value of the token.
	Val string
	// Line number, starting at 1.
	Line int
	// Column number, starting at 1 (character count).
	Col int
}

func (tok Token) String() string {
	if tok.Kind == None {
		return "NONE"
	}
	if !tok.IsValid() {
		return fmt.Sprintf("<invalid> %s", tok.Val)
	}
	return tok.Val
}

// Kind is the set of lexical token types of the Go programming language. It
// contains four classes of tokens:
//    * identifiers
//    * keywords
//    * operators and delimiters
//    * literals
//
// A token is lexically invalid if its least significant bit is set.
type Kind uint8

// Token types.
const (
	// Special tokens.
	None    Kind = 0         // NONE option for tokens.
	Invalid Kind = 1         // invalid token; e.g. an unterminated rune literal.
	Comment Kind = iota << 1 // line comment or block comment.

	// Identifiers and literals.
	// Identifier.
	Ident // main

	// Literals.
	Int    // 12345
	Float  // 123.45
	Imag   // 123.45i
	Rune   // 'a'
	String // "abc"

	// Keywords.
	Break       // break
	Case        // case
	Chan        // chan
	Const       // const
	Continue    // continue
	Default     // default
	Defer       // defer
	Else        // else
	Fallthrough // fallthrough
	For         // for
	Func        // func
	Go          // go
	Goto        // goto
	If          // if
	Import      // import
	Interface   // interface
	Map         // map
	Package     // package
	Range       // range
	Return      // return
	Select      // select
	Struct      // struct
	Switch      // switch
	Type        // type
	Var         // var

	// Operators and delimiters.
	// Unary operators.
	Not   // !
	Arrow // <-

	// Operators with precedence 5.
	Mul   // *
	Div   // /
	Mod   // %
	Shl   // <<
	Shr   // >>
	And   // &
	Clear // &^

	// Operators with precedence 4.
	Add // +
	Sub // -
	Or  // |
	Xor // ^

	// Operators with precedence 3.
	Eq  // ==
	Neq // !=
	Lt  // <
	Lte // <=
	Gt  // >
	Gte // >=

	// Operators with precedence 2.
	Land // &&

	// Operators with precedence 1.
	Lor // ||

	// Assignment operators.
	Assign      // =
	DeclAssign  // :=
	MulAssign   // *=
	DivAssign   // /=
	ModAssign   // %=
	ShlAssign   // <<=
	ShrAssign   // >>=
	AndAssign   // &=
	ClearAssign // &^=
	AddAssign   // +=
	SubAssign   // -=
	OrAssign    // |=
	XorAssign   // ^=

	// Statement operators.
	Inc // ++
	Dec // --

	// Delimiters.
	Lparen    // (
	Lbrack    // [
	Lbrace    // {
	Rparen    // )
	Rbrack    // ]
	Rbrace    // }
	Dot       // .
	Comma     // ,
	Colon     // :
	Semicolon // ;
	Ellipsis  // ...
)

// names specifies the name of each token type.
var names = [...]string{
	// Special.
	Invalid: "<invalid>",
	Comment: "comment",

	// Identifiers and literals.
	Ident:  "identifier",
	Int:    "int literal",
	Float:  "float literal",
	Imag:   "imaginary literal",
	Rune:   "rune literal",
	String: "string literal",

	// Keywords.
	Break:       "break",
	Case:        "case",
	Chan:        "chan",
	Const:       "const",
	Continue:    "continue",
	Default:     "default",
	Defer:       "defer",
	Else:        "else",
	Fallthrough: "fallthrough",
	For:         "for",
	Func:        "func",
	Go:          "go",
	Goto:        "goto",
	If:          "if",
	Import:      "import",
	Interface:   "interface",
	Map:         "map",
	Package:     "package",
	Range:       "range",
	Return:      "return",
	Select:      "select",
	Struct:      "struct",
	Switch:      "switch",
	Type:        "type",
	Var:         "var",

	// Operators and delimiters.
	Not:         "!",
	Arrow:       "<-",
	Mul:         "*",
	Div:         "/",
	Mod:         "%",
	Shl:         "<<",
	Shr:         ">>",
	And:         "&",
	Clear:       "&^",
	Add:         "+",
	Sub:         "-",
	Or:          "|",
	Xor:         "^",
	Eq:          "==",
	Neq:         "!=",
	Lt:          "<",
	Lte:         "<=",
	Gt:          ">",
	Gte:         ">=",
	Land:        "&&",
	Lor:         "||",
	Assign:      "=",
	DeclAssign:  ":=",
	MulAssign:   "*=",
	DivAssign:   "/=",
	ModAssign:   "%=",
	ShlAssign:   "<<=",
	ShrAssign:   ">>=",
	AndAssign:   "&=",
	ClearAssign: "&^=",
	AddAssign:   "+=",
	SubAssign:   "-=",
	OrAssign:    "|=",
	XorAssign:   "^=",
	Inc:         "++",
	Dec:         "--",
	Lparen:      "(",
	Lbrack:      "[",
	Lbrace:      "{",
	Rparen:      ")",
	Rbrack:      "]",
	Rbrace:      "}",
	Dot:         ".",
	Comma:       ",",
	Colon:       ":",
	Semicolon:   ";",
	Ellipsis:    "...",
}

func (kind Kind) String() string {
	if !kind.IsValid() {
		kind &^= Invalid
		return "<invalid> " + names[kind]
	}
	return names[kind]
}

// IsValid returns true if the token is lexically valid, and false otherwise.
func (kind Kind) IsValid() bool {
	return kind&Invalid == 0
}

// IsKeyword returns true if kind is a keyword, and false otherwise.
func (kind Kind) IsKeyword() bool {
	return Break <= kind && kind <= Var
}

// IsOperator returns true if kind is an operator or a delimiter, and false
// otherwise.
func (kind Kind) IsOperator() bool {
	return Not <= kind && kind <= Ellipsis
}

// IsLiteral returns true if kind is an identifier or a basic literal, and false
// otherwise.
func (kind Kind) IsLiteral() bool {
	return Ident <= kind && kind <= String
}
