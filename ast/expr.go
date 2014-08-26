package ast

// An Expr specifies the computation of a value by applying operators and
// functions to operands.
//
// Operators
//
// Operators combine operands into expressions.
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
type Expr interface {
	// exprNode ensures that only expression nodes can be assigned to the Expr
	// interface.
	exprNode()
}
