// Package token defines constants representing the lexical tokens of the Go
// programming language.
package token

// A Token represents a lexical token of the Go programming language.
type Token struct {
	// The token type.
	Kind
	// The string value of the token.
	Val string
}

// Kind is the set of lexical token types of the Go programming language. It
// contains four classes of tokens:
//    * identifiers
//    * keywords
//    * operators and delimiters
//    * literals
type Kind uint8

// Token types.
const (
	// Special tokens.
	Error   Kind = iota // an error occured; value contains the error message.
	EOF                 // end of file.
	Comment             // line comment or general comment.

	// Identifiers.
	Ident // main

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

	// Deliminators.
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

	// Literals.
	Int    // 12345
	Float  // 123.45
	Imag   // 123.45i
	Rune   // 'a'
	String // "abc"
)

// IsKeyword returns true if kind is a keyword, and false otherwise.
func (kind Kind) IsKeyword() bool {
	return Break <= kind && kind <= Var
}

// IsOperator returns true if kind is an operator or a deliminator, and false
// otherwise.
func (kind Kind) IsOperator() bool {
	return Not <= kind && kind <= Ellipsis
}

// IsLiteral returns true if kind is a literal, and false otherwise.
func (kind Kind) IsLiteral() bool {
	return Int <= kind && kind <= String
}
