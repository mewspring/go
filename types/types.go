// Package types declares the data types of the Go programming language.
package types

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
	// isType ensures that only type nodes can be assigned to the Type interface.
	isType()
}

// Basic specifies a predeclared boolean, numeric, or string type of the Go
// programming language. The following types are implicitly declared in the
// universe block:
//
//    bool byte complex64 complex128 error float32 float64
//    int int8 int16 int32 int64 rune string
//    uint uint8 uint16 uint32 uint64 uintptr
//
// ref: http://golang.org/ref/spec#Predeclared_identifiers
type Basic uint8

// Basic types.
const (
	Bool Basic = iota
	Byte
	Complex64
	Complex128
	Error
	Float32
	Float64
	Int
	Int8
	Int16
	Int32
	Int64
	Rune
	String
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Uintptr
)

// A Name binds an identifier, the type name, to a new type that has the same
// underlying type as an existing type, and operations defined for the existing
// type are also defined for the new type.
//
//    TypeSpec = identifier Type .
//
// ref: http://golang.org/ref/spec#Type_declarations
type Name struct {
	// Type name.
	Name token.Token
	// Underlying type.
	Type Type
}

// An Array is a numbered sequence of elements of a single type, called the
// element type. The number of elements is called the length and is never
// negative.
//
//    ArrayType   = "[" ArrayLength "]" ElementType .
//    ArrayLength = Expression .
//    ElementType = Type .
//
// ref: http://golang.org/ref/spec#Array_types
type Array struct {
	// Array length; holds an ast.ConstExpr.
	Len interface{}
	// Element type.
	Elem Type
}

// A Struct consists of zero or more fields.
//
//    StructType     = "struct" "{" { FieldDecl ";" } "}" .
//
// ref: http://golang.org/ref/spec#Struct_types
type Struct []Field

// A Field specifies the named elements of a struct, called fields, each of
// which has a name and a type. Field names may be specified explicitly
// (IdentifierList) or implicitly (AnonymousField).
//
// A field declared with a type but no explicit field name is an anonymous
// field, also called an embedded field or an embedding of the type in the
// struct. The unqualified type name acts as the field name.
//
//    FieldDecl      = (IdentifierList Type | AnonymousField) [ Tag ] .
//    AnonymousField = [ "*" ] TypeName .
//    Tag            = string_lit .
//
// ref: http://golang.org/ref/spec#Struct_types
type Field struct {
	// Field names; or nil.
	Names []token.Token
	// Field type; holds an anonymous field (a type name or a pointer to a type
	// name) if Names is nil.
	Type Type
	// Field tag; or NONE.
	Tag token.Token
}

// A Pointer denotes the set of all pointers to variables of a given type,
// called the base type of the pointer.
//
//    PointerType = "*" BaseType .
//    BaseType    = Type .
//
// ref: http://golang.org/ref/spec#Pointer_types
type Pointer struct {
	// Pointer base type.
	Base Type
}

// A Func denotes the set of all functions with the same parameter and result
// types.
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
type Func struct {
	// Zero or more parameters.
	Params []Parameter
	// Zero or more results.
	Results []Parameter
	// IsVariadic is true if the final parameter has an ellipsis type prefix, and
	// false otherwise.
	IsVariadic bool
}

// A Parameter declares a list of parameters or results.
type Parameter struct {
	// Parameter or result names; or nil.
	Names []token.Token
	// Parameter or result type.
	Type Type
}

// An Interface specifies a method set called its interface. A variable of
// interface type can store a value of any type with a method set that is any
// superset of the interface. Such a type is said to implement the interface.
//
//    InterfaceType     = "interface" "{" { MethodSpec ";" } "}" .
//    MethodSpec        = MethodName Signature | InterfaceTypeName .
//    MethodName        = identifier .
//    InterfaceTypeName = TypeName .
//
// ref: http://golang.org/ref/spec#Interface_types
type Interface []Method

// A Method denotes the set of all methods with the same method name, and
// parameter and result types.
type Method struct {
	// Method name (if Sig != nil) or interface type name.
	Name token.Token
	// Method signature; or nil.
	Sig Func
}

// A Slice is a descriptor for a contiguous segment of an underlying array and
// provides access to a numbered sequence of elements from that array.
//
//    SliceType = "[" "]" ElementType .
//
// ref: http://golang.org/ref/spec#Slice_types
type Slice struct {
	// Element type.
	Elem Type
}

// A Map is an unordered group of elements of one type, called the element type,
// indexed by a set of unique keys of another type, called the key type.
//
//    MapType = "map" "[" KeyType "]" ElementType .
//    KeyType = Type .
//
// ref: http://golang.org/ref/spec#Map_types
type Map struct {
	// Key type.
	Key Type
	// Element type.
	Elem Type
}

// A Chan provides a mechanism for concurrently executing functions to
// communicate by sending and receiving values of a specified element type.
//
//    ChannelType = ( "chan" | "chan" "<-" | "<-" "chan" ) ElementType .
//
// ref: http://golang.org/ref/spec#Channel_types
type Chan struct {
	// Channel direction.
	Dir ChanDir
	// Element type.
	Elem Type
}

// ChanDir is a bitfield which specifies the channel direction; send, receive or
// bidirectional.
type ChanDir uint8

// Channel directions.
const (
	Send ChanDir = 1 << iota
	Recv
)

// isType ensures that only type nodes can be assigned to the Type interface.
func (Basic) isType()     {}
func (Name) isType()      {}
func (Array) isType()     {}
func (Struct) isType()    {}
func (Pointer) isType()   {}
func (Func) isType()      {}
func (Interface) isType() {}
func (Slice) isType()     {}
func (Map) isType()       {}
func (Chan) isType()      {}
