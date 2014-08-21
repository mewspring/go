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
}

func (tok Token) String() string {
	return fmt.Sprintf("%v %q", tok.Kind, tok.Val)
}

// Kind is the set of lexical token types of the Go programming language. It
// contains four classes of tokens:
//    * identifiers
//    * keywords
//    * operators and delimiters
//    * literals
type Kind uint8

// TODO(u): Evaluate if EOF should be removed, and if Invalid or Illegal should
// be added.

// Token types.
const (
	// Special tokens.
	EOF     Kind = iota // end of file.
	Comment             // line comment or general comment.

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
	//Pos // +
	//Neg // -
	Not // !
	//Comp  // ^
	//Deref // *
	//Addr  // &
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
	EOF:     "EOF",
	Comment: "comment",

	// Identifier.
	Ident: "identifier",

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

	// Literals.
	Int:    "int literal",
	Float:  "float literal",
	Imag:   "imaginary literal",
	Rune:   "rune literal",
	String: "string literal",
}

func (kind Kind) String() string {
	return names[kind]
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
