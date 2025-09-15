package model

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCue_String(t *testing.T) {
	cue := Cue{
		Index: 1,
		Start: Duration(2 * time.Second),
		End:   Duration(5 * time.Second),
		Text:  "Hello, World!",
	}

	expected := "1\n00:00:02.000 --> 00:00:05.000\nHello, World!"
	assert.Equal(t, expected, fmt.Sprintf("%s", cue), "Expected Cue string to be:\n%s\nGot:\n%s", expected, cue.String())
}

func TestCue_Shift(t *testing.T) {
	cue := Cue{
		Index: 1,
		Start: Duration(2 * time.Second),
		End:   Duration(5 * time.Second),
		Text:  "Hello, World!",
	}

	assert.Equal(t, "00:00:02.000", cue.Start.String(), "Expected Start time to be 2 seconds, got %v", cue.Start)
	assert.Equal(t, "00:00:05.000", cue.End.String(), "Expected End time to be 5 seconds, got %v", cue.End)

	shiftedCue := cue.Shift(3*time.Second + 123*time.Millisecond)

	assert.Equal(t, "00:00:05.123", shiftedCue.Start.String(), "Expected Start time to be 5 seconds, got %v", shiftedCue.Start)
	assert.Equal(t, "00:00:08.123", shiftedCue.End.String(), "Expected End time to be 8 seconds, got %v", shiftedCue.End)
}
