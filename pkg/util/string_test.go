package util_test

import (
	"strings"
	"testing"

	"goapp/pkg/util"
)

func TestRandString(t *testing.T) {
	const hexDigits = "0123456789ABCDEF"

	tests := []struct {
		name      string
		length    int
		wantPanic bool
	}{
		{
			name:      "negative length",
			length:    -1,
			wantPanic: true,
		},
		{
			name:   "zero length",
			length: 0,
		},
		{
			name:   "one character",
			length: 1,
		},
		{
			name:   "ten characters",
			length: 10,
		},
		{
			name:   "one hundred characters",
			length: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				defer func() {
					if recover() == nil {
						t.Errorf("RandString(%d) did not panic", tt.length)
					}
				}()
			}

			value := util.RandString(tt.length)

			if len(value) != tt.length {
				t.Errorf("RandString(%d) length = %d, want %d", tt.length, len(value), tt.length)
			}

			for _, char := range value {
				if !strings.ContainsRune(hexDigits, char) {
					t.Errorf("RandString(%d) returned non-hex character %q in %q", tt.length, char, value)
				}
			}
		})
	}
}

var benchmarkRandString string

func BenchmarkRandString(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		benchmarkRandString = util.RandString(10)
	}
}
