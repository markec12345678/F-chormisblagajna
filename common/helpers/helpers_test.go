package helpers

import (
	"fmt"
	"os"
	"testing"
)

func TestRandStringBytesMaskImprSrc(t *testing.T) {
	tests := []struct {
		name string
		n    int
	}{
		{"short", 8},
		{"medium", 20},
		{"long", 64},
		{"odd", 15},
		{"one", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := RandStringBytesMaskImprSrc(tt.n)
			if err != nil {
				t.Fatalf("RandStringBytesMaskImprSrc(%d) returned error: %v", tt.n, err)
			}
			if len(result) != tt.n {
				t.Errorf("RandStringBytesMaskImprSrc(%d) returned length %d, want %d", tt.n, len(result), tt.n)
			}
		})
	}

	t.Run("uniqueness", func(t *testing.T) {
		a, _ := RandStringBytesMaskImprSrc(32)
		b, _ := RandStringBytesMaskImprSrc(32)
		if a == b {
			t.Errorf("Two random strings should differ, both got %q", a)
		}
	})

	t.Run("hex_chars_only", func(t *testing.T) {
		result, _ := RandStringBytesMaskImprSrc(40)
		for i, c := range result {
			if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
				t.Errorf("character at index %d is %q, expected hex char", i, c)
			}
		}
	})
}

func TestRandStringBytesMaskImprSrc_ValidLength(t *testing.T) {
	lengths := []int{1, 4, 7, 8, 15, 16, 32, 63, 64, 128, 256}
	for _, n := range lengths {
		t.Run(fmt.Sprintf("len_%d", n), func(t *testing.T) {
			result, err := RandStringBytesMaskImprSrc(n)
			if err != nil {
				t.Fatalf("RandStringBytesMaskImprSrc(%d) returned error: %v", n, err)
			}
			if len(result) != n {
				t.Errorf("RandStringBytesMaskImprSrc(%d) returned string of length %d, want %d", n, len(result), n)
			}
		})
	}
}

func TestRandStringBytesMaskImprSrc_SpecialChars(t *testing.T) {
	for i := 0; i < 200; i++ {
		result, err := RandStringBytesMaskImprSrc(64)
		if err != nil {
			t.Fatalf("iteration %d: unexpected error: %v", i, err)
		}
		for j, c := range result {
			if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
				t.Errorf("iteration %d, index %d: got char %q, expected only hex chars [0-9a-f]", i, j, c)
			}
		}
	}
}

func TestResolveOsEnvPath(t *testing.T) {
	t.Run("env_var_expansion", func(t *testing.T) {
		os.Setenv("TEST_HELPER_PATH", "/tmp/test")
		defer os.Unsetenv("TEST_HELPER_PATH")

		result := ResolveOsEnvPath("$TEST_HELPER_PATH")
		if result == "" {
			t.Error("ResolveOsEnvPath should expand env vars")
		}
	})

	t.Run("plain_path", func(t *testing.T) {
		result := ResolveOsEnvPath("/usr/local/bin")
		if result == "" {
			t.Error("ResolveOsEnvPath should handle plain paths")
		}
	})
}
