package utils

import (
	"testing"
)

func TestNormalizeToV4(t *testing.T) {
	tests := []struct {
		name      string
		major     int
		minor     int
		wantMajor int
		wantMinor int
		wantErr   bool
	}{
		{"4.18 passthrough", 4, 18, 4, 18, false},
		{"4.22 passthrough", 4, 22, 4, 22, false},
		{"5.0 normalizes to 4.23", 5, 0, 4, 23, false},
		{"5.1 normalizes to 4.24", 5, 1, 4, 24, false},
		{"5.2 normalizes to 4.25", 5, 2, 4, 25, false},
		{"major 3 unsupported", 3, 0, 0, 0, true},
		{"major 6 unsupported", 6, 0, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMajor, gotMinor, err := NormalizeToV4(tt.major, tt.minor)
			if (err != nil) != tt.wantErr {
				t.Errorf("NormalizeToV4(%d, %d) error = %v, wantErr %v", tt.major, tt.minor, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if gotMajor != tt.wantMajor || gotMinor != tt.wantMinor {
					t.Errorf("NormalizeToV4(%d, %d) = (%d, %d), want (%d, %d)",
						tt.major, tt.minor, gotMajor, gotMinor, tt.wantMajor, tt.wantMinor)
				}
			}
		})
	}
}

func TestDenormalizeFromV4(t *testing.T) {
	tests := []struct {
		name      string
		minor     int
		wantMajor int
		wantMinor int
	}{
		{"minor 18 stays 4.18", 18, 4, 18},
		{"minor 22 stays 4.22", 22, 4, 22},
		{"minor 23 becomes 5.0", 23, 5, 0},
		{"minor 24 becomes 5.1", 24, 5, 1},
		{"minor 25 becomes 5.2", 25, 5, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMajor, gotMinor := DenormalizeFromV4(tt.minor)
			if gotMajor != tt.wantMajor || gotMinor != tt.wantMinor {
				t.Errorf("DenormalizeFromV4(%d) = (%d, %d), want (%d, %d)",
					tt.minor, gotMajor, gotMinor, tt.wantMajor, tt.wantMinor)
			}
		})
	}
}
