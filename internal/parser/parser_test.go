package parser

import (
	"testing"
)

func TestParseLine(t *testing.T) {

	tests := []struct {
		name       string
		input      string
		wantErr    bool
		wantPlayer int
		wantEvent  int
	}{
		{
			name:       "valid standard event",
			input:      "[14:10:00] 2 11 60",
			wantErr:    false,
			wantPlayer: 2,
			wantEvent:  11,
		},
		{
			name:       "empty line",
			input:      "   ",
			wantErr:    false,
			wantPlayer: 0,
		},
		{
			name:    "invalid format missing brackets",
			input:   "14:10:00 2 11",
			wantErr: true,
		},
		{
			name:    "invalid player ID",
			input:   "[14:10:00] a 11",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseLine(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != nil && !tt.wantErr {
				if got.PlayerID != tt.wantPlayer {
					t.Errorf("got PlayerID = %v, want %v", got.PlayerID, tt.wantPlayer)
				}
				if got.ID != tt.wantEvent {
					t.Errorf("got EventID = %v, want %v", got.ID, tt.wantEvent)
				}
			}
		})
	}
}
