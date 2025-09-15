package parser

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/florentsorel/srt/internal/lexer"
	"github.com/florentsorel/srt/internal/token"
	"github.com/florentsorel/srt/model"
)

type Parser struct {
	lexer        *lexer.Lexer
	currentToken token.Token
	nextToken    token.Token
}

// New creates a new parser with the given lexer and initializes its state.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l}
	p.readToken()
	p.readToken()
	return p
}

// readToken advances the parser to the next token.
func (p *Parser) readToken() {
	p.currentToken = p.nextToken
	p.nextToken = p.lexer.NextToken()
}

// Parse parses the entire input and returns a slice of cues or an error.
func (p *Parser) Parse() ([]model.Cue, error) {
	var cues []model.Cue

	for p.currentToken.Kind != token.EOF {
		cue, err := p.parseCue()
		if err != nil {
			return nil, err
		}
		cues = append(cues, *cue)
	}

	return cues, nil
}

// parseCue parses a single cue from the token stream.
func (p *Parser) parseCue() (*model.Cue, error) {
	var c model.Cue

	// Cue index
	if p.currentToken.Kind != token.INDEX {
		return nil, fmt.Errorf("expected INDEX, got %s at line %d, column %d", p.currentToken.Kind, p.currentToken.Line, p.currentToken.Column)
	}
	index, err := strconv.Atoi(p.currentToken.Literal)
	if err != nil {
		return nil, err
	}
	c.Index = index
	p.readToken()

	// Line feed
	if p.currentToken.Kind != token.LF {
		return nil, fmt.Errorf("expected LF, got %s at line %d, column %d", p.currentToken.Kind, p.currentToken.Line, p.currentToken.Column)
	}
	p.readToken()

	// Start
	if p.currentToken.Kind != token.TIMESTAMP {
		return nil, fmt.Errorf("expected TIMESTAMP, got %s at line %d, column %d", p.currentToken.Kind, p.currentToken.Line, p.currentToken.Column)
	}
	start, err := parseSRTTime(p.currentToken.Literal)
	if err != nil {
		return nil, err
	}
	c.Start = start
	p.readToken()

	// Arrow
	if p.currentToken.Kind != token.ARROW {
		return nil, fmt.Errorf("expected ARROW, got %s at line %d, column %d", p.currentToken.Kind, p.currentToken.Line, p.currentToken.Column)
	}
	p.readToken()

	// End
	if p.currentToken.Kind != token.TIMESTAMP {
		return nil, fmt.Errorf("expected TIMESTAMP, got %s at line %d, column %d", p.currentToken.Kind, p.currentToken.Line, p.currentToken.Column)
	}
	end, err := parseSRTTime(p.currentToken.Literal)
	if err != nil {
		return nil, err
	}
	c.End = end
	p.readToken()

	// Line feed
	if p.currentToken.Kind != token.LF {
		return nil, fmt.Errorf("expected LF, got %s at line %d, column %d", p.currentToken.Kind, p.currentToken.Line, p.currentToken.Column)
	}
	p.readToken()

	// Text
	var textLines []string
	if p.currentToken.Kind != token.TEXT {
		return nil, fmt.Errorf("expected TEXT, got %s at line %d, column %d",
			p.currentToken.Kind, p.currentToken.Line, p.currentToken.Column)
	}

	for p.currentToken.Kind == token.TEXT || p.currentToken.Kind == token.LF {
		if p.currentToken.Kind == token.TEXT {
			textLines = append(textLines, p.currentToken.Literal)
		}

		p.readToken()
	}
	c.Text = strings.Join(textLines, "\n")

	// End of cue
	if p.currentToken.Kind != token.EOC && p.currentToken.Kind != token.EOF {
		return nil, fmt.Errorf("expected EOC, got %s at line %d, column %d", p.currentToken.Kind, p.currentToken.Line, p.currentToken.Column)
	}

	p.readToken()

	return &c, nil
}

// parseSRTTime parses a time string in the format "HH:MM:SS,mmm" and returns a model.Duration.
func parseSRTTime(s string) (model.Duration, error) {
	var h, m, sec, ms int
	_, err := fmt.Sscanf(s, "%02d:%02d:%02d,%03d", &h, &m, &sec, &ms)
	if err != nil {
		return 0, err
	}

	d := time.Duration(h)*time.Hour +
		time.Duration(m)*time.Minute +
		time.Duration(sec)*time.Second +
		time.Duration(ms)*time.Millisecond

	return model.Duration(d), nil
}
