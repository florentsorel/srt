package internal

import "unicode/utf8"

type lexer struct {
	input           string
	length          int
	runePosition    int
	currentPosition int
	readPosition    int
	ch              rune
	line            int
	col             int
}

// newLexer creates and initializes a new lexer for the given input string.
// It reads the first character immediately to set up the lexer's state.
// The returned *lexer is intended for internal use only.
func newLexer(input string) *lexer {
	l := &lexer{
		input:           input,
		length:          len(input),
		runePosition:    0,
		currentPosition: 0,
		readPosition:    0,
		ch:              0,
		line:            1,
		col:             0,
	}
	l.readChar()
	return l
}

// readChar advances the lexer by one rune in the input, updating its
// currentPosition, readPosition, runePosition, line, and column.
// Handles multi-byte UTF-8 characters.
func (l *lexer) readChar() {
	if l.readPosition >= l.length {
		l.ch = 0 // EOF
		return
	}

	r, size := utf8.DecodeRuneInString(l.input[l.readPosition:])
	l.ch = r
	l.currentPosition = l.readPosition
	l.readPosition += size
	l.runePosition++

	if l.ch == '\n' {
		l.line++
		l.col = 0
	} else {
		l.col++
	}
}
