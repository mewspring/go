package ast

import "github.com/mewlang/go/token"

// An Expr specifies the computation of a value by applying operators and
// functions to operands.
//
// Primary expressions
//
// Primary expressions are the operands for unary and binary expressions.
//
//    PrimaryExpr =
//       Operand |
//       Conversion |
//       BuiltinCall |
//       PrimaryExpr Selector |
//       PrimaryExpr Index |
//       PrimaryExpr Slice |
//       PrimaryExpr TypeAssertion |
//       PrimaryExpr Call .
//
//    Selector      = "." identifier .
//    Index         = "[" Expression "]" .
//    Slice         = "[" ( [ Expression ] ":" [ Expression ] ) |
//                        ( [ Expression ] ":" Expression ":" Expression )
//                    "]" .
//    TypeAssertion = "." "(" Type ")" .
//    Call          = "(" [ ArgumentList [ "," ] ] ")" .
//    ArgumentList  = ExpressionList [ "..." ] .
//
// ref: http://golang.org/ref/spec#Primary_expressions
type Expr interface {
	// exprNode ensures that only expression nodes can be assigned to the Expr
	// interface.
	exprNode()
}

// An UnaryExpr combines an unary operator and an operand into an expression.
//
// For integer operands, the unary operators +, -, and ^ are defined as follows:
//
//    +x                        is 0 + x
//    -x   negation             is 0 - x
//    ^x   bitwise complement   is m ^ x  with m = "all bits set to 1" for unsigned x
//                                        and  m = -1 for signed x
type UnaryExpr struct {
	// Unary operator; or NONE.
	Op token.Token
	// Primary (if Op is NONE) or unary expression.
	Expr Expr
}

// A BinaryExpr combines an operator and two operands into an expression.
//
//    Expression = UnaryExpr | Expression binary_op UnaryExpr .
//    UnaryExpr  = PrimaryExpr | unary_op UnaryExpr .
//
//    binary_op  = "||" | "&&" | rel_op | add_op | mul_op .
//    rel_op     = "==" | "!=" | "<" | "<=" | ">" | ">=" .
//    add_op     = "+" | "-" | "|" | "^" .
//    mul_op     = "*" | "/" | "%" | "<<" | ">>" | "&" | "&^" .
//
//    unary_op   = "+" | "-" | "!" | "^" | "*" | "&" | "<-" .
//
// ref: http://golang.org/ref/spec#Operators
type BinaryExpr struct {
	// Left-hand side expression.
	LHS Expr
	// Operator.
	Op token.Token
	// Right-hand side expression.
	RHS UnaryExpr
}

// exprNode ensures that only expression nodes can be assigned to the Expr
// interface.
func (UnaryExpr) exprNode()  {}
func (BinaryExpr) exprNode() {}
