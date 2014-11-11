WIP
---

This project is a *work in progress*. The implementation is *incomplete* and
subject to change. The documentation may be inaccurate.

go
==

[![Build Status](https://travis-ci.org/mewlang/go.svg?branch=master)](https://travis-ci.org/mewlang/go)
[![Coverage Status](https://img.shields.io/coveralls/mewlang/go.svg)](https://coveralls.io/r/mewlang/go?branch=master)
[![GoDoc](https://godoc.org/github.com/mewlang/go?status.svg)](https://godoc.org/github.com/mewlang/go)

The aim of this repository is to gain insight into compiler design by
implementing a compiler for the Go programming language.

Documentation
-------------

Documentation provided by GoDoc.

- [ast][]: declares the types used to represent abstract syntax trees of Go
source code.
- [lexer][]: implements lexical tokenization of Go source code.
- [token][]: defines constants representing the lexical tokens of the Go
programming language.
- [types][]: declares the data types of the Go programming language.

[ast]: http://godoc.org/github.com/mewlang/go/ast
[lexer]: http://godoc.org/github.com/mewlang/go/lexer
[token]: http://godoc.org/github.com/mewlang/go/token
[types]: http://godoc.org/github.com/mewlang/go/types

public domain
-------------

This code is hereby released into the *[public domain][]*.

[public domain]: https://creativecommons.org/publicdomain/zero/1.0/

BSD license
-----------

Any code or documentation directly derived from the [standard Go source code][]
is governed by a [BSD license][].

[standard Go source code]: https://code.google.com/p/go
[BSD license]: http://golang.org/LICENSE
