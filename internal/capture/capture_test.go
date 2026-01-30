package capture

import "testing"

func TestParseYoutubeTime(t *testing.T) {
	cases := []struct {
		input string
		want  int
		valid bool
	}{
		{"0", 0, true},          // Zero seconds
		{"0s", 0, true},         // Zero seconds and valid suffix
		{"0h0m0s", 0, true},     // Composite timestamps with zeros
		{"00h00m00s", 0, true},  // Composite timestamps with leading zeros
		{"42", 42, true},        // Plain seconds
		{"600s", 600, true},     // Seconds with suffix
		{"1h", 3600, true},      // Hours only
		{"10m", 600, true},      // Minutes only
		{"10m2s", 602, true},    // Minutes + seconds
		{"10m02s", 602, true},   // Leading zero seconds
		{"1h2m3s", 3723, true},  // Hours + minutes + seconds
		{"00h01m02s", 62, true}, // Leading zero hours/minutes/seconds

		{"abc", 0, false},       // Invalid
		{"10x", 0, false},       // Invalid unit
		{"m10s", 0, false},      // Malformed, with leading time unit w/o time
		{"0x1a4m10s", 0, false}, // Malformed
		{"10d", 0, false},       // Typo e.g. 'd' instead of 's'
		{"30s2m", 0, false},     // Out of order units
		{"2m1h", 0, false},      // Out of order units
		{"1h2h", 0, false},      // Duplicate hours
		{"30m40m", 0, false},    // Duplicate minutes
		{"30s3s", 0, false},     // Duplicate seconds
		{"1h2m3s45", 0, false},  // Trailing numbers
		{"1h 2m", 0, false},     // Spaces within timestamp
	}

	for _, c := range cases {
		got, err := ParseYoutubeTime(c.input)
		if c.valid && err != nil {
			t.Errorf("ParseYoutubeTime(%q) unexpected error: %v", c.input, err)
		}
		if !c.valid && err == nil {
			t.Errorf("ParseYoutubeTime(%q) expected error, got %d", c.input, got)
		}
		if got != c.want {
			t.Errorf("ParseYoutubeTime(%q) = %d, want %d", c.input, got, c.want)
		}
	}
}
