package parser

import (
	"testing"

	"github.com/florentsorel/srt/internal/lexer"
	"github.com/florentsorel/srt/internal/token"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	input := "11\n00:00:01,000 --> 00:00:04,000\nHello World!\nÇa va ? 😀\n\n2\n00:00:05,000 --> 00:00:07,000\nThis is a test."
	l, err := lexer.New(input)
	assert.NoError(t, err, "Expected no error from New")

	p := New(l)

	assert.Equal(t, token.NewToken(token.INDEX, "11", 1, 1), p.currentToken, "Expected current token to be INDEX (11), got %s (%s)", p.currentToken.Kind, p.currentToken.Literal)
	assert.Equal(t, token.NewToken(token.LF, "\n", 2, 0), p.nextToken, "Expected  next token to be LF (\n), got %s (%s)", p.currentToken.Kind, p.currentToken.Literal)
}

func TestParse(t *testing.T) {
	tests := []struct {
		input string
		err   string
	}{
		{
			input: `-->`,
			err:   "expected INDEX, got ILLEGAL at line 1, column 1",
		},
		{
			input: ` --> `,
			err:   "expected INDEX, got ARROW at line 1, column 2",
		},
		{
			input: `First line`,
			err:   "expected INDEX, got TEXT at line 1, column 1",
		},
		{
			input: `00:00:01,123`,
			err:   "expected INDEX, got TIMESTAMP at line 1, column 1",
		},
		{
			input: `1a`,
			err:   "expected INDEX, got TEXT at line 1, column 1",
		},
		{
			input: `1 00:00:01,123`,
			err:   "expected INDEX, got TEXT at line 1, column 1",
		},
		{
			input: `1 must be LF`,
			err:   "expected INDEX, got TEXT at line 1, column 1",
		},
		{
			input: `First line`,
			err:   "expected INDEX, got TEXT at line 1, column 1",
		},
		{
			input: `00:00:01,123`,
			err:   "expected INDEX, got TIMESTAMP at line 1, column 1",
		},
		{
			input: "1\n --> ",
			err:   "expected TIMESTAMP, got ARROW at line 2, column 2",
		},
		{
			input: "1\ntest --> 00:00:01,456",
			err:   "expected TIMESTAMP, got TEXT at line 2, column 1",
		},
		{
			input: "1\n00:00:00,456 test",
			err:   "expected ARROW, got TEXT at line 2, column 14",
		},
		{
			input: "1\n00:00:00,456 --> 1",
			err:   "expected TIMESTAMP, got INDEX at line 2, column 18",
		},
		{
			input: "1\n00:00:00,45",
			err:   "expected TIMESTAMP, got TEXT at line 2, column 1",
		},
		{
			input: "1\n00:00:00,456 --> 00:00:01,456",
			err:   "expected LF, got EOF at line 2, column 29",
		},
		{
			input: "1\n00:00:00,456 --> 00:0:01,456",
			err:   "expected TIMESTAMP, got TEXT at line 2, column 18",
		},
		{
			input: `1
00:00:00,456 --> 00:00:01,456
Hello world!

2
00:00:02,456 --> 00:00:03,456
`,
			err: "expected TEXT, got EOF at line 7, column 0",
		},
		{
			input: `1
00:00:00,000 --> 00:00:01,000
Hello
2
00:00:02,000 --> 00:00:03,000
World`,
			err: "expected EOC, got INDEX at line 4, column 1",
		},
		{
			input: `1
00:00:00,456 --> 00:00:01,456
Hello world!

2
00:00:02,456 --> 00:00:03,456
Test`,
			err: "",
		},
		{
			input: `1
00:00:00,456 --> 00:00:01,456
Hello world!
Welcome
Another line`,
			err: "",
		},
		{
			input: `1
00:00:00,456 --> 00:00:01,456
89 street`,
			err: "",
		},
		{
			input: `1
00:00:00,456 --> 00:00:01,456
13,23 euros`,
			err: "",
		},
		{
			input: `1
00:00:01,123 --> 00:00:01,456
First line
Second line

2
00:00:01,123 --> 00:00:01,456
First line

3
00:00:01,123 --> 00:00:01,456
- What are you doing?
- Nothing!`,
			err: "",
		},
	}

	for i, expected := range tests {
		l, _ := lexer.New(expected.input)
		p := New(l)
		_, err := p.Parse()

		if expected.err == "" {
			assert.NoError(t, err, "[%d] expected no error, got=%q.", i, err)
		} else {
			assert.EqualError(t, err, expected.err, "[%d] expected=%q, got=%q.", i, expected.err, err)
		}
	}
}
