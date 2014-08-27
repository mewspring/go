package ast

import "github.com/mewlang/go/token"

// An ImportDecl consists of zero or more import specifiers.
//
//    ImportDecl = "import" ( ImportSpec | "(" { ImportSpec ";" } ")" ) .
//    ImportSpec = [ "." | PackageName ] ImportPath .
//    ImportPath = string_lit .
//
// ref: http://golang.org/ref/spec#Import_declarations
type ImportDecl []ImportSpec

// An ImportSpec states that the source file containing the import declaration
// depends on functionality of the imported package and enables access to
// exported identifiers of that package. The import names an identifier
// (PackageName) to be used for access and an ImportPath that specifies the
// package to be imported.
type ImportSpec struct {
	// Package name; or NONE.
	Name token.Token
	// Import path.
	Path token.Token
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
//    TopLevelDecl = Declaration | FunctionDecl | MethodDecl .
//
// ref: http://golang.org/ref/spec#Declarations_and_scope
type Decl interface {
	// declNode ensures that only declaration nodes can be assigned to the Decl
	// interface.
	declNode()
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
	// Constant or variable type; or NONE.
	Type Type
	// Constant or variable value expressions; or nil.
	Vals []Expr
}

// A TypeDecl consists of zero or more type specifiers.
//
//    TypeDecl = "type" ( TypeSpec | "(" { TypeSpec ";" } ")" ) .
//    TypeSpec = identifier Type .
//
// ref: http://golang.org/ref/spec#Type_declarations
type TypeDecl []TypeSpec

// A TypeSpec binds an identifier, the type name, to a new type that has the
// same underlying type as an existing type, and operations defined for the
// existing type are also defined for the new type.
type TypeSpec struct {
	// Type name.
	Name token.Token
	// Type.
	Type Type
}

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
	Type FuncType
	// Function body; or nil.
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
	Receiver ParameterDecl
	// Method name.
	Name token.Token
	// Method signature.
	Type FuncType
	// Method body; or nil.
	Body Block
}

// declNode ensures that only declaration nodes can be assigned to the Decl
// interface.
func (ConstDecl) declNode()  {}
func (FuncDecl) declNode()   {}
func (MethodDecl) declNode() {}
func (TypeDecl) declNode()   {}
func (VarDecl) declNode()    {}
