package model

import (
	"io"
	"strings"
	"time"
)

type Subtitles struct {
	Items []Cue
}

// Shift returns a new Subtitles with all Cue times shifted by the given offset.
func (s Subtitles) Shift(offset time.Duration) Subtitles {
	shiftedCues := make([]Cue, len(s.Items))
	for i, cue := range s.Items {
		shiftedCues[i] = cue.Shift(offset)
	}
	return Subtitles{Items: shiftedCues}
}

// RemoveAt removes the Cue at the specified index and returns a new Subtitles.
func (s Subtitles) RemoveAt(index int) Subtitles {
	if index < 0 || index >= len(s.Items) {
		return s
	}

	s.Items = append(s.Items[:index], s.Items[index+1:]...)

	for i := range s.Items {
		s.Items[i].Index = i + 1
	}

	return s
}

// RemoveAtIndices removes the Cues at the specified indices and returns a new Subtitles.
func (s Subtitles) RemoveAtIndices(indices []int) Subtitles {
	indexMap := make(map[int]struct{}, len(indices))
	for _, index := range indices {
		if index >= 0 && index < len(s.Items) {
			indexMap[index] = struct{}{}
		}
	}

	var newItems []Cue
	for i, cue := range s.Items {
		if _, found := indexMap[i]; !found {
			newItems = append(newItems, cue)
		}
	}

	for i := range newItems {
		newItems[i].Index = i + 1
	}

	return Subtitles{Items: newItems}
}

// Write writes the Subtitles in SRT format to the given io.Writer.
func (s Subtitles) Write(writer io.Writer) (int, error) {
	var b strings.Builder

	for i, cue := range s.Items {
		b.WriteString(cue.String())
		if i < len(s.Items)-1 {
			b.WriteByte('\n')
			b.WriteByte('\n')
		}
	}

	n, err := writer.Write([]byte(b.String()))
	return n, err
}
