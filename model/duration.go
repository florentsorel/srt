package model

import (
	"fmt"
	"time"
)

type Duration time.Duration

func (d Duration) String() string {
	duration := time.Duration(d)

	if duration < 0 {
		duration = -duration
	}

	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60
	milliseconds := int(duration.Milliseconds()) % 1000

	return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, seconds, milliseconds)
}

func (d Duration) Add(offset time.Duration) Duration {
	return Duration(time.Duration(d) + offset)
}
