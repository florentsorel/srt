package srt

import (
	"io"
	"os"

	"github.com/florentsorel/srt/internal/lexer"
	"github.com/florentsorel/srt/internal/parser"
	"github.com/florentsorel/srt/model"
)

// Open reads the SRT file at the given path, parses its content,
// and returns a Subtitles struct or an error if reading or parsing fails.
func Open(path string) (*model.Subtitles, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return Parse(f)
}

// Parse reads from the provided io.Reader, parses the SRT content,
// and returns a Subtitles struct or an error if parsing fails.
func Parse(r io.Reader) (*model.Subtitles, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	l, err := lexer.New(string(b))
	if err != nil {
		return nil, err
	}

	p := parser.New(l)
	cues, err := p.Parse()
	if err != nil {
		return nil, err
	}

	return &model.Subtitles{Items: cues}, nil
}
