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
	{in: "/* a comment */", want: token.Token{Kind: token.Comment, Val: "/* a comment */", Line: 1, Col: 1}},
	{in: "// a comment \n", want: token.Token{Kind: token.Comment, Val: "// a comment ", Line: 4, Col: 1}},
	{in: "/*\r*/", want: token.Token{Kind: token.Comment, Val: "/**/", Line: 8, Col: 1}},
	{in: "//\r\n", want: token.Token{Kind: token.Comment, Val: "//", Line: 11, Col: 1}},

	// Identifiers and basic type literals
	{in: "foobar", want: token.Token{Kind: token.Ident, Val: "foobar", Line: 15, Col: 1}},
	{in: "a۰۱۸", want: token.Token{Kind: token.Ident, Val: "a۰۱۸", Line: 18, Col: 1}},
	{in: "foo६४", want: token.Token{Kind: token.Ident, Val: "foo६४", Line: 21, Col: 1}},
	{in: "bar９８７６", want: token.Token{Kind: token.Ident, Val: "bar９８７６", Line: 24, Col: 1}},
	{in: "ŝ", want: token.Token{Kind: token.Ident, Val: "ŝ", Line: 27, Col: 1}},       // was bug (issue 4000)
	{in: "ŝfoo", want: token.Token{Kind: token.Ident, Val: "ŝfoo", Line: 30, Col: 1}}, // was bug (issue 4000)
	{in: "0", want: token.Token{Kind: token.Int, Val: "0", Line: 33, Col: 1}},
	{in: "1", want: token.Token{Kind: token.Int, Val: "1", Line: 36, Col: 1}},
	{in: "123456789012345678890", want: token.Token{Kind: token.Int, Val: "123456789012345678890", Line: 39, Col: 1}},
	{in: "01234567", want: token.Token{Kind: token.Int, Val: "01234567", Line: 42, Col: 1}},
	{in: "0xcafebabe", want: token.Token{Kind: token.Int, Val: "0xcafebabe", Line: 45, Col: 1}},
	{in: "0.", want: token.Token{Kind: token.Float, Val: "0.", Line: 48, Col: 1}},
	{in: ".0", want: token.Token{Kind: token.Float, Val: ".0", Line: 51, Col: 1}},
	{in: "3.14159265", want: token.Token{Kind: token.Float, Val: "3.14159265", Line: 54, Col: 1}},
	{in: "1e0", want: token.Token{Kind: token.Float, Val: "1e0", Line: 57, Col: 1}},
	{in: "1e+100", want: token.Token{Kind: token.Float, Val: "1e+100", Line: 60, Col: 1}},
	{in: "1e-100", want: token.Token{Kind: token.Float, Val: "1e-100", Line: 63, Col: 1}},
	{in: "2.71828e-1000", want: token.Token{Kind: token.Float, Val: "2.71828e-1000", Line: 66, Col: 1}},
	{in: "0i", want: token.Token{Kind: token.Imag, Val: "0i", Line: 69, Col: 1}},
	{in: "1i", want: token.Token{Kind: token.Imag, Val: "1i", Line: 72, Col: 1}},
	{in: "012345678901234567889i", want: token.Token{Kind: token.Imag, Val: "012345678901234567889i", Line: 75, Col: 1}},
	{in: "123456789012345678890i", want: token.Token{Kind: token.Imag, Val: "123456789012345678890i", Line: 78, Col: 1}},
	{in: "0.i", want: token.Token{Kind: token.Imag, Val: "0.i", Line: 81, Col: 1}},
	{in: ".0i", want: token.Token{Kind: token.Imag, Val: ".0i", Line: 84, Col: 1}},
	{in: "3.14159265i", want: token.Token{Kind: token.Imag, Val: "3.14159265i", Line: 87, Col: 1}},
	{in: "1e0i", want: token.Token{Kind: token.Imag, Val: "1e0i", Line: 90, Col: 1}},
	{in: "1e+100i", want: token.Token{Kind: token.Imag, Val: "1e+100i", Line: 93, Col: 1}},
	{in: "1e-100i", want: token.Token{Kind: token.Imag, Val: "1e-100i", Line: 96, Col: 1}},
	{in: "2.71828e-1000i", want: token.Token{Kind: token.Imag, Val: "2.71828e-1000i", Line: 99, Col: 1}},
	{in: "'a'", want: token.Token{Kind: token.Rune, Val: "'a'", Line: 102, Col: 1}},
	{in: "'\\000'", want: token.Token{Kind: token.Rune, Val: "'\\000'", Line: 105, Col: 1}},
	{in: "'\\xFF'", want: token.Token{Kind: token.Rune, Val: "'\\xFF'", Line: 108, Col: 1}},
	{in: "'\\uff16'", want: token.Token{Kind: token.Rune, Val: "'\\uff16'", Line: 111, Col: 1}},
	{in: "'\\U0000ff16'", want: token.Token{Kind: token.Rune, Val: "'\\U0000ff16'", Line: 114, Col: 1}},
	{in: "`foobar`", want: token.Token{Kind: token.String, Val: "`foobar`", Line: 117, Col: 1}},
	{in: `"\a\b\f\n\r\t\v\\\""`, want: token.Token{Kind: token.String, Val: `"\a\b\f\n\r\t\v\\\""`, Line: 120, Col: 1}},
	{in: "`foo\n\t                        bar`", want: token.Token{Kind: token.String, Val: "`foo\n\t                        bar`", Line: 123, Col: 1}},
	{in: "`\r`", want: token.Token{Kind: token.String, Val: "``", Line: 127, Col: 1}},
	{in: "`foo\r\nbar`", want: token.Token{Kind: token.String, Val: "`foo\nbar`", Line: 130, Col: 1}},

	// Operators and delimiters
	{in: "+", want: token.Token{Kind: token.Add, Val: "+", Line: 134, Col: 1}},
	{in: "-", want: token.Token{Kind: token.Sub, Val: "-", Line: 137, Col: 1}},
	{in: "*", want: token.Token{Kind: token.Mul, Val: "*", Line: 140, Col: 1}},
	{in: "/", want: token.Token{Kind: token.Div, Val: "/", Line: 143, Col: 1}},
	{in: "%", want: token.Token{Kind: token.Mod, Val: "%", Line: 146, Col: 1}},
	{in: "&", want: token.Token{Kind: token.And, Val: "&", Line: 149, Col: 1}},
	{in: "|", want: token.Token{Kind: token.Or, Val: "|", Line: 152, Col: 1}},
	{in: "^", want: token.Token{Kind: token.Xor, Val: "^", Line: 155, Col: 1}},
	{in: "<<", want: token.Token{Kind: token.Shl, Val: "<<", Line: 158, Col: 1}},
	{in: ">>", want: token.Token{Kind: token.Shr, Val: ">>", Line: 161, Col: 1}},
	{in: "&^", want: token.Token{Kind: token.Clear, Val: "&^", Line: 164, Col: 1}},
	{in: "+=", want: token.Token{Kind: token.AddAssign, Val: "+=", Line: 167, Col: 1}},
	{in: "-=", want: token.Token{Kind: token.SubAssign, Val: "-=", Line: 170, Col: 1}},
	{in: "*=", want: token.Token{Kind: token.MulAssign, Val: "*=", Line: 173, Col: 1}},
	{in: "/=", want: token.Token{Kind: token.DivAssign, Val: "/=", Line: 176, Col: 1}},
	{in: "%=", want: token.Token{Kind: token.ModAssign, Val: "%=", Line: 179, Col: 1}},
	{in: "&=", want: token.Token{Kind: token.AndAssign, Val: "&=", Line: 182, Col: 1}},
	{in: "|=", want: token.Token{Kind: token.OrAssign, Val: "|=", Line: 185, Col: 1}},
	{in: "^=", want: token.Token{Kind: token.XorAssign, Val: "^=", Line: 188, Col: 1}},
	{in: "<<=", want: token.Token{Kind: token.ShlAssign, Val: "<<=", Line: 191, Col: 1}},
	{in: ">>=", want: token.Token{Kind: token.ShrAssign, Val: ">>=", Line: 194, Col: 1}},
	{in: "&^=", want: token.Token{Kind: token.ClearAssign, Val: "&^=", Line: 197, Col: 1}},
	{in: "&&", want: token.Token{Kind: token.Land, Val: "&&", Line: 200, Col: 1}},
	{in: "||", want: token.Token{Kind: token.Lor, Val: "||", Line: 203, Col: 1}},
	{in: "<-", want: token.Token{Kind: token.Arrow, Val: "<-", Line: 206, Col: 1}},
	{in: "++", want: token.Token{Kind: token.Inc, Val: "++", Line: 209, Col: 1}},
	{in: "--", want: token.Token{Kind: token.Dec, Val: "--", Line: 212, Col: 1}},
	{in: "==", want: token.Token{Kind: token.Eq, Val: "==", Line: 215, Col: 1}},
	{in: "<", want: token.Token{Kind: token.Lt, Val: "<", Line: 218, Col: 1}},
	{in: ">", want: token.Token{Kind: token.Gt, Val: ">", Line: 221, Col: 1}},
	{in: "=", want: token.Token{Kind: token.Assign, Val: "=", Line: 224, Col: 1}},
	{in: "!", want: token.Token{Kind: token.Not, Val: "!", Line: 227, Col: 1}},
	{in: "!=", want: token.Token{Kind: token.Neq, Val: "!=", Line: 230, Col: 1}},
	{in: "<=", want: token.Token{Kind: token.Lte, Val: "<=", Line: 233, Col: 1}},
	{in: ">=", want: token.Token{Kind: token.Gte, Val: ">=", Line: 236, Col: 1}},
	{in: ":=", want: token.Token{Kind: token.DeclAssign, Val: ":=", Line: 239, Col: 1}},
	{in: "...", want: token.Token{Kind: token.Ellipsis, Val: "...", Line: 242, Col: 1}},
	{in: "(", want: token.Token{Kind: token.Lparen, Val: "(", Line: 245, Col: 1}},
	{in: "[", want: token.Token{Kind: token.Lbrack, Val: "[", Line: 248, Col: 1}},
	{in: "{", want: token.Token{Kind: token.Lbrace, Val: "{", Line: 251, Col: 1}},
	{in: ",", want: token.Token{Kind: token.Comma, Val: ",", Line: 254, Col: 1}},
	{in: ".", want: token.Token{Kind: token.Dot, Val: ".", Line: 257, Col: 1}},
	{in: ")", want: token.Token{Kind: token.Rparen, Val: ")", Line: 260, Col: 1}},
	{in: "]", want: token.Token{Kind: token.Rbrack, Val: "]", Line: 263, Col: 1}},
	{in: "}", want: token.Token{Kind: token.Rbrace, Val: "}", Line: 266, Col: 1}},
	{in: ";", want: token.Token{Kind: token.Semicolon, Val: ";", Line: 269, Col: 1}},
	{in: ":", want: token.Token{Kind: token.Colon, Val: ":", Line: 272, Col: 1}},

	// Keywords
	{in: "break", want: token.Token{Kind: token.Break, Val: "break", Line: 275, Col: 1}},
	{in: "case", want: token.Token{Kind: token.Case, Val: "case", Line: 278, Col: 1}},
	{in: "chan", want: token.Token{Kind: token.Chan, Val: "chan", Line: 281, Col: 1}},
	{in: "const", want: token.Token{Kind: token.Const, Val: "const", Line: 284, Col: 1}},
	{in: "continue", want: token.Token{Kind: token.Continue, Val: "continue", Line: 287, Col: 1}},
	{in: "default", want: token.Token{Kind: token.Default, Val: "default", Line: 290, Col: 1}},
	{in: "defer", want: token.Token{Kind: token.Defer, Val: "defer", Line: 293, Col: 1}},
	{in: "else", want: token.Token{Kind: token.Else, Val: "else", Line: 296, Col: 1}},
	{in: "fallthrough", want: token.Token{Kind: token.Fallthrough, Val: "fallthrough", Line: 299, Col: 1}},
	{in: "for", want: token.Token{Kind: token.For, Val: "for", Line: 302, Col: 1}},
	{in: "func", want: token.Token{Kind: token.Func, Val: "func", Line: 305, Col: 1}},
	{in: "go", want: token.Token{Kind: token.Go, Val: "go", Line: 308, Col: 1}},
	{in: "goto", want: token.Token{Kind: token.Goto, Val: "goto", Line: 311, Col: 1}},
	{in: "if", want: token.Token{Kind: token.If, Val: "if", Line: 314, Col: 1}},
	{in: "import", want: token.Token{Kind: token.Import, Val: "import", Line: 317, Col: 1}},
	{in: "interface", want: token.Token{Kind: token.Interface, Val: "interface", Line: 320, Col: 1}},
	{in: "map", want: token.Token{Kind: token.Map, Val: "map", Line: 323, Col: 1}},
	{in: "package", want: token.Token{Kind: token.Package, Val: "package", Line: 326, Col: 1}},
	{in: "range", want: token.Token{Kind: token.Range, Val: "range", Line: 329, Col: 1}},
	{in: "return", want: token.Token{Kind: token.Return, Val: "return", Line: 332, Col: 1}},
	{in: "select", want: token.Token{Kind: token.Select, Val: "select", Line: 335, Col: 1}},
	{in: "struct", want: token.Token{Kind: token.Struct, Val: "struct", Line: 338, Col: 1}},
	{in: "switch", want: token.Token{Kind: token.Switch, Val: "switch", Line: 341, Col: 1}},
	{in: "type", want: token.Token{Kind: token.Type, Val: "type", Line: 344, Col: 1}},
	{in: "var", want: token.Token{Kind: token.Var, Val: "var", Line: 347, Col: 1}},
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

	tokens, err := Parse(source)
	if err != nil {
		t.Fatalf("Parse failed; %v", err)
	}
	for i, g := range golden {
		if i >= len(tokens) {
			t.Fatalf("i=%d: too few tokens; expected >= %d, got %d.", i, len(golden), len(tokens))
		}
		got := tokens[i]
		if got != g.want {
			t.Errorf("i=%d: token mismatch; expected %#v, got %#v.", i, g.want, got)
		}
	}
}

func TestParseInsertSemicolon(t *testing.T) {
	// test cases derived from lines in go/src/pkg/scanner/scanner_test.go
	golden := []struct {
		in   string
		err  string
		want []token.Token
	}{
		{in: "", want: []token.Token{}},
		{in: "\ufeff;", want: []token.Token{{Kind: token.Semicolon, Val: ";", Line: 1, Col: 1}}},                                                       // first BOM is ignored; a semicolon is present in the source
		{in: ";", want: []token.Token{{Kind: token.Semicolon, Val: ";", Line: 1, Col: 1}}},                                                             // a semicolon is present in the source
		{in: "foo\n", want: []token.Token{{Kind: token.Ident, Val: "foo", Line: 1, Col: 1}, {Kind: token.Semicolon, Val: ";", Line: 1, Col: 4}}},       // a semicolon was automatically inserted.
		{in: "123\n", want: []token.Token{{Kind: token.Int, Val: "123", Line: 1, Col: 1}, {Kind: token.Semicolon, Val: ";", Line: 1, Col: 4}}},         // a semicolon was automatically inserted.
		{in: "1.2\n", want: []token.Token{{Kind: token.Float, Val: "1.2", Line: 1, Col: 1}, {Kind: token.Semicolon, Val: ";", Line: 1, Col: 4}}},       // a semicolon was automatically inserted.
		{in: "'x'\n", want: []token.Token{{Kind: token.Rune, Val: "'x'", Line: 1, Col: 1}, {Kind: token.Semicolon, Val: ";", Line: 1, Col: 4}}},        // a semicolon was automatically inserted.
		{in: `"x"` + "\n", want: []token.Token{{Kind: token.String, Val: `"x"`, Line: 1, Col: 1}, {Kind: token.Semicolon, Val: ";", Line: 1, Col: 4}}}, // a semicolon was automatically inserted.
		{in: "`x`\n", want: []token.Token{{Kind: token.String, Val: "`x`", Line: 1, Col: 1}, {Kind: token.Semicolon, Val: ";", Line: 1, Col: 4}}},      // a semicolon was automatically inserted.

		{in: "+\n", want: []token.Token{{Kind: token.Add, Val: "+", Line: 1, Col: 1}}},
		{in: "-\n", want: []token.Token{{Kind: token.Sub, Val: "-", Line: 1, Col: 1}}},
		{in: "*\n", want: []token.Token{{Kind: token.Mul, Val: "*", Line: 1, Col: 1}}},
		{in: "/\n", want: []token.Token{{Kind: token.Div, Val: "/", Line: 1, Col: 1}}},
		{in: "%\n", want: []token.Token{{Kind: token.Mod, Val: "%", Line: 1, Col: 1}}},

		{in: "&\n", want: []token.Token{{Kind: token.And, Val: "&", Line: 1, Col: 1}}},
		{in: "|\n", want: []token.Token{{Kind: token.Or, Val: "|", Line: 1, Col: 1}}},
		{in: "^\n", want: []token.Token{{Kind: token.Xor, Val: "^", Line: 1, Col: 1}}},
		{in: "<<\n", want: []token.Token{{Kind: token.Shl, Val: "<<", Line: 1, Col: 1}}},
		{in: ">>\n", want: []token.Token{{Kind: token.Shr, Val: ">>", Line: 1, Col: 1}}},
		{in: "&^\n", want: []token.Token{{Kind: token.Clear, Val: "&^", Line: 1, Col: 1}}},

		{in: "+=\n", want: []token.Token{{Kind: token.AddAssign, Val: "+=", Line: 1, Col: 1}}},
		{in: "-=\n", want: []token.Token{{Kind: token.SubAssign, Val: "-=", Line: 1, Col: 1}}},
		{in: "*=\n", want: []token.Token{{Kind: token.MulAssign, Val: "*=", Line: 1, Col: 1}}},
		{in: "/=\n", want: []token.Token{{Kind: token.DivAssign, Val: "/=", Line: 1, Col: 1}}},
		{in: "%=\n", want: []token.Token{{Kind: token.ModAssign, Val: "%=", Line: 1, Col: 1}}},

		{in: "&=\n", want: []token.Token{{Kind: token.AndAssign, Val: "&=", Line: 1, Col: 1}}},
		{in: "|=\n", want: []token.Token{{Kind: token.OrAssign, Val: "|=", Line: 1, Col: 1}}},
		{in: "^=\n", want: []token.Token{{Kind: token.XorAssign, Val: "^=", Line: 1, Col: 1}}},
		{in: "<<=\n", want: []token.Token{{Kind: token.ShlAssign, Val: "<<=", Line: 1, Col: 1}}},
		{in: ">>=\n", want: []token.Token{{Kind: token.ShrAssign, Val: ">>=", Line: 1, Col: 1}}},
		{in: "&^=\n", want: []token.Token{{Kind: token.ClearAssign, Val: "&^=", Line: 1, Col: 1}}},

		{in: "&&\n", want: []token.Token{{Kind: token.Land, Val: "&&", Line: 1, Col: 1}}},
		{in: "||\n", want: []token.Token{{Kind: token.Lor, Val: "||", Line: 1, Col: 1}}},
		{in: "<-\n", want: []token.Token{{Kind: token.Arrow, Val: "<-", Line: 1, Col: 1}}},
		{in: "++\n", want: []token.Token{{Kind: token.Inc, Val: "++", Line: 1, Col: 1}, {Kind: token.Semicolon, Val: ";", Line: 1, Col: 3}}}, // a semicolon was automatically inserted.
		{in: "--\n", want: []token.Token{{Kind: token.Dec, Val: "--", Line: 1, Col: 1}, {Kind: token.Semicolon, Val: ";", Line: 1, Col: 3}}}, // a semicolon was automatically inserted.

		{in: "==\n", want: []token.Token{{Kind: token.Eq, Val: "==", Line: 1, Col: 1}}},
		{in: "<\n", want: []token.Token{{Kind: token.Lt, Val: "<", Line: 1, Col: 1}}},
		{in: ">\n", want: []token.Token{{Kind: token.Gt, Val: ">", Line: 1, Col: 1}}},
		{in: "=\n", want: []token.Token{{Kind: token.Assign, Val: "=", Line: 1, Col: 1}}},
		{in: "!\n", want: []token.Token{{Kind: token.Not, Val: "!", Line: 1, Col: 1}}},

		{in: "!=\n", want: []token.Token{{Kind: token.Neq, Val: "!=", Line: 1, Col: 1}}},
		{in: "<=\n", want: []token.Token{{Kind: token.Lte, Val: "<=", Line: 1, Col: 1}}},
		{in: ">=\n", want: []token.Token{{Kind: token.Gte, Val: ">=", Line: 1, Col: 1}}},
		{in: ":=\n", want: []token.Token{{Kind: token.DeclAssign, Val: ":=", Line: 1, Col: 1}}},
		{in: "...\n", want: []token.Token{{Kind: token.Ellipsis, Val: "...", Line: 1, Col: 1}}},

		{in: "(\n", want: []token.Token{{Kind: token.Lparen, Val: "(", Line: 1, Col: 1}}},
		{in: "[\n", want: []token.Token{{Kind: token.Lbrack, Val: "[", Line: 1, Col: 1}}},
		{in: "{\n", want: []token.Token{{Kind: token.Lbrace, Val: "{", Line: 1, Col: 1}}},
		{in: ",\n", want: []token.Token{{Kind: token.Comma, Val: ",", Line: 1, Col: 1}}},
		{in: ".\n", want: []token.Token{{Kind: token.Dot, Val: ".", Line: 1, Col: 1}}},

		{in: ")\n", want: []token.Token{{Kind: token.Rparen, Val: ")", Line: 1, Col: 1}, {Kind: token.Semicolon, Val: ";", Line: 1, Col: 2}}}, // a semicolon was automatically inserted.
		{in: "]\n", want: []token.Token{{Kind: token.Rbrack, Val: "]", Line: 1, Col: 1}, {Kind: token.Semicolon, Val: ";", Line: 1, Col: 2}}}, // a semicolon was automatically inserted.
		{in: "}\n", want: []token.Token{{Kind: token.Rbrace, Val: "}", Line: 1, Col: 1}, {Kind: token.Semicolon, Val: ";", Line: 1, Col: 2}}}, // a semicolon was automatically inserted.
		{in: ";\n", want: []token.Token{{Kind: token.Semicolon, Val: ";", Line: 1, Col: 1}}},                                                  // a semicolon is present in the source
		{in: ":\n", want: []token.Token{{Kind: token.Colon, Val: ":", Line: 1, Col: 1}}},

		{in: "break\n", want: []token.Token{{Kind: token.Break, Val: "break", Line: 1, Col: 1}, {Kind: token.Semicolon, Val: ";", Line: 1, Col: 6}}}, // a semicolon was automatically inserted.
		{in: "case\n", want: []token.Token{{Kind: token.Case, Val: "case", Line: 1, Col: 1}}},
		{in: "chan\n", want: []token.Token{{Kind: token.Chan, Val: "chan", Line: 1, Col: 1}}},
		{in: "const\n", want: []token.Token{{Kind: token.Const, Val: "const", Line: 1, Col: 1}}},
		{in: "continue\n", want: []token.Token{{Kind: token.Continue, Val: "continue", Line: 1, Col: 1}, {Kind: token.Semicolon, Val: ";", Line: 1, Col: 9}}}, // a semicolon was automatically inserted.

		{in: "default\n", want: []token.Token{{Kind: token.Default, Val: "default", Line: 1, Col: 1}}},
		{in: "defer\n", want: []token.Token{{Kind: token.Defer, Val: "defer", Line: 1, Col: 1}}},
		{in: "else\n", want: []token.Token{{Kind: token.Else, Val: "else", Line: 1, Col: 1}}},
		{in: "fallthrough\n", want: []token.Token{{Kind: token.Fallthrough, Val: "fallthrough", Line: 1, Col: 1}, {Kind: token.Semicolon, Val: ";", Line: 1, Col: 12}}}, // a semicolon was automatically inserted.
		{in: "for\n", want: []token.Token{{Kind: token.For, Val: "for", Line: 1, Col: 1}}},

		{in: "func\n", want: []token.Token{{Kind: token.Func, Val: "func", Line: 1, Col: 1}}},
		{in: "go\n", want: []token.Token{{Kind: token.Go, Val: "go", Line: 1, Col: 1}}},
		{in: "goto\n", want: []token.Token{{Kind: token.Goto, Val: "goto", Line: 1, Col: 1}}},
		{in: "if\n", want: []token.Token{{Kind: token.If, Val: "if", Line: 1, Col: 1}}},
		{in: "import\n", want: []token.Token{{Kind: token.Import, Val: "import", Line: 1, Col: 1}}},

		{in: "interface\n", want: []token.Token{{Kind: token.Interface, Val: "interface", Line: 1, Col: 1}}},
		{in: "map\n", want: []token.Token{{Kind: token.Map, Val: "map", Line: 1, Col: 1}}},
		{in: "package\n", want: []token.Token{{Kind: token.Package, Val: "package", Line: 1, Col: 1}}},
		{in: "range\n", want: []token.Token{{Kind: token.Range, Val: "range", Line: 1, Col: 1}}},
		{in: "return\n", want: []token.Token{{Kind: token.Return, Val: "return", Line: 1, Col: 1}, {Kind: token.Semicolon, Val: ";", Line: 1, Col: 7}}}, // a semicolon was automatically inserted.

		{in: "select\n", want: []token.Token{{Kind: token.Select, Val: "select", Line: 1, Col: 1}}},
		{in: "struct\n", want: []token.Token{{Kind: token.Struct, Val: "struct", Line: 1, Col: 1}}},
		{in: "switch\n", want: []token.Token{{Kind: token.Switch, Val: "switch", Line: 1, Col: 1}}},
		{in: "type\n", want: []token.Token{{Kind: token.Type, Val: "type", Line: 1, Col: 1}}},
		{in: "var\n", want: []token.Token{{Kind: token.Var, Val: "var", Line: 1, Col: 1}}},

		{in: "foo//comment\n", want: []token.Token{{Kind: token.Ident, Val: "foo", Line: 1, Col: 1}, {Kind: token.Semicolon, Val: ";", Line: 1, Col: 4}, {Kind: token.Comment, Val: "//comment", Line: 1, Col: 4}}},         // a semicolon was automatically inserted.
		{in: "foo//comment", want: []token.Token{{Kind: token.Ident, Val: "foo", Line: 1, Col: 1}, {Kind: token.Semicolon, Val: ";", Line: 1, Col: 4}, {Kind: token.Comment, Val: "//comment", Line: 1, Col: 4}}},           // a semicolon was automatically inserted.
		{in: "foo/*comment*/\n", want: []token.Token{{Kind: token.Ident, Val: "foo", Line: 1, Col: 1}, {Kind: token.Semicolon, Val: ";", Line: 1, Col: 4}, {Kind: token.Comment, Val: "/*comment*/", Line: 1, Col: 4}}},     // a semicolon was automatically inserted.
		{in: "foo/*\n*/", want: []token.Token{{Kind: token.Ident, Val: "foo", Line: 1, Col: 1}, {Kind: token.Semicolon, Val: ";", Line: 1, Col: 4}, {Kind: token.Comment, Val: "/*\n*/", Line: 1, Col: 4}}},                 // a semicolon was automatically inserted.
		{in: "foo/*comment*/    \n", want: []token.Token{{Kind: token.Ident, Val: "foo", Line: 1, Col: 1}, {Kind: token.Semicolon, Val: ";", Line: 1, Col: 4}, {Kind: token.Comment, Val: "/*comment*/", Line: 1, Col: 4}}}, // a semicolon was automatically inserted.
		{in: "foo/*\n*/    ", want: []token.Token{{Kind: token.Ident, Val: "foo", Line: 1, Col: 1}, {Kind: token.Semicolon, Val: ";", Line: 1, Col: 4}, {Kind: token.Comment, Val: "/*\n*/", Line: 1, Col: 4}}},             // a semicolon was automatically inserted.

		{in: "foo    // comment\n", want: []token.Token{{Kind: token.Ident, Val: "foo", Line: 1, Col: 1}, {Kind: token.Semicolon, Val: ";", Line: 1, Col: 4}, {Kind: token.Comment, Val: "// comment", Line: 1, Col: 8}}},                                                                                                                                                                                                                               // a semicolon was automatically inserted.
		{in: "foo    // comment", want: []token.Token{{Kind: token.Ident, Val: "foo", Line: 1, Col: 1}, {Kind: token.Semicolon, Val: ";", Line: 1, Col: 4}, {Kind: token.Comment, Val: "// comment", Line: 1, Col: 8}}},                                                                                                                                                                                                                                 // a semicolon was automatically inserted.
		{in: "foo    /*comment*/\n", want: []token.Token{{Kind: token.Ident, Val: "foo", Line: 1, Col: 1}, {Kind: token.Semicolon, Val: ";", Line: 1, Col: 4}, {Kind: token.Comment, Val: "/*comment*/", Line: 1, Col: 8}}},                                                                                                                                                                                                                             // a semicolon was automatically inserted.
		{in: "foo    /*\n*/", want: []token.Token{{Kind: token.Ident, Val: "foo", Line: 1, Col: 1}, {Kind: token.Semicolon, Val: ";", Line: 1, Col: 4}, {Kind: token.Comment, Val: "/*\n*/", Line: 1, Col: 8}}},                                                                                                                                                                                                                                         // a semicolon was automatically inserted.
		{in: "foo    /*  */ /* \n */ bar/**/\n", want: []token.Token{{Kind: token.Ident, Val: "foo", Line: 1, Col: 1}, {Kind: token.Semicolon, Val: ";", Line: 1, Col: 4}, {Kind: token.Comment, Val: "/*  */", Line: 1, Col: 8}, {Kind: token.Comment, Val: "/* \n */", Line: 1, Col: 15}, {Kind: token.Ident, Val: "bar", Line: 2, Col: 5}, {Kind: token.Semicolon, Val: ";", Line: 2, Col: 8}, {Kind: token.Comment, Val: "/**/", Line: 2, Col: 8}}}, // a semicolon was automatically inserted.
		{in: "foo    /*0*/ /*1*/ /*2*/\n", want: []token.Token{{Kind: token.Ident, Val: "foo", Line: 1, Col: 1}, {Kind: token.Semicolon, Val: ";", Line: 1, Col: 4}, {Kind: token.Comment, Val: "/*0*/", Line: 1, Col: 8}, {Kind: token.Comment, Val: "/*1*/", Line: 1, Col: 14}, {Kind: token.Comment, Val: "/*2*/", Line: 1, Col: 20}}},                                                                                                               // a semicolon was automatically inserted.

		{in: "foo    /*comment*/    \n", want: []token.Token{{Kind: token.Ident, Val: "foo", Line: 1, Col: 1}, {Kind: token.Semicolon, Val: ";", Line: 1, Col: 4}, {Kind: token.Comment, Val: "/*comment*/", Line: 1, Col: 8}}},                                                                                                               // a semicolon was automatically inserted.
		{in: "foo    /*0*/ /*1*/ /*2*/    \n", want: []token.Token{{Kind: token.Ident, Val: "foo", Line: 1, Col: 1}, {Kind: token.Semicolon, Val: ";", Line: 1, Col: 4}, {Kind: token.Comment, Val: "/*0*/", Line: 1, Col: 8}, {Kind: token.Comment, Val: "/*1*/", Line: 1, Col: 14}, {Kind: token.Comment, Val: "/*2*/", Line: 1, Col: 20}}}, // a semicolon was automatically inserted.
		{in: "foo	/**/ /*-------------*/       /*----\n*/bar       /*  \n*/baa\n", want: []token.Token{{Kind: token.Ident, Val: "foo", Line: 1, Col: 1}, {Kind: token.Semicolon, Val: ";", Line: 1, Col: 4}, {Kind: token.Comment, Val: "/**/", Line: 1, Col: 5}, {Kind: token.Comment, Val: "/*-------------*/", Line: 1, Col: 10}, {Kind: token.Comment, Val: "/*----\n*/", Line: 1, Col: 34}, {Kind: token.Ident, Val: "bar", Line: 2, Col: 3}, {Kind: token.Semicolon, Val: ";", Line: 2, Col: 6}, {Kind: token.Comment, Val: "/*  \n*/", Line: 2, Col: 13}, {Kind: token.Ident, Val: "baa", Line: 3, Col: 3}, {Kind: token.Semicolon, Val: ";", Line: 3, Col: 6}}}, // a semicolon was automatically inserted.
		{in: "foo    /* an EOF terminates a line */", want: []token.Token{{Kind: token.Ident, Val: "foo", Line: 1, Col: 1}, {Kind: token.Semicolon, Val: ";", Line: 1, Col: 4}, {Kind: token.Comment, Val: "/* an EOF terminates a line */", Line: 1, Col: 8}}},                                                                                                          // a semicolon was automatically inserted.
		{in: "foo    /* an EOF terminates a line */ /*", want: []token.Token{{Kind: token.Ident, Val: "foo", Line: 1, Col: 1}, {Kind: token.Semicolon, Val: ";", Line: 1, Col: 4}, {Kind: token.Comment, Val: "/* an EOF terminates a line */", Line: 1, Col: 8}, {Kind: token.Comment | token.Invalid, Val: "/*", Line: 1, Col: 39}}, err: "unexpected eof in comment"}, // a semicolon was automatically inserted.
		{in: "foo    /* an EOF terminates a line */ //", want: []token.Token{{Kind: token.Ident, Val: "foo", Line: 1, Col: 1}, {Kind: token.Semicolon, Val: ";", Line: 1, Col: 4}, {Kind: token.Comment, Val: "/* an EOF terminates a line */", Line: 1, Col: 8}, {Kind: token.Comment, Val: "//", Line: 1, Col: 39}}},                                                   // a semicolon was automatically inserted.

		{in: "package main\n\nfunc main() {\n\tif {\n\t\treturn /* */ }\n}\n", want: []token.Token{{Kind: token.Package, Val: "package", Line: 1, Col: 1}, {Kind: token.Ident, Val: "main", Line: 1, Col: 9}, {Kind: token.Semicolon, Val: ";", Line: 1, Col: 13}, {Kind: token.Func, Val: "func", Line: 3, Col: 1}, {Kind: token.Ident, Val: "main", Line: 3, Col: 6}, {Kind: token.Lparen, Val: "(", Line: 3, Col: 10}, {Kind: token.Rparen, Val: ")", Line: 3, Col: 11}, {Kind: token.Lbrace, Val: "{", Line: 3, Col: 13}, {Kind: token.If, Val: "if", Line: 4, Col: 2}, {Kind: token.Lbrace, Val: "{", Line: 4, Col: 5}, {Kind: token.Return, Val: "return", Line: 5, Col: 3}, {Kind: token.Comment, Val: "/* */", Line: 5, Col: 10}, {Kind: token.Rbrace, Val: "}", Line: 5, Col: 16}, {Kind: token.Semicolon, Val: ";", Line: 5, Col: 17}, {Kind: token.Rbrace, Val: "}", Line: 6, Col: 1}, {Kind: token.Semicolon, Val: ";", Line: 6, Col: 2}}}, // a semicolon was automatically inserted.
		{in: "package main", want: []token.Token{{Kind: token.Package, Val: "package", Line: 1, Col: 1}, {Kind: token.Ident, Val: "main", Line: 1, Col: 9}, {Kind: token.Semicolon, Val: ";", Line: 1, Col: 13}}}, // a semicolon was automatically inserted.
	}

	for i, g := range golden {
		got, err := Parse(g.in)
		if err != nil && err.Error() != g.err {
			t.Errorf("i=%d: Parse failed; %v", i, err)
			continue
		}
		if !reflect.DeepEqual(got, g.want) {
			t.Errorf("i=%d: token mismatch; expected %#v, got %#v.", i, g.want, got)
		}
	}
}

func TestParseErrors(t *testing.T) {
	// test cases derived from errors in go/src/pkg/scanner/scanner_test.go
	golden := []struct {
		in   string
		err  string
		want token.Token
	}{
		{in: "\a", err: "syntax error: unexpected U+0007", want: token.Token{Kind: token.Invalid, Val: "\a", Line: 1, Col: 1}},
		{in: `#`, err: "syntax error: unexpected U+0023 '#'", want: token.Token{Kind: token.Invalid, Val: `#`, Line: 1, Col: 1}},
		{in: `…`, err: "syntax error: unexpected U+2026 '…'", want: token.Token{Kind: token.Invalid, Val: `…`, Line: 1, Col: 1}},
		{in: `' '`, want: token.Token{Kind: token.Rune, Val: "' '", Line: 1, Col: 1}},
		{in: `''`, err: "empty rune literal or unescaped ' in rune literal", want: token.Token{Kind: token.Rune | token.Invalid, Val: "''", Line: 1, Col: 1}},
		{in: `'12'`, err: "too many characters in rune literal", want: token.Token{Kind: token.Rune | token.Invalid, Val: "'12'", Line: 1, Col: 1}},
		{in: `'123'`, err: "too many characters in rune literal", want: token.Token{Kind: token.Rune | token.Invalid, Val: "'123'", Line: 1, Col: 1}},
		{in: `'\0'`, err: "too few digits in octal escape; expected 3, got 1", want: token.Token{Kind: token.Rune | token.Invalid, Val: `'\0'`, Line: 1, Col: 1}},
		{in: `'\07'`, err: "too few digits in octal escape; expected 3, got 2", want: token.Token{Kind: token.Rune | token.Invalid, Val: `'\07'`, Line: 1, Col: 1}},
		{in: `'\8'`, err: "unknown escape sequence U+0038 '8'", want: token.Token{Kind: token.Rune | token.Invalid, Val: `'\8'`, Line: 1, Col: 1}},
		{in: `'\08'`, err: "non-octal character U+0038 '8' in octal escape", want: token.Token{Kind: token.Rune | token.Invalid, Val: `'\08'`, Line: 1, Col: 1}},
		{in: `'\0`, err: "unexpected eof in octal escape", want: token.Token{Kind: token.Rune | token.Invalid, Val: `'\0`, Line: 1, Col: 1}},
		{in: `'\00`, err: "unexpected eof in octal escape", want: token.Token{Kind: token.Rune | token.Invalid, Val: `'\00`, Line: 1, Col: 1}},
		{in: `'\000`, err: "unexpected eof in rune literal", want: token.Token{Kind: token.Rune | token.Invalid, Val: `'\000`, Line: 1, Col: 1}},
		{in: `'\x'`, err: "too few digits in hex escape; expected 2, got 0", want: token.Token{Kind: token.Rune | token.Invalid, Val: `'\x'`, Line: 1, Col: 1}},
		{in: `'\x0'`, err: "too few digits in hex escape; expected 2, got 1", want: token.Token{Kind: token.Rune | token.Invalid, Val: `'\x0'`, Line: 1, Col: 1}},
		{in: `'\x0g'`, err: "non-hex character U+0067 'g' in hex escape", want: token.Token{Kind: token.Rune | token.Invalid, Val: `'\x0g'`, Line: 1, Col: 1}},
		{in: `'\x`, err: "unexpected eof in hex escape", want: token.Token{Kind: token.Rune | token.Invalid, Val: `'\x`, Line: 1, Col: 1}},
		{in: `'\x0`, err: "unexpected eof in hex escape", want: token.Token{Kind: token.Rune | token.Invalid, Val: `'\x0`, Line: 1, Col: 1}},
		{in: `'\x00`, err: "unexpected eof in rune literal", want: token.Token{Kind: token.Rune | token.Invalid, Val: `'\x00`, Line: 1, Col: 1}},
		{in: `'\u'`, err: "too few digits in Unicode escape; expected 4, got 0", want: token.Token{Kind: token.Rune | token.Invalid, Val: `'\u'`, Line: 1, Col: 1}},
		{in: `'\u0'`, err: "too few digits in Unicode escape; expected 4, got 1", want: token.Token{Kind: token.Rune | token.Invalid, Val: `'\u0'`, Line: 1, Col: 1}},
		{in: `'\u00'`, err: "too few digits in Unicode escape; expected 4, got 2", want: token.Token{Kind: token.Rune | token.Invalid, Val: `'\u00'`, Line: 1, Col: 1}},
		{in: `'\u000'`, err: "too few digits in Unicode escape; expected 4, got 3", want: token.Token{Kind: token.Rune | token.Invalid, Val: `'\u000'`, Line: 1, Col: 1}},
		{in: `'\u000`, err: "unexpected eof in Unicode escape", want: token.Token{Kind: token.Rune | token.Invalid, Val: `'\u000`, Line: 1, Col: 1}},
		{in: `'\u0000'`, want: token.Token{Kind: token.Rune, Val: `'\u0000'`, Line: 1, Col: 1}},
		{in: `'\U'`, err: "too few digits in Unicode escape; expected 8, got 0", want: token.Token{Kind: token.Rune | token.Invalid, Val: `'\U'`, Line: 1, Col: 1}},
		{in: `'\U0'`, err: "too few digits in Unicode escape; expected 8, got 1", want: token.Token{Kind: token.Rune | token.Invalid, Val: `'\U0'`, Line: 1, Col: 1}},
		{in: `'\U00'`, err: "too few digits in Unicode escape; expected 8, got 2", want: token.Token{Kind: token.Rune | token.Invalid, Val: `'\U00'`, Line: 1, Col: 1}},
		{in: `'\U000'`, err: "too few digits in Unicode escape; expected 8, got 3", want: token.Token{Kind: token.Rune | token.Invalid, Val: `'\U000'`, Line: 1, Col: 1}},
		{in: `'\U0000'`, err: "too few digits in Unicode escape; expected 8, got 4", want: token.Token{Kind: token.Rune | token.Invalid, Val: `'\U0000'`, Line: 1, Col: 1}},
		{in: `'\U00000'`, err: "too few digits in Unicode escape; expected 8, got 5", want: token.Token{Kind: token.Rune | token.Invalid, Val: `'\U00000'`, Line: 1, Col: 1}},
		{in: `'\U000000'`, err: "too few digits in Unicode escape; expected 8, got 6", want: token.Token{Kind: token.Rune | token.Invalid, Val: `'\U000000'`, Line: 1, Col: 1}},
		{in: `'\U0000000'`, err: "too few digits in Unicode escape; expected 8, got 7", want: token.Token{Kind: token.Rune | token.Invalid, Val: `'\U0000000'`, Line: 1, Col: 1}},
		{in: `'\U0000000`, err: "unexpected eof in Unicode escape", want: token.Token{Kind: token.Rune | token.Invalid, Val: `'\U0000000`, Line: 1, Col: 1}},
		{in: `'\U00000000'`, want: token.Token{Kind: token.Rune, Val: `'\U00000000'`, Line: 1, Col: 1}},
		{in: `'\Uffffffff'`, err: "invalid Unicode code point U+FFFFFFFFFFFFFFFF in escape sequence", want: token.Token{Kind: token.Rune | token.Invalid, Val: `'\Uffffffff'`, Line: 1, Col: 1}},
		{in: `'\U0g'`, err: "non-hex character U+0067 'g' in Unicode escape", want: token.Token{Kind: token.Rune | token.Invalid, Val: `'\U0g'`, Line: 1, Col: 1}},
		{in: `'`, err: "unexpected eof in rune literal", want: token.Token{Kind: token.Rune | token.Invalid, Val: `'`, Line: 1, Col: 1}},
		{in: `'\`, err: "unexpected eof in escape sequence", want: token.Token{Kind: token.Rune | token.Invalid, Val: `'\`, Line: 1, Col: 1}},
		{in: "'\n", err: "unexpected newline in rune literal", want: token.Token{Kind: token.Rune | token.Invalid, Val: "'", Line: 1, Col: 1}},
		{in: "'\n ", err: "unexpected newline in rune literal", want: token.Token{Kind: token.Rune | token.Invalid, Val: "'", Line: 1, Col: 1}},
		{in: "'x", err: "unexpected eof in rune literal", want: token.Token{Kind: token.Rune | token.Invalid, Val: "'x", Line: 1, Col: 1}},
		{in: "'x\n", err: "unexpected newline in rune literal", want: token.Token{Kind: token.Rune | token.Invalid, Val: "'x", Line: 1, Col: 1}},
		{in: `""`, want: token.Token{Kind: token.String, Val: `""`, Line: 1, Col: 1}},
		{in: `"abc`, err: "unexpected eof in string literal", want: token.Token{Kind: token.String | token.Invalid, Val: `"abc`, Line: 1, Col: 1}},
		{in: "\"abc\n", err: "unexpected newline in string literal", want: token.Token{Kind: token.String | token.Invalid, Val: "\"abc", Line: 1, Col: 1}},
		{in: "\"abc\n ", err: "unexpected newline in string literal", want: token.Token{Kind: token.String | token.Invalid, Val: "\"abc", Line: 1, Col: 1}},
		{in: `"\q"`, err: "unknown escape sequence U+0071 'q'", want: token.Token{Kind: token.String | token.Invalid, Val: `"\q"`, Line: 1, Col: 1}},
		{in: "``", want: token.Token{Kind: token.String, Val: "``", Line: 1, Col: 1}},
		{in: "`", err: "unexpected eof in raw string literal", want: token.Token{Kind: token.String | token.Invalid, Val: "`", Line: 1, Col: 1}},
		{in: "/**/", want: token.Token{Kind: token.Comment, Val: "/**/", Line: 1, Col: 1}},
		{in: "/*", err: "unexpected eof in comment", want: token.Token{Kind: token.Comment | token.Invalid, Val: "/*", Line: 1, Col: 1}},
		{in: "077", want: token.Token{Kind: token.Int, Val: "077", Line: 1, Col: 1}},
		{in: "078.", want: token.Token{Kind: token.Float, Val: "078.", Line: 1, Col: 1}},
		{in: "07801234567.", want: token.Token{Kind: token.Float, Val: "07801234567.", Line: 1, Col: 1}},
		{in: "078e0", want: token.Token{Kind: token.Float, Val: "078e0", Line: 1, Col: 1}},
		{in: "078", err: "invalid digit '8' in octal constant", want: token.Token{Kind: token.Int | token.Invalid, Val: "078", Line: 1, Col: 1}},
		{in: "07800000009", err: "invalid digit '8' in octal constant", want: token.Token{Kind: token.Int | token.Invalid, Val: "07800000009", Line: 1, Col: 1}},
		{in: "079", err: "invalid digit '9' in octal constant", want: token.Token{Kind: token.Int | token.Invalid, Val: "079", Line: 1, Col: 1}},
		{in: "0x", err: "missing digits in hexadecimal constant", want: token.Token{Kind: token.Int | token.Invalid, Val: "0x", Line: 1, Col: 1}},
		{in: "0X", err: "missing digits in hexadecimal constant", want: token.Token{Kind: token.Int | token.Invalid, Val: "0X", Line: 1, Col: 1}},
		{in: ".3e", err: "missing digits in floating-point exponent", want: token.Token{Kind: token.Float | token.Invalid, Val: ".3e", Line: 1, Col: 1}},
		{in: "3.14E", err: "missing digits in floating-point exponent", want: token.Token{Kind: token.Float | token.Invalid, Val: "3.14E", Line: 1, Col: 1}},
		{in: "5e", err: "missing digits in floating-point exponent", want: token.Token{Kind: token.Float | token.Invalid, Val: "5e", Line: 1, Col: 1}},
		{in: "//abc\x00def", err: "illegal NUL character", want: token.Token{Kind: token.Comment | token.Invalid, Val: "//abc\x00def", Line: 1, Col: 1}},
		{in: "/*abc\x00def*/", err: "illegal NUL character", want: token.Token{Kind: token.Comment | token.Invalid, Val: "/*abc\x00def*/", Line: 1, Col: 1}},
		{in: "'\x00'", err: "illegal NUL character", want: token.Token{Kind: token.Rune | token.Invalid, Val: "'\x00'", Line: 1, Col: 1}},
		{in: "\"abc\x00def\"", err: "illegal NUL character", want: token.Token{Kind: token.String | token.Invalid, Val: "\"abc\x00def\"", Line: 1, Col: 1}},
		{in: "`abc\x00def`", err: "illegal NUL character", want: token.Token{Kind: token.String | token.Invalid, Val: "`abc\x00def`", Line: 1, Col: 1}},
		{in: "//abc\x80def", err: "illegal UTF-8 encoding", want: token.Token{Kind: token.Comment | token.Invalid, Val: "//abc\x80def", Line: 1, Col: 1}},
		{in: "/*abc\x80def*/", err: "illegal UTF-8 encoding", want: token.Token{Kind: token.Comment | token.Invalid, Val: "/*abc\x80def*/", Line: 1, Col: 1}},
		{in: "'\x80'", err: "illegal UTF-8 encoding", want: token.Token{Kind: token.Rune | token.Invalid, Val: "'\x80'", Line: 1, Col: 1}},
		{in: "\"abc\x80def\"", err: "illegal UTF-8 encoding", want: token.Token{Kind: token.String | token.Invalid, Val: "\"abc\x80def\"", Line: 1, Col: 1}},
		{in: "`abc\x80def`", err: "illegal UTF-8 encoding", want: token.Token{Kind: token.String | token.Invalid, Val: "`abc\x80def`", Line: 1, Col: 1}},
		{in: "\ufeff\ufeff", err: "illegal byte order mark", want: token.Token{Kind: token.Invalid, Val: "\ufeff", Line: 1, Col: 1}},                               // only first BOM is ignored.
		{in: "//abc\ufeffdef", err: "illegal byte order mark", want: token.Token{Kind: token.Comment | token.Invalid, Val: "//abc\ufeffdef", Line: 1, Col: 1}},     // only first BOM is ignored.
		{in: "/*abc\ufeffdef*/", err: "illegal byte order mark", want: token.Token{Kind: token.Comment | token.Invalid, Val: "/*abc\ufeffdef*/", Line: 1, Col: 1}}, // only first BOM is ignored.
		{in: "'\ufeff'", err: "illegal byte order mark", want: token.Token{Kind: token.Rune | token.Invalid, Val: "'\ufeff'", Line: 1, Col: 1}},                    // only first BOM is ignored.
		{in: "\"abc\ufeffdef\"", err: "illegal byte order mark", want: token.Token{Kind: token.String | token.Invalid, Val: "\"abc\ufeffdef\"", Line: 1, Col: 1}},  // only first BOM is ignored.
		{in: "`abc\ufeffdef`", err: "illegal byte order mark", want: token.Token{Kind: token.String | token.Invalid, Val: "`abc\ufeffdef`", Line: 1, Col: 1}},      // only first BOM is ignored.
	}

	for i, g := range golden {
		tokens, err := Parse(g.in)
		errstr := ""
		if err != nil {
			errstr = err.Error()
		}
		if g.err != errstr {
			t.Errorf("i=%d: error mismatch; expected %v, got %v.", i, g.err, errstr)
		}
		if len(tokens) < 1 {
			t.Errorf("i=%d: too few tokens; expected >= 1, got %d.", i, len(tokens))
			continue
		}
		got := tokens[0]
		if got != g.want {
			t.Errorf("i=%d: token mismatch; expected %#v, got %#v.", i, g.want, got)
		}
	}
}

func TestParsePosition(t *testing.T) {
	input := `// Package p implements …
package p

import "strings"

// T is a bitfield which specifies …
type T uint16

// T bitfield values.
const (
	FooA T = 1<<iota /* bitfield … */ + 0x10   /* Foo start value */
	FooB                                       /* FooB specifies … */
	FooC                                       /* FooC specifies … */
	BarA T = 1<<iota /* bitfield … */ + 0x100  /* Bar start value */
	BarB                                       /* BarB specifies … */
	BarC                                       /* BarC specifies … */
	BazA T = 1<<iota /* bitfield … */ + 0x1000 /* Baz start value */
	BazB                                       /* BazB specifies … */
	BazC                                       /* BazC specifies … */
)

// names specifies the name of each …
var names = map[T]string{
	FooA: "foo A",
	FooB: "foo B",
	FooC: "foo C",
	BarA: "bar A",
	BarB: "bar B",
	BarC: "bar C",
	BazA: "baz A",
	BazB: "baz B",
	BazC: "baz C",
}

func (t T) String() string {
	var ss []string
	for i := uint(0); i < 16; i++ {
		mask := T(1 << i)
		if v := t & mask; v != 0 {
			if s, ok := names[v]; ok {
				ss = append(ss, s)
			}
		}
	}
	return strings.Join(ss, " ")
}

// Merge merges … into a single T.
func Merge(ts ...T) T {
	var t T
	for i := range ts {
		t |= ts[i]
	}
	return t
}
`
	want := []token.Token{
		{Kind: token.Comment, Val: "// Package p implements …", Line: 1, Col: 1},
		{Kind: token.Package, Val: "package", Line: 2, Col: 1},
		{Kind: token.Ident, Val: "p", Line: 2, Col: 9},
		{Kind: token.Semicolon, Val: ";", Line: 2, Col: 10},
		{Kind: token.Import, Val: "import", Line: 4, Col: 1},
		{Kind: token.String, Val: "\"strings\"", Line: 4, Col: 8},
		{Kind: token.Semicolon, Val: ";", Line: 4, Col: 17},
		{Kind: token.Comment, Val: "// T is a bitfield which specifies …", Line: 6, Col: 1},
		{Kind: token.Type, Val: "type", Line: 7, Col: 1},
		{Kind: token.Ident, Val: "T", Line: 7, Col: 6},
		{Kind: token.Ident, Val: "uint16", Line: 7, Col: 8},
		{Kind: token.Semicolon, Val: ";", Line: 7, Col: 14},
		{Kind: token.Comment, Val: "// T bitfield values.", Line: 9, Col: 1},
		{Kind: token.Const, Val: "const", Line: 10, Col: 1},
		{Kind: token.Lparen, Val: "(", Line: 10, Col: 7},
		{Kind: token.Ident, Val: "FooA", Line: 11, Col: 2},
		{Kind: token.Ident, Val: "T", Line: 11, Col: 7},
		{Kind: token.Assign, Val: "=", Line: 11, Col: 9},
		{Kind: token.Int, Val: "1", Line: 11, Col: 11},
		{Kind: token.Shl, Val: "<<", Line: 11, Col: 12},
		{Kind: token.Ident, Val: "iota", Line: 11, Col: 14},
		{Kind: token.Comment, Val: "/* bitfield … */", Line: 11, Col: 19},
		{Kind: token.Add, Val: "+", Line: 11, Col: 36},
		{Kind: token.Int, Val: "0x10", Line: 11, Col: 38},
		{Kind: token.Semicolon, Val: ";", Line: 11, Col: 42},
		{Kind: token.Comment, Val: "/* Foo start value */", Line: 11, Col: 45},
		{Kind: token.Ident, Val: "FooB", Line: 12, Col: 2},
		{Kind: token.Semicolon, Val: ";", Line: 12, Col: 6},
		{Kind: token.Comment, Val: "/* FooB specifies … */", Line: 12, Col: 45},
		{Kind: token.Ident, Val: "FooC", Line: 13, Col: 2},
		{Kind: token.Semicolon, Val: ";", Line: 13, Col: 6},
		{Kind: token.Comment, Val: "/* FooC specifies … */", Line: 13, Col: 45},
		{Kind: token.Ident, Val: "BarA", Line: 14, Col: 2},
		{Kind: token.Ident, Val: "T", Line: 14, Col: 7},
		{Kind: token.Assign, Val: "=", Line: 14, Col: 9},
		{Kind: token.Int, Val: "1", Line: 14, Col: 11},
		{Kind: token.Shl, Val: "<<", Line: 14, Col: 12},
		{Kind: token.Ident, Val: "iota", Line: 14, Col: 14},
		{Kind: token.Comment, Val: "/* bitfield … */", Line: 14, Col: 19},
		{Kind: token.Add, Val: "+", Line: 14, Col: 36},
		{Kind: token.Int, Val: "0x100", Line: 14, Col: 38},
		{Kind: token.Semicolon, Val: ";", Line: 14, Col: 43},
		{Kind: token.Comment, Val: "/* Bar start value */", Line: 14, Col: 45},
		{Kind: token.Ident, Val: "BarB", Line: 15, Col: 2},
		{Kind: token.Semicolon, Val: ";", Line: 15, Col: 6},
		{Kind: token.Comment, Val: "/* BarB specifies … */", Line: 15, Col: 45},
		{Kind: token.Ident, Val: "BarC", Line: 16, Col: 2},
		{Kind: token.Semicolon, Val: ";", Line: 16, Col: 6},
		{Kind: token.Comment, Val: "/* BarC specifies … */", Line: 16, Col: 45},
		{Kind: token.Ident, Val: "BazA", Line: 17, Col: 2},
		{Kind: token.Ident, Val: "T", Line: 17, Col: 7},
		{Kind: token.Assign, Val: "=", Line: 17, Col: 9},
		{Kind: token.Int, Val: "1", Line: 17, Col: 11},
		{Kind: token.Shl, Val: "<<", Line: 17, Col: 12},
		{Kind: token.Ident, Val: "iota", Line: 17, Col: 14},
		{Kind: token.Comment, Val: "/* bitfield … */", Line: 17, Col: 19},
		{Kind: token.Add, Val: "+", Line: 17, Col: 36},
		{Kind: token.Int, Val: "0x1000", Line: 17, Col: 38},
		{Kind: token.Semicolon, Val: ";", Line: 17, Col: 44},
		{Kind: token.Comment, Val: "/* Baz start value */", Line: 17, Col: 45},
		{Kind: token.Ident, Val: "BazB", Line: 18, Col: 2},
		{Kind: token.Semicolon, Val: ";", Line: 18, Col: 6},
		{Kind: token.Comment, Val: "/* BazB specifies … */", Line: 18, Col: 45},
		{Kind: token.Ident, Val: "BazC", Line: 19, Col: 2},
		{Kind: token.Semicolon, Val: ";", Line: 19, Col: 6},
		{Kind: token.Comment, Val: "/* BazC specifies … */", Line: 19, Col: 45},
		{Kind: token.Rparen, Val: ")", Line: 20, Col: 1},
		{Kind: token.Semicolon, Val: ";", Line: 20, Col: 2},
		{Kind: token.Comment, Val: "// names specifies the name of each …", Line: 22, Col: 1},
		{Kind: token.Var, Val: "var", Line: 23, Col: 1},
		{Kind: token.Ident, Val: "names", Line: 23, Col: 5},
		{Kind: token.Assign, Val: "=", Line: 23, Col: 11},
		{Kind: token.Map, Val: "map", Line: 23, Col: 13},
		{Kind: token.Lbrack, Val: "[", Line: 23, Col: 16},
		{Kind: token.Ident, Val: "T", Line: 23, Col: 17},
		{Kind: token.Rbrack, Val: "]", Line: 23, Col: 18},
		{Kind: token.Ident, Val: "string", Line: 23, Col: 19},
		{Kind: token.Lbrace, Val: "{", Line: 23, Col: 25},
		{Kind: token.Ident, Val: "FooA", Line: 24, Col: 2},
		{Kind: token.Colon, Val: ":", Line: 24, Col: 6},
		{Kind: token.String, Val: "\"foo A\"", Line: 24, Col: 8},
		{Kind: token.Comma, Val: ",", Line: 24, Col: 15},
		{Kind: token.Ident, Val: "FooB", Line: 25, Col: 2},
		{Kind: token.Colon, Val: ":", Line: 25, Col: 6},
		{Kind: token.String, Val: "\"foo B\"", Line: 25, Col: 8},
		{Kind: token.Comma, Val: ",", Line: 25, Col: 15},
		{Kind: token.Ident, Val: "FooC", Line: 26, Col: 2},
		{Kind: token.Colon, Val: ":", Line: 26, Col: 6},
		{Kind: token.String, Val: "\"foo C\"", Line: 26, Col: 8},
		{Kind: token.Comma, Val: ",", Line: 26, Col: 15},
		{Kind: token.Ident, Val: "BarA", Line: 27, Col: 2},
		{Kind: token.Colon, Val: ":", Line: 27, Col: 6},
		{Kind: token.String, Val: "\"bar A\"", Line: 27, Col: 8},
		{Kind: token.Comma, Val: ",", Line: 27, Col: 15},
		{Kind: token.Ident, Val: "BarB", Line: 28, Col: 2},
		{Kind: token.Colon, Val: ":", Line: 28, Col: 6},
		{Kind: token.String, Val: "\"bar B\"", Line: 28, Col: 8},
		{Kind: token.Comma, Val: ",", Line: 28, Col: 15},
		{Kind: token.Ident, Val: "BarC", Line: 29, Col: 2},
		{Kind: token.Colon, Val: ":", Line: 29, Col: 6},
		{Kind: token.String, Val: "\"bar C\"", Line: 29, Col: 8},
		{Kind: token.Comma, Val: ",", Line: 29, Col: 15},
		{Kind: token.Ident, Val: "BazA", Line: 30, Col: 2},
		{Kind: token.Colon, Val: ":", Line: 30, Col: 6},
		{Kind: token.String, Val: "\"baz A\"", Line: 30, Col: 8},
		{Kind: token.Comma, Val: ",", Line: 30, Col: 15},
		{Kind: token.Ident, Val: "BazB", Line: 31, Col: 2},
		{Kind: token.Colon, Val: ":", Line: 31, Col: 6},
		{Kind: token.String, Val: "\"baz B\"", Line: 31, Col: 8},
		{Kind: token.Comma, Val: ",", Line: 31, Col: 15},
		{Kind: token.Ident, Val: "BazC", Line: 32, Col: 2},
		{Kind: token.Colon, Val: ":", Line: 32, Col: 6},
		{Kind: token.String, Val: "\"baz C\"", Line: 32, Col: 8},
		{Kind: token.Comma, Val: ",", Line: 32, Col: 15},
		{Kind: token.Rbrace, Val: "}", Line: 33, Col: 1},
		{Kind: token.Semicolon, Val: ";", Line: 33, Col: 2},
		{Kind: token.Func, Val: "func", Line: 35, Col: 1},
		{Kind: token.Lparen, Val: "(", Line: 35, Col: 6},
		{Kind: token.Ident, Val: "t", Line: 35, Col: 7},
		{Kind: token.Ident, Val: "T", Line: 35, Col: 9},
		{Kind: token.Rparen, Val: ")", Line: 35, Col: 10},
		{Kind: token.Ident, Val: "String", Line: 35, Col: 12},
		{Kind: token.Lparen, Val: "(", Line: 35, Col: 18},
		{Kind: token.Rparen, Val: ")", Line: 35, Col: 19},
		{Kind: token.Ident, Val: "string", Line: 35, Col: 21},
		{Kind: token.Lbrace, Val: "{", Line: 35, Col: 28},
		{Kind: token.Var, Val: "var", Line: 36, Col: 2},
		{Kind: token.Ident, Val: "ss", Line: 36, Col: 6},
		{Kind: token.Lbrack, Val: "[", Line: 36, Col: 9},
		{Kind: token.Rbrack, Val: "]", Line: 36, Col: 10},
		{Kind: token.Ident, Val: "string", Line: 36, Col: 11},
		{Kind: token.Semicolon, Val: ";", Line: 36, Col: 17},
		{Kind: token.For, Val: "for", Line: 37, Col: 2},
		{Kind: token.Ident, Val: "i", Line: 37, Col: 6},
		{Kind: token.DeclAssign, Val: ":=", Line: 37, Col: 8},
		{Kind: token.Ident, Val: "uint", Line: 37, Col: 11},
		{Kind: token.Lparen, Val: "(", Line: 37, Col: 15},
		{Kind: token.Int, Val: "0", Line: 37, Col: 16},
		{Kind: token.Rparen, Val: ")", Line: 37, Col: 17},
		{Kind: token.Semicolon, Val: ";", Line: 37, Col: 18},
		{Kind: token.Ident, Val: "i", Line: 37, Col: 20},
		{Kind: token.Lt, Val: "<", Line: 37, Col: 22},
		{Kind: token.Int, Val: "16", Line: 37, Col: 24},
		{Kind: token.Semicolon, Val: ";", Line: 37, Col: 26},
		{Kind: token.Ident, Val: "i", Line: 37, Col: 28},
		{Kind: token.Inc, Val: "++", Line: 37, Col: 29},
		{Kind: token.Lbrace, Val: "{", Line: 37, Col: 32},
		{Kind: token.Ident, Val: "mask", Line: 38, Col: 3},
		{Kind: token.DeclAssign, Val: ":=", Line: 38, Col: 8},
		{Kind: token.Ident, Val: "T", Line: 38, Col: 11},
		{Kind: token.Lparen, Val: "(", Line: 38, Col: 12},
		{Kind: token.Int, Val: "1", Line: 38, Col: 13},
		{Kind: token.Shl, Val: "<<", Line: 38, Col: 15},
		{Kind: token.Ident, Val: "i", Line: 38, Col: 18},
		{Kind: token.Rparen, Val: ")", Line: 38, Col: 19},
		{Kind: token.Semicolon, Val: ";", Line: 38, Col: 20},
		{Kind: token.If, Val: "if", Line: 39, Col: 3},
		{Kind: token.Ident, Val: "v", Line: 39, Col: 6},
		{Kind: token.DeclAssign, Val: ":=", Line: 39, Col: 8},
		{Kind: token.Ident, Val: "t", Line: 39, Col: 11},
		{Kind: token.And, Val: "&", Line: 39, Col: 13},
		{Kind: token.Ident, Val: "mask", Line: 39, Col: 15},
		{Kind: token.Semicolon, Val: ";", Line: 39, Col: 19},
		{Kind: token.Ident, Val: "v", Line: 39, Col: 21},
		{Kind: token.Neq, Val: "!=", Line: 39, Col: 23},
		{Kind: token.Int, Val: "0", Line: 39, Col: 26},
		{Kind: token.Lbrace, Val: "{", Line: 39, Col: 28},
		{Kind: token.If, Val: "if", Line: 40, Col: 4},
		{Kind: token.Ident, Val: "s", Line: 40, Col: 7},
		{Kind: token.Comma, Val: ",", Line: 40, Col: 8},
		{Kind: token.Ident, Val: "ok", Line: 40, Col: 10},
		{Kind: token.DeclAssign, Val: ":=", Line: 40, Col: 13},
		{Kind: token.Ident, Val: "names", Line: 40, Col: 16},
		{Kind: token.Lbrack, Val: "[", Line: 40, Col: 21},
		{Kind: token.Ident, Val: "v", Line: 40, Col: 22},
		{Kind: token.Rbrack, Val: "]", Line: 40, Col: 23},
		{Kind: token.Semicolon, Val: ";", Line: 40, Col: 24},
		{Kind: token.Ident, Val: "ok", Line: 40, Col: 26},
		{Kind: token.Lbrace, Val: "{", Line: 40, Col: 29},
		{Kind: token.Ident, Val: "ss", Line: 41, Col: 5},
		{Kind: token.Assign, Val: "=", Line: 41, Col: 8},
		{Kind: token.Ident, Val: "append", Line: 41, Col: 10},
		{Kind: token.Lparen, Val: "(", Line: 41, Col: 16},
		{Kind: token.Ident, Val: "ss", Line: 41, Col: 17},
		{Kind: token.Comma, Val: ",", Line: 41, Col: 19},
		{Kind: token.Ident, Val: "s", Line: 41, Col: 21},
		{Kind: token.Rparen, Val: ")", Line: 41, Col: 22},
		{Kind: token.Semicolon, Val: ";", Line: 41, Col: 23},
		{Kind: token.Rbrace, Val: "}", Line: 42, Col: 4},
		{Kind: token.Semicolon, Val: ";", Line: 42, Col: 5},
		{Kind: token.Rbrace, Val: "}", Line: 43, Col: 3},
		{Kind: token.Semicolon, Val: ";", Line: 43, Col: 4},
		{Kind: token.Rbrace, Val: "}", Line: 44, Col: 2},
		{Kind: token.Semicolon, Val: ";", Line: 44, Col: 3},
		{Kind: token.Return, Val: "return", Line: 45, Col: 2},
		{Kind: token.Ident, Val: "strings", Line: 45, Col: 9},
		{Kind: token.Dot, Val: ".", Line: 45, Col: 16},
		{Kind: token.Ident, Val: "Join", Line: 45, Col: 17},
		{Kind: token.Lparen, Val: "(", Line: 45, Col: 21},
		{Kind: token.Ident, Val: "ss", Line: 45, Col: 22},
		{Kind: token.Comma, Val: ",", Line: 45, Col: 24},
		{Kind: token.String, Val: "\" \"", Line: 45, Col: 26},
		{Kind: token.Rparen, Val: ")", Line: 45, Col: 29},
		{Kind: token.Semicolon, Val: ";", Line: 45, Col: 30},
		{Kind: token.Rbrace, Val: "}", Line: 46, Col: 1},
		{Kind: token.Semicolon, Val: ";", Line: 46, Col: 2},
		{Kind: token.Comment, Val: "// Merge merges … into a single T.", Line: 48, Col: 1},
		{Kind: token.Func, Val: "func", Line: 49, Col: 1},
		{Kind: token.Ident, Val: "Merge", Line: 49, Col: 6},
		{Kind: token.Lparen, Val: "(", Line: 49, Col: 11},
		{Kind: token.Ident, Val: "ts", Line: 49, Col: 12},
		{Kind: token.Ellipsis, Val: "...", Line: 49, Col: 15},
		{Kind: token.Ident, Val: "T", Line: 49, Col: 18},
		{Kind: token.Rparen, Val: ")", Line: 49, Col: 19},
		{Kind: token.Ident, Val: "T", Line: 49, Col: 21},
		{Kind: token.Lbrace, Val: "{", Line: 49, Col: 23},
		{Kind: token.Var, Val: "var", Line: 50, Col: 2},
		{Kind: token.Ident, Val: "t", Line: 50, Col: 6},
		{Kind: token.Ident, Val: "T", Line: 50, Col: 8},
		{Kind: token.Semicolon, Val: ";", Line: 50, Col: 9},
		{Kind: token.For, Val: "for", Line: 51, Col: 2},
		{Kind: token.Ident, Val: "i", Line: 51, Col: 6},
		{Kind: token.DeclAssign, Val: ":=", Line: 51, Col: 8},
		{Kind: token.Range, Val: "range", Line: 51, Col: 11},
		{Kind: token.Ident, Val: "ts", Line: 51, Col: 17},
		{Kind: token.Lbrace, Val: "{", Line: 51, Col: 20},
		{Kind: token.Ident, Val: "t", Line: 52, Col: 3},
		{Kind: token.OrAssign, Val: "|=", Line: 52, Col: 5},
		{Kind: token.Ident, Val: "ts", Line: 52, Col: 8},
		{Kind: token.Lbrack, Val: "[", Line: 52, Col: 10},
		{Kind: token.Ident, Val: "i", Line: 52, Col: 11},
		{Kind: token.Rbrack, Val: "]", Line: 52, Col: 12},
		{Kind: token.Semicolon, Val: ";", Line: 52, Col: 13},
		{Kind: token.Rbrace, Val: "}", Line: 53, Col: 2},
		{Kind: token.Semicolon, Val: ";", Line: 53, Col: 3},
		{Kind: token.Return, Val: "return", Line: 54, Col: 2},
		{Kind: token.Ident, Val: "t", Line: 54, Col: 9},
		{Kind: token.Semicolon, Val: ";", Line: 54, Col: 10},
		{Kind: token.Rbrace, Val: "}", Line: 55, Col: 1},
		{Kind: token.Semicolon, Val: ";", Line: 55, Col: 2},
	}
	got, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error; %v", err)
	}
	for i := range want {
		if i >= len(got) {
			t.Fatalf("too few tokens; expected >= %d, got %d.", len(want), len(got))
			continue
		}
		if got[i] != want[i] {
			t.Errorf("i=%d: token mismatch; expected %#v, got %#v.", i, want[i], got[i])
		}
	}
}

func BenchmarkParse(b *testing.B) {
	b.SetBytes(int64(len(source)))
	for i := 0; i < b.N; i++ {
		Parse(source)
	}
}
