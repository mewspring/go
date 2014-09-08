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
//    TypeAssertion = "." "(" Type ")" .
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

// A CallExpr is a function call or a method invocation.
//
//    PrimaryExpr Call .
//
//    Call          = "(" [ ArgumentList [ "," ] ] ")" .
//    ArgumentList  = ExpressionList [ "..." ] .
//
// ref: http://golang.org/ref/spec#Calls
//
// Built-in functions are predeclared. They are called like any other function
// but some of them accept a type instead of an expression as the first
// argument.
//
//    BuiltinCall = identifier "(" [ BuiltinArgs [ "," ] ] ")" .
//    BuiltinArgs = Type [ "," ArgumentList ] | ArgumentList .
//
// ref: http://golang.org/ref/spec#Built-in_functions
type CallExpr struct {
	// Function or method expression.
	Func PrimaryExpr
	// Function or method arguments; each argument is an Expr, except for when
	// one of the built-in functions make or new is invoked, in which case the
	// first argument is a types.Type.
	Args []interface{}
	// Specifies if the final argument is suffixed with an ellipsis.
	HasEllipsis bool
}

// A SelectorExpr denotes a field or method of a primary expression with an
// identifier called the selector.
//
//    PrimaryExpr Selector .
//
//    Selector = "." identifier .
//
// ref: http://golang.org/ref/spec#Selectors
type SelectorExpr struct {
	// Expression.
	Expr PrimaryExpr
	// Field or method selector.
	Selector token.Token
}

// An IndexExpr denotes an element of an array, pointer to array, slice, string,
// or map.
//
//    PrimaryExpr Index |
//
//    Index = "[" Expression "]" .
//
// ref: http://golang.org/ref/spec#Index_expressions
type IndexExpr struct {
	// Expression.
	Expr PrimaryExpr
	// Index expression.
	Index Expr
}

// A SliceExpr constructs a substring or slice from a string, array, pointer to
// array, or slice. There are two variants: a simple form that specifies a low
// and high bound, and a full form that also specifies a bound on the capacity.
//
//    PrimaryExpr Slice .
//
//    Slice         = "[" ( [ Expression ] ":" [ Expression ] ) |
//                        ( [ Expression ] ":" Expression ":" Expression )
//                    "]" .
//
// ref: http://golang.org/ref/spec#Slice_expressions
type SliceExpr struct {
	// Lower bound.
	Low Expr
	// Higher bound.
	High Expr
	// Capacity.
	Cap Expr
}

// isExpr ensures that only expression nodes can be assigned to the Expr
// interface.
func (UnaryExpr) isExpr()    {}
func (BinaryExpr) isExpr()   {}
func (Conversion) isExpr()   {}
func (CallExpr) isExpr()     {}
func (SelectorExpr) isExpr() {}
func (IndexExpr) isExpr()    {}
func (SliceExpr) isExpr()    {}

// isPrimaryExpr ensures that only primary expression nodes can be assigned to
// the PrimaryExpr interface.
func (Conversion) isPrimaryExpr()   {}
func (CallExpr) isPrimaryExpr()     {}
func (SelectorExpr) isPrimaryExpr() {}
func (IndexExpr) isPrimaryExpr()    {}
func (SliceExpr) isPrimaryExpr()    {}
