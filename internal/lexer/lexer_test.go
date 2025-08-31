package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLexer(t *testing.T) {
	input := "Hello, World!"
	lexer, err := New(input)
	assert.NoError(t, err, "Expected no error from newLexer")

	assert.NotNil(t, lexer, "Expected lexer to be initialized, got nil")
	assert.Equal(t, 13, lexer.length, "Expected lexer length to be 13, got %d", lexer.length)
	assert.Equal(t, 1, lexer.runePosition, "Expected initial rune position to be 1, got %d", lexer.runePosition)
	assert.Equal(t, 0, lexer.currentPosition, "Expected initial current position to be 0, got %d", lexer.currentPosition)
	assert.Equal(t, 1, lexer.readPosition, "Expected initial read position to be 1, got %d", lexer.readPosition)
	assert.Equal(t, 'H', lexer.ch, "Expected character to be 'H', got '%c'", lexer.ch)
	assert.Equal(t, 1, lexer.line, "Expected initial line to be 1, got %d", lexer.line)
	assert.Equal(t, 1, lexer.col, "Expected initial column to be 1, got %d", lexer.col)
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
	assert.Equal(t, 0, lexer.col, "Expected initial column to be 0, got %d", lexer.col)
}

func TestNewLexerWithInvalidUTF8String(t *testing.T) {
	input := []byte{0x48, 0x65, 0x6C, 0x6C, 0x6F, 0x20, 0xC3, 0x28} // "Hello " + broken sequence
	b := []byte(input)
	lexer, err := New(string(b))
	assert.Error(t, err, "input string is not valid UTF-8")
	assert.Nil(t, lexer, "Expected lexer to be nil on error, got initialized lexer")
}

// TestReadChar verifies the lexer correctly reads characters,
// including multi-byte UTF-8 runes. For example, 'é' is 2 bytes,
// '漢' is 3 bytes, and '😀' is 4 bytes.
func TestReadChar(t *testing.T) {
	input := "Hello éç漢😀\nè"
	lexer, err := New(input)
	assert.NoError(t, err, "Expected no error from newLexer")

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
		{rune(0), 18, 20, 12, 1, 2},
	}

	for i, c := range cases {
		assert.Equal(t, c.expectedRune, lexer.ch, "[%d] Expected rune %q, got %q", i, c.expectedRune, lexer.ch)
		assert.Equal(t, c.expectedCurrPos, lexer.currentPosition, "[%d] Expected currentPosition %d, got %d", i, c.expectedCurrPos, lexer.currentPosition)
		assert.Equal(t, c.expectedReadPos, lexer.readPosition, "[%d] Expected readPosition %d, got %d", i, c.expectedReadPos, lexer.readPosition)
		assert.Equal(t, c.expectedRunePosition, lexer.runePosition, "[%d] Expected runePosition %d, got %d", i, c.expectedRunePosition, lexer.runePosition)
		assert.Equal(t, c.expectedCol, lexer.col, "[%d] Expected column %d, got %d", i, c.expectedCol, lexer.col)
		assert.Equal(t, c.expectedLine, lexer.line, "[%d] Expected line %d, got %d", i, c.expectedLine, lexer.line)

		lexer.readChar()
	}
}
