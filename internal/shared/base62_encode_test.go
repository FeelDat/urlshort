package shared

import (
	"testing"
)

func TestBase62Encode(t *testing.T) {
	tests := []struct {
		name   string
		number uint64
		want   string
	}{
		{
			name:   "Encode Zero",
			number: 0,
			want:   "",
		},
		{
			name:   "Encode One",
			number: 1,
			want:   "b",
		},
		{
			name:   "Encode Large Number",
			number: 1234567890,
			want:   "uFhIvb",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Base62Encode(tt.number); got != tt.want {
				t.Errorf("Base62Encode(%v) = %v, want %v", tt.number, got, tt.want)
			}
		})
	}
}
