package lexer

import (
	"testing"

	"github.com/florentsorel/srt/internal/token"
	"github.com/stretchr/testify/assert"
)

func TestNewLexer(t *testing.T) {
	input := "Hello, World!\r\n"
	lexer, err := New(input)
	assert.NoError(t, err, "Expected no error from New")

	assert.NotNil(t, lexer, "Expected lexer to be initialized, got nil")
	assert.Equal(t, 14, lexer.length, "Expected lexer length to be 14, got %d", lexer.length)
	assert.Equal(t, 1, lexer.runePosition, "Expected initial rune position to be 1, got %d", lexer.runePosition)
	assert.Equal(t, 0, lexer.currentPosition, "Expected initial current position to be 0, got %d", lexer.currentPosition)
	assert.Equal(t, 1, lexer.readPosition, "Expected initial read position to be 1, got %d", lexer.readPosition)
	assert.Equal(t, 'H', lexer.ch, "Expected character to be 'H', got '%c'", lexer.ch)
	assert.Equal(t, 1, lexer.line, "Expected initial line to be 1, got %d", lexer.line)
	assert.Equal(t, 1, lexer.column, "Expected initial column to be 1, got %d", lexer.column)

	for i := 0; i <= 12; i++ {
		lexer.readChar()
	}
	assert.Equal(t, '\n', lexer.ch, "Expected character to be '\\n' after reading 5 chars, got '%c'", lexer.ch)
}

func TestNewLexerWithEmptyInput(t *testing.T) {
	input := ""
	lexer, err := New(input)
	assert.NoError(t, err, "Expected no error from newLexer")

	assert.NotNil(t, lexer, "Expected lexer to be initialized, got nil")
	assert.Equal(t, 0, lexer.length, "Expected lexer length to be 0, got %d", lexer.length)
	assert.Equal(t, 0, lexer.runePosition, "Expected initial rune position to be 0, got %d", lexer.runePosition)
	assert.Equal(t, 0, lexer.currentPosition, "Expected current position to be 0, got %d", lexer.currentPosition)
	assert.Equal(t, 0, lexer.readPosition, "Expected read position to be 0, got %d", lexer.readPosition)
	assert.Equal(t, rune(0), lexer.ch, "Expected character to be '0' (EOF), got '%c'", lexer.ch)
	assert.Equal(t, 1, lexer.line, "Expected initial line to be 1, got %d", lexer.line)
	assert.Equal(t, 0, lexer.column, "Expected initial column to be 0, got %d", lexer.column)
}

func TestNewLexerWithInvalidUTF8String(t *testing.T) {
	input := []byte{0x48, 0x65, 0x6C, 0x6C, 0x6F, 0x20, 0xC3, 0x28} // "Hello " + broken sequence
	lexer, err := New(string(input))
	assert.Error(t, err, "input string is not valid UTF-8")
	assert.Nil(t, lexer, "Expected lexer to be nil on error, got initialized lexer")
}

// TestReadChar verifies the lexer correctly reads characters,
// including multibyte UTF-8 runes. For example, 'é' is 2 bytes,
// '漢' is 3 bytes, and '😀' is 4 bytes.
func TestReadChar(t *testing.T) {
	input := "Hello éç漢😀\nè"
	lexer, err := New(input)
	assert.NoError(t, err, "Expected no error from New")

	cases := []struct {
		expectedRune         rune
		expectedCurrPos      int
		expectedReadPos      int
		expectedRunePosition int
		expectedCol          int
		expectedLine         int
	}{
		{'H', 0, 1, 1, 1, 1},
		{'e', 1, 2, 2, 2, 1},
		{'l', 2, 3, 3, 3, 1},
		{'l', 3, 4, 4, 4, 1},
		{'o', 4, 5, 5, 5, 1},
		{' ', 5, 6, 6, 6, 1},
		{'é', 6, 8, 7, 7, 1},     // 2-byte UTF-8 character
		{'ç', 8, 10, 8, 8, 1},    // 2-byte UTF-8 character
		{'漢', 10, 13, 9, 9, 1},   // 3-byte UTF-8 character
		{'😀', 13, 17, 10, 10, 1}, // 4-byte UTF-8 character
		{'\n', 17, 18, 11, 0, 2}, // newline resets column
		{'è', 18, 20, 12, 1, 2},  // 2-byte UTF-8 character
		{rune(0), 20, 20, 12, 1, 2},
	}

	for i, c := range cases {
		assert.Equal(t, c.expectedRune, lexer.ch, "[%d] Expected rune %q, got %q", i, c.expectedRune, lexer.ch)
		assert.Equal(t, c.expectedCurrPos, lexer.currentPosition, "[%d] Expected currentPosition %d, got %d", i, c.expectedCurrPos, lexer.currentPosition)
		assert.Equal(t, c.expectedReadPos, lexer.readPosition, "[%d] Expected readPosition %d, got %d", i, c.expectedReadPos, lexer.readPosition)
		assert.Equal(t, c.expectedRunePosition, lexer.runePosition, "[%d] Expected runePosition %d, got %d", i, c.expectedRunePosition, lexer.runePosition)
		assert.Equal(t, c.expectedCol, lexer.column, "[%d] Expected column %d, got %d", i, c.expectedCol, lexer.column)
		assert.Equal(t, c.expectedLine, lexer.line, "[%d] Expected line %d, got %d", i, c.expectedLine, lexer.line)

		lexer.readChar()
	}
}

func TestNextToken(t *testing.T) {
	input := "11\n00:00:01,000 --> 00:00:04,000\nHello World!\nÇa va ? 😀\n\n2\n00:00:05,000 --> 00:00:07,000\nThis is a test."
	lexer, err := New(input)
	assert.NoError(t, err, "Expected no error from New")

	tests := []struct {
		kind    token.TokenKind
		literal string
		line    int
		column  int
	}{
		{token.INDEX, "11", 1, 1},
		{token.LF, "\n", 2, 0},
		{token.TIMESTAMP, "00:00:01,000", 2, 1},
		{token.ARROW, "-->", 2, 14},
		{token.TIMESTAMP, "00:00:04,000", 2, 18},
		{token.LF, "\n", 3, 0},
		{token.TEXT, "Hello World!", 3, 1},
		{token.LF, "\n", 4, 0},
		{token.TEXT, "Ça va ? 😀", 4, 1},
		{token.EOC, "\n\n", 5, 0},
		{token.INDEX, "2", 6, 1},
		{token.LF, "\n", 7, 0},
		{token.TIMESTAMP, "00:00:05,000", 7, 1},
		{token.ARROW, "-->", 7, 14},
		{token.TIMESTAMP, "00:00:07,000", 7, 18},
		{token.LF, "\n", 8, 0},
		{token.TEXT, "This is a test.", 8, 1},
		{token.EOF, "", 8, 15},
	}

	for i, expected := range tests {
		tok := lexer.NextToken()
		assert.Equal(t, expected.kind, tok.Kind, "[%d] Expected token kind %q, got %q", i, expected.kind, tok.Kind)
		assert.Equal(t, expected.literal, tok.Literal, "[%d] Expected token literal %q, got %q", i, expected.literal, tok.Literal)
		assert.Equal(t, expected.line, tok.Line, "[%d] Expected token line %d, got %d", i, expected.line, tok.Line)
		assert.Equal(t, expected.column, tok.Column, "[%d] Expected token column %d, got %d", i, expected.column, tok.Column)
	}
}
