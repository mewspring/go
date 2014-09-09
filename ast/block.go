package ast

// A Block is a possibly empty sequence of declarations and statements within
// matching brace brackets.
//
//    Block         = "{" StatementList "}" .
//    StatementList = { Statement ";" } .
//
// ref: http://golang.org/ref/spec#Blocks
type Block []Stmt
