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

// An ArrayType is a numbered sequence of elements of a single type, called the
// element type. The number of elements is called the length and is never
// negative.
//
//    ArrayType   = "[" ArrayLength "]" ElementType .
//    ArrayLength = Expression .
//    ElementType = Type .
//
// ref: http://golang.org/ref/spec#Array_types
type ArrayType struct {
	// Array length.
	Len Expr
	// Element type.
	Type Type
}

// A StructType consists of zero or more field declarations.
//
//    StructType     = "struct" "{" { FieldDecl ";" } "}" .
//    FieldDecl      = (IdentifierList Type | AnonymousField) [ Tag ] .
//    AnonymousField = [ "*" ] TypeName .
//    Tag            = string_lit .
//
// ref: http://golang.org/ref/spec#Struct_types
type StructType []FieldDecl

// A FieldDecl specifies the named elements of a struct, called fields, each of
// which has a name and a type. Field names may be specified explicitly
// (IdentifierList) or implicitly (AnonymousField).
//
// A field declared with a type but no explicit field name is an anonymous
// field, also called an embedded field or an embedding of the type in the
// struct. The unqualified type name acts as the field name.
type FieldDecl struct {
	// Field names; or nil.
	Names []token.Token
	// Field type.
	Type Type
	// Field tag.
	Tag token.Token
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
func (ArrayType) typeNode()  {}
func (StructType) typeNode() {}
func (FuncType) typeNode()   {}
