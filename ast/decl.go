package ast

import (
	"github.com/mewlang/go/token"
	"github.com/mewlang/go/types"
)

// An ImportDecl consists of zero or more import specifiers.
//
//    ImportDecl = "import" ( ImportSpec | "(" { ImportSpec ";" } ")" ) .
//
// ref: http://golang.org/ref/spec#Import_declarations
type ImportDecl []ImportSpec

// An ImportSpec states that the source file containing the import declaration
// depends on functionality of the imported package and enables access to
// exported identifiers of that package. The import names an identifier
// (PackageName) to be used for access and an ImportPath that specifies the
// package to be imported.
//
//    ImportSpec = [ "." | PackageName ] ImportPath .
//    ImportPath = string_lit .
//
// ref: http://golang.org/ref/spec#Import_declarations
type ImportSpec struct {
	// Package name, or NONE.
	Name token.Token
	// Import path.
	Path token.Token
}

// A TopLevelDecl declares a constant, type, variable, function or method at the
// top level scope.
//
//    TopLevelDecl = Declaration | FunctionDecl | MethodDecl .
//
// ref: http://golang.org/ref/spec#Declarations_and_scope
type TopLevelDecl interface {
	// isTopLevelDecl ensures that only top level declaration nodes can be
	// assigned to the TopLevelDecl interface.
	isTopLevelDecl()
}

// A Decl binds a non-blank identifier to a function, method, label or package,
// or one or more non-blank identifiers to the same number of constants, types,
// or variables.
//
// The blank identifier may be used like any other identifier in a declaration,
// but it does not introduce a binding and thus is not declared. In the package
// block, the identifier init may only be used for init function declarations,
// and like the blank identifier it does not introduce a new binding.
//
//    Declaration  = ConstDecl | TypeDecl | VarDecl .
//
// ref: http://golang.org/ref/spec#Declarations_and_scope
type Decl interface {
	// isDecl ensures that only declaration nodes can be assigned to the Decl
	// interface.
	isDecl()
}

// A ConstDecl consists of zero or more constant specifiers.
//
//    ConstDecl      = "const" ( ConstSpec | "(" { ConstSpec ";" } ")" ) .
//    ConstSpec      = IdentifierList [ [ Type ] "=" ExpressionList ] .
//
//    IdentifierList = identifier { "," identifier } .
//    ExpressionList = Expression { "," Expression } .
//
// ref: http://golang.org/ref/spec#Constant_declarations
type ConstDecl []ValueSpec

// A VarDecl consists of zero or more variable specifiers.
//
//    VarDecl = "var" ( VarSpec | "(" { VarSpec ";" } ")" ) .
//    VarSpec = IdentifierList ( Type [ "=" ExpressionList ] | "=" ExpressionList ) .
//
// ref: http://golang.org/ref/spec#Variable_declarations
type VarDecl []ValueSpec

// A ValueSpec binds a list of constant or variable identifiers to the values of
// a list of constant or variable expressions respectively.
type ValueSpec struct {
	// Constant or variable names.
	Names []token.Token
	// Constant or variable type, or NONE.
	Type types.Type
	// Constant or variable value expressions, or nil.
	Vals []Expr
}

// A TypeDecl consists of zero or more type specifiers.
//
//    TypeDecl = "type" ( TypeSpec | "(" { TypeSpec ";" } ")" ) .
//
// ref: http://golang.org/ref/spec#Type_declarations
type TypeDecl []types.Name

// A FuncDecl binds an identifier, the function name, to a function.
//
//    FunctionDecl = "func" FunctionName ( Function | Signature ) .
//    FunctionName = identifier .
//    Function     = Signature FunctionBody .
//    FunctionBody = Block .
//
// ref: http://golang.org/ref/spec#Function_declarations
type FuncDecl struct {
	// Function name.
	Name token.Token
	// Function signature.
	Sig types.Func
	// Function body, or nil.
	Body Block
}

// A MethodDecl binds an identifier, the method name, to a method, and
// associates the method with the receiver's base type. A method is a function
// with a receiver.
//
//    MethodDecl = "func" Receiver MethodName ( Function | Signature ) .
//    Receiver   = Parameters .
//
// ref: http://golang.org/ref/spec#Method_declarations
type MethodDecl struct {
	// Receiver; must declare a single parameter.
	Receiver types.Parameter
	// Method name.
	Name token.Token
	// Method signature.
	Sig types.Func
	// Method body, or nil.
	Body Block
}

// isDecl ensures that only declaration nodes can be assigned to the Decl
// interface.
func (ConstDecl) isDecl() {}
func (TypeDecl) isDecl()  {}
func (VarDecl) isDecl()   {}

// isTopLevelDecl ensures that only top level declaration nodes can be assigned
// to the TopLevelDecl interface.
func (ConstDecl) isTopLevelDecl()  {}
func (TypeDecl) isTopLevelDecl()   {}
func (VarDecl) isTopLevelDecl()    {}
func (FuncDecl) isTopLevelDecl()   {}
func (MethodDecl) isTopLevelDecl() {}
