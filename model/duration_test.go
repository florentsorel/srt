package model

import (
	"testing"
	"time"
)

func TestDurationString(t *testing.T) {
	tests := []struct {
		name     string
		d        Duration
		expected string
	}{
		{"zero", Duration(0), "00:00:00.000"},
		{"milliseconds", Duration(123 * time.Millisecond), "00:00:00.123"},
		{"under_one_second_round_down", Duration(999 * time.Millisecond), "00:00:00.999"},
		{"one_second", Duration(1 * time.Second), "00:00:01.000"},
		{"seconds_and_millis", Duration(1*time.Second + 42*time.Millisecond), "00:00:01.042"},
		{"one_minute", Duration(1 * time.Minute), "00:01:00.000"},
		{"minutes_seconds_millis", Duration(2*time.Minute + 3*time.Second + 7*time.Millisecond), "00:02:03.007"},
		{"one_hour", Duration(1 * time.Hour), "01:00:00.000"},
		{"hours_minutes_seconds_millis", Duration(10*time.Hour + 9*time.Minute + 8*time.Second + 765*time.Millisecond), "10:09:08.765"},
		{"negative_duration", Duration(-1*time.Hour - 2*time.Minute - 3*time.Second - 4*time.Millisecond), "01:02:03.004"},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.String(); got != tt.expected {
				t.Fatalf("[%d] expected %q, got %q", i, tt.expected, got)
			}
		})
	}
}
