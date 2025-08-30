package internal

import (
	"testing"
)

func TestNewLexer(t *testing.T) {
	input := "Hello, World!"
	lexer := newLexer(input)

	if lexer == nil {
		t.Fatal("Expected lexer to be initialized, got nil")
	}

	if lexer.length != 13 {
		t.Fatalf("Expected lexer length to be %d, got %d", 13, lexer.length)
	}

	if lexer.currentPosition != 0 {
		t.Fatalf("Expected current position to be 0, got %d", lexer.currentPosition)
	}

	if lexer.readPosition != 1 {
		t.Fatalf("Expected read position to be 1, got %d", lexer.readPosition)
	}

	if lexer.ch != 'H' {
		t.Fatalf("Expected character to be 'H', got '%c'", lexer.ch)
	}
}

func TestNewLexerWithEmptyInput(t *testing.T) {
	input := ""
	lexer := newLexer(input)

	if lexer == nil {
		t.Fatal("Expected lexer to be initialized, got nil")
	}

	if lexer.length != 0 {
		t.Fatalf("Expected lexer length to be %d, got %d", 0, lexer.length)
	}

	if lexer.currentPosition != 0 {
		t.Fatalf("Expected current position to be 0, got %d", lexer.currentPosition)
	}

	if lexer.readPosition != 0 {
		t.Fatalf("Expected read position to be 0, got %d", lexer.readPosition)
	}

	if lexer.ch != 0 {
		t.Fatalf("Expected character to be '0' (EOF), got '%c'", lexer.ch)
	}
}

// TestReadChar verifies the lexer correctly reads characters,
// including multi-byte UTF-8 runes. For example, 'é' is 2 bytes,
// '漢' is 3 bytes, and '😀' is 4 bytes.
func TestReadChar(t *testing.T) {
	input := "Hello éç漢😀\nè"
	lexer := newLexer(input)

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
	}

	for i, c := range cases {
		if lexer.ch != c.expectedRune {
			t.Errorf("[%d] Expected rune %q, got %q", i, c.expectedRune, lexer.ch)
		}
		if lexer.currentPosition != c.expectedCurrPos {
			t.Errorf("[%d] Expected currentPosition %d, got %d", i, c.expectedCurrPos, lexer.currentPosition)
		}
		if lexer.readPosition != c.expectedReadPos {
			t.Errorf("[%d] Expected readPosition %d, got %d", i, c.expectedReadPos, lexer.readPosition)
		}
		if lexer.runePosition != c.expectedRunePosition {
			t.Errorf("[%d] Expected runePosition %d, got %d", i, c.expectedRunePosition, lexer.runePosition)
		}
		if lexer.col != c.expectedCol {
			t.Errorf("[%d] Expected column %d, got %d", i, c.expectedCol, lexer.col)
		}
		if lexer.line != c.expectedLine {
			t.Errorf("[%d] Expected line %d, got %d", i, c.expectedLine, lexer.line)
		}

		lexer.readChar()
	}
}
