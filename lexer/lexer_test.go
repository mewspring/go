package lexer

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/mewlang/go/token"
)

// test cases derived from tokens in go/src/pkg/scanner/scanner_test.go
var golden = []struct {
	in   string
	want token.Token
}{
	// Special tokens
	{in: "/* a comment */", want: token.Token{Kind: token.Comment, Val: "/* a comment */"}},
	{in: "// a comment \n", want: token.Token{Kind: token.Comment, Val: "// a comment "}},
	{in: "/*\r*/", want: token.Token{Kind: token.Comment, Val: "/**/"}},
	{in: "//\r\n", want: token.Token{Kind: token.Comment, Val: "//"}},

	// Identifiers and basic type literals
	{in: "foobar", want: token.Token{Kind: token.Ident, Val: "foobar"}},
	{in: "a۰۱۸", want: token.Token{Kind: token.Ident, Val: "a۰۱۸"}},
	{in: "foo६४", want: token.Token{Kind: token.Ident, Val: "foo६४"}},
	{in: "bar９８７６", want: token.Token{Kind: token.Ident, Val: "bar９８７６"}},
	{in: "ŝ", want: token.Token{Kind: token.Ident, Val: "ŝ"}},       // was bug (issue 4000)
	{in: "ŝfoo", want: token.Token{Kind: token.Ident, Val: "ŝfoo"}}, // was bug (issue 4000)
	{in: "0", want: token.Token{Kind: token.Int, Val: "0"}},
	{in: "1", want: token.Token{Kind: token.Int, Val: "1"}},
	{in: "123456789012345678890", want: token.Token{Kind: token.Int, Val: "123456789012345678890"}},
	{in: "01234567", want: token.Token{Kind: token.Int, Val: "01234567"}},
	{in: "0xcafebabe", want: token.Token{Kind: token.Int, Val: "0xcafebabe"}},
	{in: "0.", want: token.Token{Kind: token.Float, Val: "0."}},
	{in: ".0", want: token.Token{Kind: token.Float, Val: ".0"}},
	{in: "3.14159265", want: token.Token{Kind: token.Float, Val: "3.14159265"}},
	{in: "1e0", want: token.Token{Kind: token.Float, Val: "1e0"}},
	{in: "1e+100", want: token.Token{Kind: token.Float, Val: "1e+100"}},
	{in: "1e-100", want: token.Token{Kind: token.Float, Val: "1e-100"}},
	{in: "2.71828e-1000", want: token.Token{Kind: token.Float, Val: "2.71828e-1000"}},
	{in: "0i", want: token.Token{Kind: token.Imag, Val: "0i"}},
	{in: "1i", want: token.Token{Kind: token.Imag, Val: "1i"}},
	{in: "012345678901234567889i", want: token.Token{Kind: token.Imag, Val: "012345678901234567889i"}},
	{in: "123456789012345678890i", want: token.Token{Kind: token.Imag, Val: "123456789012345678890i"}},
	{in: "0.i", want: token.Token{Kind: token.Imag, Val: "0.i"}},
	{in: ".0i", want: token.Token{Kind: token.Imag, Val: ".0i"}},
	{in: "3.14159265i", want: token.Token{Kind: token.Imag, Val: "3.14159265i"}},
	{in: "1e0i", want: token.Token{Kind: token.Imag, Val: "1e0i"}},
	{in: "1e+100i", want: token.Token{Kind: token.Imag, Val: "1e+100i"}},
	{in: "1e-100i", want: token.Token{Kind: token.Imag, Val: "1e-100i"}},
	{in: "2.71828e-1000i", want: token.Token{Kind: token.Imag, Val: "2.71828e-1000i"}},
	{in: "'a'", want: token.Token{Kind: token.Rune, Val: "'a'"}},
	{in: "'\\000'", want: token.Token{Kind: token.Rune, Val: "'\\000'"}},
	{in: "'\\xFF'", want: token.Token{Kind: token.Rune, Val: "'\\xFF'"}},
	{in: "'\\uff16'", want: token.Token{Kind: token.Rune, Val: "'\\uff16'"}},
	{in: "'\\U0000ff16'", want: token.Token{Kind: token.Rune, Val: "'\\U0000ff16'"}},
	{in: "`foobar`", want: token.Token{Kind: token.String, Val: "`foobar`"}},
	{in: `"\a\b\f\n\r\t\v\\\""`, want: token.Token{Kind: token.String, Val: `"\a\b\f\n\r\t\v\\\""`}},
	{in: "`foo\n\t                        bar`", want: token.Token{Kind: token.String, Val: "`foo\n\t                        bar`"}},
	{in: "`\r`", want: token.Token{Kind: token.String, Val: "``"}},
	{in: "`foo\r\nbar`", want: token.Token{Kind: token.String, Val: "`foo\nbar`"}},

	// Operators and delimiters
	{in: "+", want: token.Token{Kind: token.Add, Val: "+"}},
	{in: "-", want: token.Token{Kind: token.Sub, Val: "-"}},
	{in: "*", want: token.Token{Kind: token.Mul, Val: "*"}},
	{in: "/", want: token.Token{Kind: token.Div, Val: "/"}},
	{in: "%", want: token.Token{Kind: token.Mod, Val: "%"}},
	{in: "&", want: token.Token{Kind: token.And, Val: "&"}},
	{in: "|", want: token.Token{Kind: token.Or, Val: "|"}},
	{in: "^", want: token.Token{Kind: token.Xor, Val: "^"}},
	{in: "<<", want: token.Token{Kind: token.Shl, Val: "<<"}},
	{in: ">>", want: token.Token{Kind: token.Shr, Val: ">>"}},
	{in: "&^", want: token.Token{Kind: token.Clear, Val: "&^"}},
	{in: "+=", want: token.Token{Kind: token.AddAssign, Val: "+="}},
	{in: "-=", want: token.Token{Kind: token.SubAssign, Val: "-="}},
	{in: "*=", want: token.Token{Kind: token.MulAssign, Val: "*="}},
	{in: "/=", want: token.Token{Kind: token.DivAssign, Val: "/="}},
	{in: "%=", want: token.Token{Kind: token.ModAssign, Val: "%="}},
	{in: "&=", want: token.Token{Kind: token.AndAssign, Val: "&="}},
	{in: "|=", want: token.Token{Kind: token.OrAssign, Val: "|="}},
	{in: "^=", want: token.Token{Kind: token.XorAssign, Val: "^="}},
	{in: "<<=", want: token.Token{Kind: token.ShlAssign, Val: "<<="}},
	{in: ">>=", want: token.Token{Kind: token.ShrAssign, Val: ">>="}},
	{in: "&^=", want: token.Token{Kind: token.ClearAssign, Val: "&^="}},
	{in: "&&", want: token.Token{Kind: token.Land, Val: "&&"}},
	{in: "||", want: token.Token{Kind: token.Lor, Val: "||"}},
	{in: "<-", want: token.Token{Kind: token.Arrow, Val: "<-"}},
	{in: "++", want: token.Token{Kind: token.Inc, Val: "++"}},
	{in: "--", want: token.Token{Kind: token.Dec, Val: "--"}},
	{in: "==", want: token.Token{Kind: token.Eq, Val: "=="}},
	{in: "<", want: token.Token{Kind: token.Lt, Val: "<"}},
	{in: ">", want: token.Token{Kind: token.Gt, Val: ">"}},
	{in: "=", want: token.Token{Kind: token.Assign, Val: "="}},
	{in: "!", want: token.Token{Kind: token.Not, Val: "!"}},
	{in: "!=", want: token.Token{Kind: token.Neq, Val: "!="}},
	{in: "<=", want: token.Token{Kind: token.Lte, Val: "<="}},
	{in: ">=", want: token.Token{Kind: token.Gte, Val: ">="}},
	{in: ":=", want: token.Token{Kind: token.DeclAssign, Val: ":="}},
	{in: "...", want: token.Token{Kind: token.Ellipsis, Val: "..."}},
	{in: "(", want: token.Token{Kind: token.Lparen, Val: "("}},
	{in: "[", want: token.Token{Kind: token.Lbrack, Val: "["}},
	{in: "{", want: token.Token{Kind: token.Lbrace, Val: "{"}},
	{in: ",", want: token.Token{Kind: token.Comma, Val: ","}},
	{in: ".", want: token.Token{Kind: token.Dot, Val: "."}},
	{in: ")", want: token.Token{Kind: token.Rparen, Val: ")"}},
	{in: "]", want: token.Token{Kind: token.Rbrack, Val: "]"}},
	{in: "}", want: token.Token{Kind: token.Rbrace, Val: "}"}},
	{in: ";", want: token.Token{Kind: token.Semicolon, Val: ";"}},
	{in: ":", want: token.Token{Kind: token.Colon, Val: ":"}},

	// Keywords
	{in: "break", want: token.Token{Kind: token.Break, Val: "break"}},
	{in: "case", want: token.Token{Kind: token.Case, Val: "case"}},
	{in: "chan", want: token.Token{Kind: token.Chan, Val: "chan"}},
	{in: "const", want: token.Token{Kind: token.Const, Val: "const"}},
	{in: "continue", want: token.Token{Kind: token.Continue, Val: "continue"}},
	{in: "default", want: token.Token{Kind: token.Default, Val: "default"}},
	{in: "defer", want: token.Token{Kind: token.Defer, Val: "defer"}},
	{in: "else", want: token.Token{Kind: token.Else, Val: "else"}},
	{in: "fallthrough", want: token.Token{Kind: token.Fallthrough, Val: "fallthrough"}},
	{in: "for", want: token.Token{Kind: token.For, Val: "for"}},
	{in: "func", want: token.Token{Kind: token.Func, Val: "func"}},
	{in: "go", want: token.Token{Kind: token.Go, Val: "go"}},
	{in: "goto", want: token.Token{Kind: token.Goto, Val: "goto"}},
	{in: "if", want: token.Token{Kind: token.If, Val: "if"}},
	{in: "import", want: token.Token{Kind: token.Import, Val: "import"}},
	{in: "interface", want: token.Token{Kind: token.Interface, Val: "interface"}},
	{in: "map", want: token.Token{Kind: token.Map, Val: "map"}},
	{in: "package", want: token.Token{Kind: token.Package, Val: "package"}},
	{in: "range", want: token.Token{Kind: token.Range, Val: "range"}},
	{in: "return", want: token.Token{Kind: token.Return, Val: "return"}},
	{in: "select", want: token.Token{Kind: token.Select, Val: "select"}},
	{in: "struct", want: token.Token{Kind: token.Struct, Val: "struct"}},
	{in: "switch", want: token.Token{Kind: token.Switch, Val: "switch"}},
	{in: "type", want: token.Token{Kind: token.Type, Val: "type"}},
	{in: "var", want: token.Token{Kind: token.Var, Val: "var"}},
}

// source contains each token of golden separated by white space.
var source string

func init() {
	const whitespace = "  \t  \n\n\n" // to separate tokens
	src := new(bytes.Buffer)
	for _, g := range golden {
		src.WriteString(g.in)
		src.WriteString(whitespace)
	}
	source = src.String()
}

func TestParse(t *testing.T) {
	// Disable insertion of semicolons.
	f := insertSemicolon
	insertSemicolon = func(*lexer) {}
	defer func() {
		// Enable insertion of semicolons.
		insertSemicolon = f
	}()

	tokens := Parse(source)
	for i, g := range golden {
		if i >= len(tokens) {
			t.Fatalf("too few tokens; expected %d, got %d.", len(golden), len(tokens))
		}
		got := tokens[i]
		if got != g.want {
			t.Errorf("i=%d: token mismatch; expected %v, got %v.", i, g.want, got)
		}
	}
}

func TestInsertSemicolon(t *testing.T) {
	// test cases derived from lines in go/src/pkg/scanner/scanner_test.go
	golden := []struct {
		in   string
		want []token.Token
	}{
		{in: "", want: []token.Token{}},
		{in: "\ufeff;", want: []token.Token{{Kind: token.Semicolon, Val: ";"}}},                                      // first BOM is ignored; a semicolon is present in the source
		{in: ";", want: []token.Token{{Kind: token.Semicolon, Val: ";"}}},                                            // a semicolon is present in the source
		{in: "foo\n", want: []token.Token{{Kind: token.Ident, Val: "foo"}, {Kind: token.Semicolon, Val: ";"}}},       // a semicolon was automatically inserted.
		{in: "123\n", want: []token.Token{{Kind: token.Int, Val: "123"}, {Kind: token.Semicolon, Val: ";"}}},         // a semicolon was automatically inserted.
		{in: "1.2\n", want: []token.Token{{Kind: token.Float, Val: "1.2"}, {Kind: token.Semicolon, Val: ";"}}},       // a semicolon was automatically inserted.
		{in: "'x'\n", want: []token.Token{{Kind: token.Rune, Val: "'x'"}, {Kind: token.Semicolon, Val: ";"}}},        // a semicolon was automatically inserted.
		{in: `"x"` + "\n", want: []token.Token{{Kind: token.String, Val: `"x"`}, {Kind: token.Semicolon, Val: ";"}}}, // a semicolon was automatically inserted.
		{in: "`x`\n", want: []token.Token{{Kind: token.String, Val: "`x`"}, {Kind: token.Semicolon, Val: ";"}}},      // a semicolon was automatically inserted.

		{in: "+\n", want: []token.Token{{Kind: token.Add, Val: "+"}}},
		{in: "-\n", want: []token.Token{{Kind: token.Sub, Val: "-"}}},
		{in: "*\n", want: []token.Token{{Kind: token.Mul, Val: "*"}}},
		{in: "/\n", want: []token.Token{{Kind: token.Div, Val: "/"}}},
		{in: "%\n", want: []token.Token{{Kind: token.Mod, Val: "%"}}},

		{in: "&\n", want: []token.Token{{Kind: token.And, Val: "&"}}},
		{in: "|\n", want: []token.Token{{Kind: token.Or, Val: "|"}}},
		{in: "^\n", want: []token.Token{{Kind: token.Xor, Val: "^"}}},
		{in: "<<\n", want: []token.Token{{Kind: token.Shl, Val: "<<"}}},
		{in: ">>\n", want: []token.Token{{Kind: token.Shr, Val: ">>"}}},
		{in: "&^\n", want: []token.Token{{Kind: token.Clear, Val: "&^"}}},

		{in: "+=\n", want: []token.Token{{Kind: token.AddAssign, Val: "+="}}},
		{in: "-=\n", want: []token.Token{{Kind: token.SubAssign, Val: "-="}}},
		{in: "*=\n", want: []token.Token{{Kind: token.MulAssign, Val: "*="}}},
		{in: "/=\n", want: []token.Token{{Kind: token.DivAssign, Val: "/="}}},
		{in: "%=\n", want: []token.Token{{Kind: token.ModAssign, Val: "%="}}},

		{in: "&=\n", want: []token.Token{{Kind: token.AndAssign, Val: "&="}}},
		{in: "|=\n", want: []token.Token{{Kind: token.OrAssign, Val: "|="}}},
		{in: "^=\n", want: []token.Token{{Kind: token.XorAssign, Val: "^="}}},
		{in: "<<=\n", want: []token.Token{{Kind: token.ShlAssign, Val: "<<="}}},
		{in: ">>=\n", want: []token.Token{{Kind: token.ShrAssign, Val: ">>="}}},
		{in: "&^=\n", want: []token.Token{{Kind: token.ClearAssign, Val: "&^="}}},

		{in: "&&\n", want: []token.Token{{Kind: token.Land, Val: "&&"}}},
		{in: "||\n", want: []token.Token{{Kind: token.Lor, Val: "||"}}},
		{in: "<-\n", want: []token.Token{{Kind: token.Arrow, Val: "<-"}}},
		{in: "++\n", want: []token.Token{{Kind: token.Inc, Val: "++"}, {Kind: token.Semicolon, Val: ";"}}}, // a semicolon was automatically inserted.
		{in: "--\n", want: []token.Token{{Kind: token.Dec, Val: "--"}, {Kind: token.Semicolon, Val: ";"}}}, // a semicolon was automatically inserted.

		{in: "==\n", want: []token.Token{{Kind: token.Eq, Val: "=="}}},
		{in: "<\n", want: []token.Token{{Kind: token.Lt, Val: "<"}}},
		{in: ">\n", want: []token.Token{{Kind: token.Gt, Val: ">"}}},
		{in: "=\n", want: []token.Token{{Kind: token.Assign, Val: "="}}},
		{in: "!\n", want: []token.Token{{Kind: token.Not, Val: "!"}}},

		{in: "!=\n", want: []token.Token{{Kind: token.Neq, Val: "!="}}},
		{in: "<=\n", want: []token.Token{{Kind: token.Lte, Val: "<="}}},
		{in: ">=\n", want: []token.Token{{Kind: token.Gte, Val: ">="}}},
		{in: ":=\n", want: []token.Token{{Kind: token.DeclAssign, Val: ":="}}},
		{in: "...\n", want: []token.Token{{Kind: token.Ellipsis, Val: "..."}}},

		{in: "(\n", want: []token.Token{{Kind: token.Lparen, Val: "("}}},
		{in: "[\n", want: []token.Token{{Kind: token.Lbrack, Val: "["}}},
		{in: "{\n", want: []token.Token{{Kind: token.Lbrace, Val: "{"}}},
		{in: ",\n", want: []token.Token{{Kind: token.Comma, Val: ","}}},
		{in: ".\n", want: []token.Token{{Kind: token.Dot, Val: "."}}},

		{in: ")\n", want: []token.Token{{Kind: token.Rparen, Val: ")"}, {Kind: token.Semicolon, Val: ";"}}}, // a semicolon was automatically inserted.
		{in: "]\n", want: []token.Token{{Kind: token.Rbrack, Val: "]"}, {Kind: token.Semicolon, Val: ";"}}}, // a semicolon was automatically inserted.
		{in: "}\n", want: []token.Token{{Kind: token.Rbrace, Val: "}"}, {Kind: token.Semicolon, Val: ";"}}}, // a semicolon was automatically inserted.
		{in: ";\n", want: []token.Token{{Kind: token.Semicolon, Val: ";"}}},                                 // a semicolon is present in the source
		{in: ":\n", want: []token.Token{{Kind: token.Colon, Val: ":"}}},

		{in: "break\n", want: []token.Token{{Kind: token.Break, Val: "break"}, {Kind: token.Semicolon, Val: ";"}}}, // a semicolon was automatically inserted.
		{in: "case\n", want: []token.Token{{Kind: token.Case, Val: "case"}}},
		{in: "chan\n", want: []token.Token{{Kind: token.Chan, Val: "chan"}}},
		{in: "const\n", want: []token.Token{{Kind: token.Const, Val: "const"}}},
		{in: "continue\n", want: []token.Token{{Kind: token.Continue, Val: "continue"}, {Kind: token.Semicolon, Val: ";"}}}, // a semicolon was automatically inserted.

		{in: "default\n", want: []token.Token{{Kind: token.Default, Val: "default"}}},
		{in: "defer\n", want: []token.Token{{Kind: token.Defer, Val: "defer"}}},
		{in: "else\n", want: []token.Token{{Kind: token.Else, Val: "else"}}},
		{in: "fallthrough\n", want: []token.Token{{Kind: token.Fallthrough, Val: "fallthrough"}, {Kind: token.Semicolon, Val: ";"}}}, // a semicolon was automatically inserted.
		{in: "for\n", want: []token.Token{{Kind: token.For, Val: "for"}}},

		{in: "func\n", want: []token.Token{{Kind: token.Func, Val: "func"}}},
		{in: "go\n", want: []token.Token{{Kind: token.Go, Val: "go"}}},
		{in: "goto\n", want: []token.Token{{Kind: token.Goto, Val: "goto"}}},
		{in: "if\n", want: []token.Token{{Kind: token.If, Val: "if"}}},
		{in: "import\n", want: []token.Token{{Kind: token.Import, Val: "import"}}},

		{in: "interface\n", want: []token.Token{{Kind: token.Interface, Val: "interface"}}},
		{in: "map\n", want: []token.Token{{Kind: token.Map, Val: "map"}}},
		{in: "package\n", want: []token.Token{{Kind: token.Package, Val: "package"}}},
		{in: "range\n", want: []token.Token{{Kind: token.Range, Val: "range"}}},
		{in: "return\n", want: []token.Token{{Kind: token.Return, Val: "return"}, {Kind: token.Semicolon, Val: ";"}}}, // a semicolon was automatically inserted.

		{in: "select\n", want: []token.Token{{Kind: token.Select, Val: "select"}}},
		{in: "struct\n", want: []token.Token{{Kind: token.Struct, Val: "struct"}}},
		{in: "switch\n", want: []token.Token{{Kind: token.Switch, Val: "switch"}}},
		{in: "type\n", want: []token.Token{{Kind: token.Type, Val: "type"}}},
		{in: "var\n", want: []token.Token{{Kind: token.Var, Val: "var"}}},

		{in: "foo//comment\n", want: []token.Token{{Kind: token.Ident, Val: "foo"}, {Kind: token.Semicolon, Val: ";"}, {Kind: token.Comment, Val: "//comment"}}},         // a semicolon was automatically inserted.
		{in: "foo//comment", want: []token.Token{{Kind: token.Ident, Val: "foo"}, {Kind: token.Semicolon, Val: ";"}, {Kind: token.Comment, Val: "//comment"}}},           // a semicolon was automatically inserted.
		{in: "foo/*comment*/\n", want: []token.Token{{Kind: token.Ident, Val: "foo"}, {Kind: token.Semicolon, Val: ";"}, {Kind: token.Comment, Val: "/*comment*/"}}},     // a semicolon was automatically inserted.
		{in: "foo/*\n*/", want: []token.Token{{Kind: token.Ident, Val: "foo"}, {Kind: token.Semicolon, Val: ";"}, {Kind: token.Comment, Val: "/*\n*/"}}},                 // a semicolon was automatically inserted.
		{in: "foo/*comment*/    \n", want: []token.Token{{Kind: token.Ident, Val: "foo"}, {Kind: token.Semicolon, Val: ";"}, {Kind: token.Comment, Val: "/*comment*/"}}}, // a semicolon was automatically inserted.
		{in: "foo/*\n*/    ", want: []token.Token{{Kind: token.Ident, Val: "foo"}, {Kind: token.Semicolon, Val: ";"}, {Kind: token.Comment, Val: "/*\n*/"}}},             // a semicolon was automatically inserted.

		{in: "foo    // comment\n", want: []token.Token{{Kind: token.Ident, Val: "foo"}, {Kind: token.Semicolon, Val: ";"}, {Kind: token.Comment, Val: "// comment"}}},                                                                                                                                                          // a semicolon was automatically inserted.
		{in: "foo    // comment", want: []token.Token{{Kind: token.Ident, Val: "foo"}, {Kind: token.Semicolon, Val: ";"}, {Kind: token.Comment, Val: "// comment"}}},                                                                                                                                                            // a semicolon was automatically inserted.
		{in: "foo    /*comment*/\n", want: []token.Token{{Kind: token.Ident, Val: "foo"}, {Kind: token.Semicolon, Val: ";"}, {Kind: token.Comment, Val: "/*comment*/"}}},                                                                                                                                                        // a semicolon was automatically inserted.
		{in: "foo    /*\n*/", want: []token.Token{{Kind: token.Ident, Val: "foo"}, {Kind: token.Semicolon, Val: ";"}, {Kind: token.Comment, Val: "/*\n*/"}}},                                                                                                                                                                    // a semicolon was automatically inserted.
		{in: "foo    /*  */ /* \n */ bar/**/\n", want: []token.Token{{Kind: token.Ident, Val: "foo"}, {Kind: token.Semicolon, Val: ";"}, {Kind: token.Comment, Val: "/*  */"}, {Kind: token.Comment, Val: "/* \n */"}, {Kind: token.Ident, Val: "bar"}, {Kind: token.Semicolon, Val: ";"}, {Kind: token.Comment, Val: "/**/"}}}, // a semicolon was automatically inserted.
		{in: "foo    /*0*/ /*1*/ /*2*/\n", want: []token.Token{{Kind: token.Ident, Val: "foo"}, {Kind: token.Semicolon, Val: ";"}, {Kind: token.Comment, Val: "/*0*/"}, {Kind: token.Comment, Val: "/*1*/"}, {Kind: token.Comment, Val: "/*2*/"}}},                                                                              // a semicolon was automatically inserted.

		{in: "foo    /*comment*/    \n", want: []token.Token{{Kind: token.Ident, Val: "foo"}, {Kind: token.Semicolon, Val: ";"}, {Kind: token.Comment, Val: "/*comment*/"}}},                                                                           // a semicolon was automatically inserted.
		{in: "foo    /*0*/ /*1*/ /*2*/    \n", want: []token.Token{{Kind: token.Ident, Val: "foo"}, {Kind: token.Semicolon, Val: ";"}, {Kind: token.Comment, Val: "/*0*/"}, {Kind: token.Comment, Val: "/*1*/"}, {Kind: token.Comment, Val: "/*2*/"}}}, // a semicolon was automatically inserted.
		{in: "foo	/**/ /*-------------*/       /*----\n*/bar       /*  \n*/baa\n", want: []token.Token{{Kind: token.Ident, Val: "foo"}, {Kind: token.Semicolon, Val: ";"}, {Kind: token.Comment, Val: "/**/"}, {Kind: token.Comment, Val: "/*-------------*/"}, {Kind: token.Comment, Val: "/*----\n*/"}, {Kind: token.Ident, Val: "bar"}, {Kind: token.Semicolon, Val: ";"}, {Kind: token.Comment, Val: "/*  \n*/"}, {Kind: token.Ident, Val: "baa"}, {Kind: token.Semicolon, Val: ";"}}}, // a semicolon was automatically inserted.
		{in: "foo    /* an EOF terminates a line */", want: []token.Token{{Kind: token.Ident, Val: "foo"}, {Kind: token.Semicolon, Val: ";"}, {Kind: token.Comment, Val: "/* an EOF terminates a line */"}}},                                                                   // a semicolon was automatically inserted.
		{in: "foo    /* an EOF terminates a line */ /*", want: []token.Token{{Kind: token.Ident, Val: "foo"}, {Kind: token.Semicolon, Val: ";"}, {Kind: token.Comment, Val: "/* an EOF terminates a line */"}, {Kind: token.Error, Val: "unexpected eof in general comment"}}}, // a semicolon was automatically inserted.
		{in: "foo    /* an EOF terminates a line */ //", want: []token.Token{{Kind: token.Ident, Val: "foo"}, {Kind: token.Semicolon, Val: ";"}, {Kind: token.Comment, Val: "/* an EOF terminates a line */"}, {Kind: token.Comment, Val: "//"}}},                              // a semicolon was automatically inserted.

		{in: "package main\n\nfunc main() {\n\tif {\n\t\treturn /* */ }\n}\n", want: []token.Token{{Kind: token.Package, Val: "package"}, {Kind: token.Ident, Val: "main"}, {Kind: token.Semicolon, Val: ";"}, {Kind: token.Func, Val: "func"}, {Kind: token.Ident, Val: "main"}, {Kind: token.Lparen, Val: "("}, {Kind: token.Rparen, Val: ")"}, {Kind: token.Lbrace, Val: "{"}, {Kind: token.If, Val: "if"}, {Kind: token.Lbrace, Val: "{"}, {Kind: token.Return, Val: "return"}, {Kind: token.Comment, Val: "/* */"}, {Kind: token.Rbrace, Val: "}"}, {Kind: token.Semicolon, Val: ";"}, {Kind: token.Rbrace, Val: "}"}, {Kind: token.Semicolon, Val: ";"}}}, // a semicolon was automatically inserted.
		{in: "package main", want: []token.Token{{Kind: token.Package, Val: "package"}, {Kind: token.Ident, Val: "main"}, {Kind: token.Semicolon, Val: ";"}}}, // a semicolon was automatically inserted.
	}

	for i, g := range golden {
		got := Parse(g.in)
		if len(got) > 0 {
			tok := got[len(got)-1]
			if tok.Kind == token.EOF {
				got = got[:len(got)-1]
			}
		}
		if !reflect.DeepEqual(got, g.want) {
			t.Errorf("i=%d: expected %v, got %v.", i, g.want, got)
		}
	}
}

func BenchmarkParse(b *testing.B) {
	b.SetBytes(int64(len(source)))
	for i := 0; i < b.N; i++ {
		Parse(source)
	}
}
