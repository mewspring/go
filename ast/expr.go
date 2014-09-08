package ast

import (
	"github.com/mewlang/go/token"
	"github.com/mewlang/go/types"
)

// An Expr specifies the computation of a value by applying operators and
// functions to operands.
type Expr interface {
	// isExpr ensures that only expression nodes can be assigned to the Expr
	// interface.
	isExpr()
}

// An UnaryExpr combines an unary operator and an operand into an expression.
//
//    UnaryExpr  = PrimaryExpr | unary_op UnaryExpr .
//
//    unary_op   = "+" | "-" | "!" | "^" | "*" | "&" | "<-" .
//
// ref: http://golang.org/ref/spec#Operators
//
// For integer operands, the unary operators +, -, and ^ are defined as follows:
//
//    +x                        is 0 + x
//    -x   negation             is 0 - x
//    ^x   bitwise complement   is m ^ x  with m = "all bits set to 1" for unsigned x
//                                        and  m = -1 for signed x
type UnaryExpr struct {
	// Unary operator.
	Op token.Token
	// Unary operand; holds a PrimaryExpr or an UnaryExpr.
	Expr Expr
}

// A BinaryExpr combines an operator and two operands into an expression.
//
//    Expression = UnaryExpr | Expression binary_op UnaryExpr .
//
//    binary_op  = "||" | "&&" | rel_op | add_op | mul_op .
//    rel_op     = "==" | "!=" | "<" | "<=" | ">" | ">=" .
//    add_op     = "+" | "-" | "|" | "^" .
//    mul_op     = "*" | "/" | "%" | "<<" | ">>" | "&" | "&^" .
//
// ref: http://golang.org/ref/spec#Operators
type BinaryExpr struct {
	// Left-hand side operand.
	Left Expr
	// Operator.
	Op token.Token
	// Right-hand side operand; holds a PrimaryExpr or an UnaryExpr.
	Right Expr
}

// A PrimaryExpr represents a primary expression. Primary expressions are the
// operands for unary and binary expressions.
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
type PrimaryExpr interface {
	// isPrimaryExpr ensures that only primary expression nodes can be assigned
	// to the PrimaryExpr interface.
	isPrimaryExpr()
}

// A Conversion is an expression of the form T(x) where T is a type and x is an
// expression that can be converted to type T.
//
//    Conversion = Type "(" Expression [ "," ] ")" .
//
// ref: http://golang.org/ref/spec#Conversions
type Conversion struct {
	// Result type.
	Type types.Type
	// Original expression.
	Expr Expr
}

// isExpr ensures that only expression nodes can be assigned to the Expr
// interface.
func (UnaryExpr) isExpr()  {}
func (BinaryExpr) isExpr() {}
func (Conversion) isExpr() {}

// isPrimaryExpr ensures that only primary expression nodes can be assigned to
// the PrimaryExpr interface.
func (Conversion) isPrimaryExpr() {}
