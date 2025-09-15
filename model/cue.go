package model

import (
	"fmt"
	"time"
)

type Cue struct {
	Index int
	Start Duration
	End   Duration
	Text  string
}

// String returns the Cue in SRT format.
func (c Cue) String() string {
	return fmt.Sprintf("%d\n%s --> %s\n%s", c.Index, c.Start.String(), c.End.String(), c.Text)
}

// Shift returns a new Cue with Start and End times shifted by the given offset.
func (c Cue) Shift(offset time.Duration) Cue {
	return Cue{
		Index: c.Index,
		Start: c.Start.Add(offset),
		End:   c.End.Add(offset),
		Text:  c.Text,
	}
}
