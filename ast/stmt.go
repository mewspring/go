package ast

// A Stmt controls execution.
//
//    Statement =
//       Declaration | LabeledStmt | SimpleStmt |
//       GoStmt | ReturnStmt | BreakStmt | ContinueStmt | GotoStmt |
//       FallthroughStmt | Block | IfStmt | SwitchStmt | SelectStmt | ForStmt |
//       DeferStmt .
//
// ref: http://golang.org/ref/spec#Statements
type Stmt interface {
	// isStmt ensures that only statement nodes can be assigned to the Stmt
	// interface.
	isStmt()
}

// A SimpleStmt may precede the expression of if and switch statements, and the
// type switch guard of type switch statements. Simple statements may also
// specify the init and post statements of a for clause.
//
//    SimpleStmt = EmptyStmt | ExpressionStmt | SendStmt | IncDecStmt | Assignment | ShortVarDecl .
//
// ref: http://golang.org/ref/spec#Statements
type SimpleStmt interface {
	// isSimpleStmt ensures that only simple statement nodes can be assigned to
	// the SimpleStmt interface.
	isSimpleStmt()
}
