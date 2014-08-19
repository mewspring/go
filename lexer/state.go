package lexer

// A stateFn represents the state of the lexer as a function that returns a
// state function.
type stateFn func(l *lexer) stateFn

func lexToken(l *lexer) stateFn {
	panic("not yet implemented.")
}
