package ast

import (
	"github.com/mewlang/go/token"
	"github.com/mewlang/go/types"
)

// Operands denote the elementary values in an expression. An operand may be a
// literal, a (possibly qualified) non-blank identifier denoting a constant,
// variable, or function, a method expression yielding a function, or a
// parenthesized expression.
//
// The blank identifier may appear as an operand only on the left-hand side of
// an assignment.
//
//    Operand     = Literal | OperandName | MethodExpr | "(" Expression ")" .
//    Literal     = BasicLit | CompositeLit | FunctionLit .
//
// ref: http://golang.org/ref/spec#Operands

// A BasicLit is an integer, floating-point, imaginary, rune, or string literal.
//
//    BasicLit    = int_lit | float_lit | imaginary_lit | rune_lit | string_lit .
//
// ref: http://golang.org/ref/spec#Operands
type BasicLit token.Token

// A CompositeLit constructs a value for a struct, array, slice, or map and
// creates a new value each time it is evaluated. Composite literals consist of
// the type of the value followed by a brace-bound list of composite elements.
//
//    CompositeLit  = LiteralType LiteralValue .
//    LiteralType   = StructType | ArrayType | "[" "..." "]" ElementType |
//                    SliceType | MapType | TypeName .
//    LiteralValue  = "{" [ ElementList [ "," ] ] "}" .
//    ElementList   = Element { "," Element } .
//    Element       = [ Key ":" ] Value .
//    Key           = FieldName | ElementIndex .
//    FieldName     = identifier .
//    ElementIndex  = Expression .
//    Value         = Expression | LiteralValue .
//
// ref: http://golang.org/ref/spec#Composite_literals
type CompositeLit struct {
	// Literal type; holds a Struct, Array, Slice, Map, or Name from the types
	// package.
	Type types.Type
	// Literal values.
	Vals []CompositeElement
}

// A CompositeElement may be a single expression or a key-value pair.
type CompositeElement struct {
	// Element key, or nil; holds an identifier (token.Token) or an Expr.
	Key interface{}
	// Element value; holds an Expr or a []CompositeElement.
	Val interface{}
}

// A FuncLit represents an anonymous function.
//
//    FunctionLit = "func" Function .
//
// ref: http://golang.org/ref/spec#Function_literals
type FuncLit struct {
	// Function signature.
	Sig types.Func
	// Function body, or nil.
	Body Block
}

// An OperandName is a (possibly qualified) non-blank identifier denoting a
// constant, variable, or function.
//
//    OperandName = identifier | QualifiedIdent.
//
// ref: http://golang.org/ref/spec#Operands
type OperandName token.Token

// A MethodExpr yields a function equivalent to the specified method with an
// explicit receiver as its first argument.
//
//    MethodExpr   = ReceiverType "." MethodName .
//    ReceiverType = TypeName | "(" "*" TypeName ")" | "(" ReceiverType ")" .
//
// ref: http://golang.org/ref/spec#Method_expressions
type MethodExpr struct {
	// Receiver type; holds a Name or a Pointer with a Name base type from the
	// types package.
	ReceiverType types.Type
	// Method name.
	Name token.Token
}

// isPrimaryExpr ensures that only primary expression nodes can be assigned to
// the PrimaryExpr interface.
func (BasicLit) isPrimaryExpr()     {}
func (CompositeLit) isPrimaryExpr() {}
func (FuncLit) isPrimaryExpr()      {}
func (OperandName) isPrimaryExpr()  {}
func (MethodExpr) isPrimaryExpr()   {}
