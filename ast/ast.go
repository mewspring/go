// TODO(u): Figure out how and where to store comments.

// Package ast declares the types used to represent abstract syntax trees of Go
// source code.
package ast

import "github.com/mewlang/go/token"

// A Package is constructed from one or more source files that together declare
// constants, types, variables and functions belonging to the package and which
// are accessible in all files of the same package.
//
// ref: http://golang.org/ref/spec#Packages
type Package struct {
	// Source files.
	Files []File
}

// A File consists of a package clause defining the package to which it belongs,
// followed by a possibly empty set of import declarations that declare packages
// whose contents it wishes to use, followed by a possibly empty set of
// declarations of functions, types, variables, and constants.
//
//    SourceFile    = PackageClause ";" { ImportDecl ";" } { TopLevelDecl ";" } .
//
//    PackageClause = "package" PackageName .
//    PackageName   = identifier .
//
// ref: http://golang.org/ref/spec#Source_file_organization
type File struct {
	// Package name.
	Pkg token.Token
	// Import declarations.
	Imps []ImportDecl
	// Top level declarations.
	Decls []Decl
}
