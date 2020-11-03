package misc

import (
	"fmt"
	"testing"
)

func TestNumberRange_IsConnected(t *testing.T) {
	tests := []struct {
		name string
		r1   Range
		r2   Range
		want bool
	}{
		{"t1", NewClosedRange(1, 3), NewClosedRange(3, 5), true},
		{"t2", NewClosedRange(3, 5), NewClosedRange(1, 3), true},
		{"t3", NewClosedRange(1, 3), NewClosedRange(3.1, 5), false},
		{"t4", NewClosedRange(1, 5), NewClosedRange(2, 6), true},
		{"t5", NewClosedRange(1, 5), NewClosedRange(2, 3), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r1.IsConnected(tt.r2); got != tt.want {
				t.Errorf("NumberRange.IsConnected() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumberRange_String(t *testing.T) {
	tests := []struct {
		name string
		r    Range
		want string
	}{
		{"t1", NewClosedRange(1, 3), "[1,3]"},
		{"t2", NewClosedRange(1.32, 3.12), "[1.32,3.12]"},
		{"t3", NewOpenRange(1.32, 3.12), "(1.32,3.12)"},
		{"t4", NewOpenClosedRange(1.32, 3.12), "(1.32,3.12]"},
		{"t5", NewClosedOpenRange(1.32, 3.12), "[1.32,3.12)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fmt.Sprint(tt.r); got != tt.want {
				t.Errorf("NumberRange.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
