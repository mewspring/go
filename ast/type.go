package ast

import "github.com/mewlang/go/token"

// A Type determines the set of values and operations specific to values of that
// type. Types may be named or unnamed. Named types are specified by a (possibly
// qualified) type name; unnamed types are specified using a type literal, which
// composes a new type from existing types.
//
//    Type     = TypeName | TypeLit | "(" Type ")" .
//    TypeName = identifier | QualifiedIdent .
//    TypeLit  = ArrayType | StructType | PointerType | FunctionType | InterfaceType |
//               SliceType | MapType | ChannelType .
//
// http://golang.org/ref/spec#Types
type Type interface {
	// typeNode ensures that only type nodes can be assigned to the Type
	// interface.
	typeNode()
}

// A FuncType denotes the set of all functions with the same parameter and
// result types.
//
// Within a list of parameters or results, the names (IdentifierList) must
// either all be present or all be absent. If present, each name stands for one
// item (parameter or result) of the specified type. If absent, each type stands
// for one item of that type.
//
//    FunctionType   = "func" Signature .
//    Signature      = Parameters [ Result ] .
//    Result         = Parameters | Type .
//    Parameters     = "(" [ ParameterList [ "," ] ] ")" .
//    ParameterList  = ParameterDecl { "," ParameterDecl } .
//    ParameterDecl  = [ IdentifierList ] [ "..." ] Type .
//
// ref: http://golang.org/ref/spec#Function_types
type FuncType struct {
	// Zero or more parameters.
	Params []ParameterDecl
	// Zero or more results.
	Results []ParameterDecl
	// IsVariadic is true if the final parameter has an ellipsis type prefix, and
	// false otherwise.
	IsVariadic bool
}

// A ParameterDecl declares a list of parameters or results.
type ParameterDecl struct {
	// Parameter or result names; or nil.
	Names []token.Token
	// Parameter or result type.
	Type Type
}

// typeNode ensures that only type nodes can be assigned to the Type interface.
func (FuncType) typeNode() {}
