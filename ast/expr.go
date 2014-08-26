package ast

type Expr interface {
	// exprNode ensures that only expression nodes can be assigned to the Expr
	// interface.
	exprNode()
}
