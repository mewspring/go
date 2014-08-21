package token

import "testing"

type test struct {
	kind Kind
	want bool
}

func TestKindIsKeyword(t *testing.T) {
	golden := []test{
		// Keywords.
		{kind: Break, want: true},
		{kind: Case, want: true},
		{kind: Chan, want: true},
		{kind: Const, want: true},
		{kind: Continue, want: true},
		{kind: Default, want: true},
		{kind: Defer, want: true},
		{kind: Else, want: true},
		{kind: Fallthrough, want: true},
		{kind: For, want: true},
		{kind: Func, want: true},
		{kind: Go, want: true},
		{kind: Goto, want: true},
		{kind: If, want: true},
		{kind: Import, want: true},
		{kind: Interface, want: true},
		{kind: Map, want: true},
		{kind: Package, want: true},
		{kind: Range, want: true},
		{kind: Return, want: true},
		{kind: Select, want: true},
		{kind: Struct, want: true},
		{kind: Switch, want: true},
		{kind: Type, want: true},
		{kind: Var, want: true},

		// Other tokens.
		{kind: Add, want: false},
		{kind: AddAssign, want: false},
		{kind: And, want: false},
		{kind: AndAssign, want: false},
		{kind: Arrow, want: false},
		{kind: Assign, want: false},
		{kind: Clear, want: false},
		{kind: ClearAssign, want: false},
		{kind: Colon, want: false},
		{kind: Comma, want: false},
		{kind: Comment, want: false},
		{kind: Dec, want: false},
		{kind: DeclAssign, want: false},
		{kind: Div, want: false},
		{kind: DivAssign, want: false},
		{kind: Dot, want: false},
		{kind: Ellipsis, want: false},
		{kind: EOF, want: false},
		{kind: Eq, want: false},
		{kind: Float, want: false},
		{kind: Gt, want: false},
		{kind: Gte, want: false},
		{kind: Ident, want: false},
		{kind: Imag, want: false},
		{kind: Inc, want: false},
		{kind: Int, want: false},
		{kind: Land, want: false},
		{kind: Lbrace, want: false},
		{kind: Lbrack, want: false},
		{kind: Lor, want: false},
		{kind: Lparen, want: false},
		{kind: Lt, want: false},
		{kind: Lte, want: false},
		{kind: Mod, want: false},
		{kind: ModAssign, want: false},
		{kind: Mul, want: false},
		{kind: MulAssign, want: false},
		{kind: Neq, want: false},
		{kind: Not, want: false},
		{kind: Or, want: false},
		{kind: OrAssign, want: false},
		{kind: Rbrace, want: false},
		{kind: Rbrack, want: false},
		{kind: Rparen, want: false},
		{kind: Rune, want: false},
		{kind: Semicolon, want: false},
		{kind: Shl, want: false},
		{kind: ShlAssign, want: false},
		{kind: Shr, want: false},
		{kind: ShrAssign, want: false},
		{kind: String, want: false},
		{kind: Sub, want: false},
		{kind: SubAssign, want: false},
		{kind: Xor, want: false},
		{kind: XorAssign, want: false},
	}

	for i, g := range golden {
		got := g.kind.IsKeyword()
		if got != g.want {
			t.Errorf("i=%d: IsKeyword mismatch for token type %v; expected %t, got %t.", i, g.kind, g.want, got)
		}
	}
}

func TestKindIsOperator(t *testing.T) {
	golden := []test{
		// Operators and delimiters.
		{kind: Add, want: true},
		{kind: AddAssign, want: true},
		{kind: And, want: true},
		{kind: AndAssign, want: true},
		{kind: Arrow, want: true},
		{kind: Assign, want: true},
		{kind: Clear, want: true},
		{kind: ClearAssign, want: true},
		{kind: Colon, want: true},
		{kind: Comma, want: true},
		{kind: Dec, want: true},
		{kind: DeclAssign, want: true},
		{kind: Div, want: true},
		{kind: DivAssign, want: true},
		{kind: Dot, want: true},
		{kind: Ellipsis, want: true},
		{kind: Eq, want: true},
		{kind: Gt, want: true},
		{kind: Gte, want: true},
		{kind: Inc, want: true},
		{kind: Land, want: true},
		{kind: Lbrace, want: true},
		{kind: Lbrack, want: true},
		{kind: Lor, want: true},
		{kind: Lparen, want: true},
		{kind: Lt, want: true},
		{kind: Lte, want: true},
		{kind: Mod, want: true},
		{kind: ModAssign, want: true},
		{kind: Mul, want: true},
		{kind: MulAssign, want: true},
		{kind: Neq, want: true},
		{kind: Not, want: true},
		{kind: Or, want: true},
		{kind: OrAssign, want: true},
		{kind: Rbrace, want: true},
		{kind: Rbrack, want: true},
		{kind: Rparen, want: true},
		{kind: Semicolon, want: true},
		{kind: Shl, want: true},
		{kind: ShlAssign, want: true},
		{kind: Shr, want: true},
		{kind: ShrAssign, want: true},
		{kind: Sub, want: true},
		{kind: SubAssign, want: true},
		{kind: Xor, want: true},
		{kind: XorAssign, want: true},

		// Other tokens.
		{kind: Break, want: false},
		{kind: Case, want: false},
		{kind: Chan, want: false},
		{kind: Comment, want: false},
		{kind: Const, want: false},
		{kind: Continue, want: false},
		{kind: Default, want: false},
		{kind: Defer, want: false},
		{kind: Else, want: false},
		{kind: EOF, want: false},
		{kind: Fallthrough, want: false},
		{kind: Float, want: false},
		{kind: For, want: false},
		{kind: Func, want: false},
		{kind: Go, want: false},
		{kind: Goto, want: false},
		{kind: Ident, want: false},
		{kind: If, want: false},
		{kind: Imag, want: false},
		{kind: Import, want: false},
		{kind: Int, want: false},
		{kind: Interface, want: false},
		{kind: Map, want: false},
		{kind: Package, want: false},
		{kind: Range, want: false},
		{kind: Return, want: false},
		{kind: Rune, want: false},
		{kind: Select, want: false},
		{kind: String, want: false},
		{kind: Struct, want: false},
		{kind: Switch, want: false},
		{kind: Type, want: false},
		{kind: Var, want: false},
	}

	for i, g := range golden {
		got := g.kind.IsOperator()
		if got != g.want {
			t.Errorf("i=%d: IsOperand mismatch for token type %v; expected %t, got %t.", i, g.kind, g.want, got)
		}
	}
}

func TestKindIsLiteral(t *testing.T) {
	golden := []test{
		// Literals.
		{kind: Float, want: true},
		{kind: Ident, want: true},
		{kind: Imag, want: true},
		{kind: Int, want: true},
		{kind: Rune, want: true},
		{kind: String, want: true},

		// Other tokens.
		{kind: Add, want: false},
		{kind: AddAssign, want: false},
		{kind: And, want: false},
		{kind: AndAssign, want: false},
		{kind: Arrow, want: false},
		{kind: Assign, want: false},
		{kind: Break, want: false},
		{kind: Case, want: false},
		{kind: Chan, want: false},
		{kind: Clear, want: false},
		{kind: ClearAssign, want: false},
		{kind: Colon, want: false},
		{kind: Comma, want: false},
		{kind: Comment, want: false},
		{kind: Const, want: false},
		{kind: Continue, want: false},
		{kind: Dec, want: false},
		{kind: DeclAssign, want: false},
		{kind: Default, want: false},
		{kind: Defer, want: false},
		{kind: Div, want: false},
		{kind: DivAssign, want: false},
		{kind: Dot, want: false},
		{kind: Ellipsis, want: false},
		{kind: Else, want: false},
		{kind: EOF, want: false},
		{kind: Eq, want: false},
		{kind: Fallthrough, want: false},
		{kind: For, want: false},
		{kind: Func, want: false},
		{kind: Go, want: false},
		{kind: Goto, want: false},
		{kind: Gt, want: false},
		{kind: Gte, want: false},
		{kind: If, want: false},
		{kind: Import, want: false},
		{kind: Inc, want: false},
		{kind: Interface, want: false},
		{kind: Land, want: false},
		{kind: Lbrace, want: false},
		{kind: Lbrack, want: false},
		{kind: Lor, want: false},
		{kind: Lparen, want: false},
		{kind: Lt, want: false},
		{kind: Lte, want: false},
		{kind: Map, want: false},
		{kind: Mod, want: false},
		{kind: ModAssign, want: false},
		{kind: Mul, want: false},
		{kind: MulAssign, want: false},
		{kind: Neq, want: false},
		{kind: Not, want: false},
		{kind: Or, want: false},
		{kind: OrAssign, want: false},
		{kind: Package, want: false},
		{kind: Range, want: false},
		{kind: Rbrace, want: false},
		{kind: Rbrack, want: false},
		{kind: Return, want: false},
		{kind: Rparen, want: false},
		{kind: Select, want: false},
		{kind: Semicolon, want: false},
		{kind: Shl, want: false},
		{kind: ShlAssign, want: false},
		{kind: Shr, want: false},
		{kind: ShrAssign, want: false},
		{kind: Struct, want: false},
		{kind: Sub, want: false},
		{kind: SubAssign, want: false},
		{kind: Switch, want: false},
		{kind: Type, want: false},
		{kind: Var, want: false},
		{kind: Xor, want: false},
		{kind: XorAssign, want: false},
	}

	for i, g := range golden {
		got := g.kind.IsLiteral()
		if got != g.want {
			t.Errorf("i=%d: IsLiteral mismatch for token type %v; expected %t, got %t.", i, g.kind, g.want, got)
		}
	}
}
