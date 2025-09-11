package lexer

import (
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/florentsorel/srt/internal/token"
)

type Lexer struct {
	input           string
	length          int
	runePosition    int
	currentPosition int
	readPosition    int
	ch              rune
	line            int
	column          int
}

// New creates and initializes a new lexer for the given input string.
// It reads the first character immediately to set up the lexer's state.
// The returned *lexer is intended for internal use only.
func New(input string) (*Lexer, error) {
	if !utf8.ValidString(input) {
		return nil, errors.New("input string is not valid UTF-8")
	}

	input = strings.Replace(input, "\r\n", "\n", -1)
	l := &Lexer{
		input:           input,
		length:          len(input),
		runePosition:    0,
		currentPosition: 0,
		readPosition:    0,
		ch:              0,
		line:            1,
		column:          0,
	}
	l.readChar()
	return l, nil
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	if l.ch == 0 {
		return token.NewToken(token.EOF, "", l.line, l.column)
	}

	switch {
	case unicode.IsDigit(l.ch):
		start := l.currentPosition
		line := l.line
		column := l.column
		literal := l.readNumber()

		if l.ch == ':' || l.ch == ',' {
			for unicode.IsDigit(l.ch) || l.ch == ':' || l.ch == ',' {
				l.readChar()
			}
			literal = l.input[start:l.currentPosition]

			if isTimestampLiteral(literal) {
				tok = token.NewToken(token.TIMESTAMP, literal, line, column)
			} else {
				tok = token.NewToken(token.TEXT, literal, line, column)
			}
		} else {
			if l.ch != '\n' && l.ch != 0 {
				for l.ch != 0 && l.ch != '\n' {
					l.readChar()
				}
				literal = l.input[start:l.currentPosition]
				tok = token.NewToken(token.TEXT, literal, line, column)
			} else {
				tok = token.NewToken(token.INDEX, literal, line, column)
			}
		}
	case l.ch == '\n':
		line := l.line
		column := l.column

		l.readChar()

		if l.ch == '\n' {
			l.readChar()
			return token.NewToken(token.EOC, "\n\n", line, column)
		}

		return token.NewToken(token.LF, "\n", line, column)
	case l.ch == '-' && l.peekChar(1) == '-' && l.peekChar(2) == '>':
		start := l.currentPosition
		column := l.column

		if start == 0 || l.input[start-1] != ' ' {
			return token.NewToken(token.ILLEGAL, "-->", l.line, column)
		}

		l.readChar()
		l.readChar()
		l.readChar()

		if l.ch != ' ' {
			return token.NewToken(token.ILLEGAL, "-->", l.line, column)
		}

		literal := l.input[start:l.currentPosition]
		return token.NewToken(token.ARROW, literal, l.line, column)
	default:
		start := l.currentPosition
		line := l.line
		column := l.column

		for l.ch != 0 && l.ch != '\n' {
			l.readChar()
		}
		return token.NewToken(token.TEXT, l.input[start:l.currentPosition], line, column)
	}

	return tok
}

// readChar advances the Lexer by one rune in the input, updating its
// currentPosition, readPosition, runePosition, line, and column.
// Handles multibyte UTF-8 characters.
func (l *Lexer) readChar() {
	if l.readPosition >= l.length {
		l.currentPosition = l.readPosition
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
		l.column = 0
	} else {
		l.column++
	}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' {
		l.readChar()
	}
}

// peekChar returns the rune that is n runes ahead without advancing the Lexer.
// Returns 0 if it reaches the end of input.
func (l *Lexer) peekChar(n int) rune {
	if n <= 0 {
		return 0
	}

	pos := l.readPosition
	var r rune
	var size int

	for i := 0; i < n; i++ {
		if pos >= len(l.input) {
			return 0
		}
		r, size = utf8.DecodeRuneInString(l.input[pos:])
		pos += size
	}

	return r
}

func (l *Lexer) readNumber() string {
	start := l.currentPosition
	for unicode.IsDigit(l.ch) {
		l.readChar()
	}
	return l.input[start:l.currentPosition]
}

// isTimestampLiteral checks if a string matches the SRT timestamp shape
// HH:MM:SS,mmm with an exact length of 12 characters (e.g., 00:00:00,000).
func isTimestampLiteral(s string) bool {
	if len(s) != 12 {
		return false
	}

	// Positions: 0-1 HH, 2 ':', 3-4 MM, 5 ':', 6-7 SS, 8 ',', 9-11 mmm
	if !unicode.IsDigit(rune(s[0])) || !unicode.IsDigit(rune(s[1])) {
		return false
	}
	if s[2] != ':' {
		return false
	}
	if !unicode.IsDigit(rune(s[3])) || !unicode.IsDigit(rune(s[4])) {
		return false
	}
	if s[5] != ':' {
		return false
	}
	if !unicode.IsDigit(rune(s[6])) || !unicode.IsDigit(rune(s[7])) {
		return false
	}
	if s[8] != ',' {
		return false
	}
	if !unicode.IsDigit(rune(s[9])) || !unicode.IsDigit(rune(s[10])) || !unicode.IsDigit(rune(s[11])) {
		return false
	}
	return true
}
