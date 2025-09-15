package model

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSubtitles_Shift(t *testing.T) {
	subtitles := Subtitles{
		Items: []Cue{
			{Index: 1, Start: Duration(1 * time.Second), End: Duration(3 * time.Second), Text: "First"},
			{Index: 2, Start: Duration(4 * time.Second), End: Duration(6 * time.Second), Text: "Second"},
		},
	}

	shifted := subtitles.Shift(2 * time.Second)

	assert.Equal(t, 2, len(shifted.Items), "Expected 2 items after shift, got %d", len(shifted.Items))
	assert.Equal(t, Duration(3*time.Second), shifted.Items[0].Start, "Expected first item's start to be 3 seconds, got %v", shifted.Items[0].Start)
	assert.Equal(t, Duration(5*time.Second), shifted.Items[0].End, "Expected first item's end to be 5 seconds, got %v", shifted.Items[0].End)
	assert.Equal(t, Duration(6*time.Second), shifted.Items[1].Start, "Expected second item's start to be 6 seconds, got %v", shifted.Items[1].Start)
	assert.Equal(t, Duration(8*time.Second), shifted.Items[1].End, "Expected second item's end to be 8 seconds, got %v", shifted.Items[1].End)
}

func TestSubtitles_RemoveAt(t *testing.T) {
	subtitles := Subtitles{
		Items: []Cue{
			{Index: 1, Start: Duration(1 * time.Second), End: Duration(3 * time.Second), Text: "First"},
			{Index: 2, Start: Duration(4 * time.Second), End: Duration(6 * time.Second), Text: "Second"},
			{Index: 3, Start: Duration(7 * time.Second), End: Duration(9 * time.Second), Text: "Third"},
		},
	}

	updated := subtitles.RemoveAt(1)

	assert.Equal(t, 2, len(updated.Items), "Expected 2 items after removal, got %d", len(updated.Items))
	assert.Equal(t, 1, updated.Items[0].Index, "Expected first item's index to be 1, got %d", updated.Items[0].Index)
	assert.Equal(t, "First", updated.Items[0].Text, "Expected first item's text to be 'First', got '%s'", updated.Items[0].Text)
	assert.Equal(t, 2, updated.Items[1].Index, "Expected second item's index to be 2, got %d", updated.Items[1].Index)
	assert.Equal(t, "Third", updated.Items[1].Text, "Expected second item's text to be 'Third', got '%s'", updated.Items[1].Text)
}

func TestSubtitles_RemoveAtIndices(t *testing.T) {
	subtitles := Subtitles{
		Items: []Cue{
			{Index: 1, Start: Duration(1 * time.Second), End: Duration(3 * time.Second), Text: "First"},
			{Index: 2, Start: Duration(4 * time.Second), End: Duration(6 * time.Second), Text: "Second"},
			{Index: 3, Start: Duration(7 * time.Second), End: Duration(9 * time.Second), Text: "Third"},
			{Index: 4, Start: Duration(10 * time.Second), End: Duration(12 * time.Second), Text: "Fourth"},
		},
	}

	updated := subtitles.RemoveAtIndices([]int{0, 2})

	assert.Equal(t, 2, len(updated.Items), "Expected 2 items after removal, got %d", len(updated.Items))
	assert.Equal(t, 1, updated.Items[0].Index, "Expected first item's index to be 1, got %d", updated.Items[0].Index)
	assert.Equal(t, "Second", updated.Items[0].Text, "Expected first item's text to be 'Second', got '%s'", updated.Items[0].Text)
	assert.Equal(t, 2, updated.Items[1].Index, "Expected second item's index to be 2, got %d", updated.Items[1].Index)
	assert.Equal(t, "Fourth", updated.Items[1].Text, "Expected second item's text to be 'Fourth', got '%s'", updated.Items[1].Text)
}

func TestSubtitles_Write(t *testing.T) {
	subtitles := Subtitles{
		Items: []Cue{
			{Index: 1, Start: Duration(1 * time.Second), End: Duration(3 * time.Second), Text: "First"},
			{Index: 2, Start: Duration(4 * time.Second), End: Duration(6 * time.Second), Text: "Second"},
		},
	}

	var sb strings.Builder
	n, err := subtitles.Write(&sb)
	assert.NoError(t, err, "Expected no error from Write")

	expected := "1\n00:00:01.000 --> 00:00:03.000\nFirst\n\n2\n00:00:04.000 --> 00:00:06.000\nSecond"
	assert.Equal(t, expected, sb.String(), "Expected output to be:\n%s\nGot:\n%s", expected, sb.String())
	assert.Equal(t, len(expected), n, "Expected number of bytes written to be %d, got %d", len(expected), n)
}
