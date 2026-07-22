package helpers

import (
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
			result := RandStringBytesMaskImprSrc(tt.n)
			if len(result) != tt.n {
				t.Errorf("RandStringBytesMaskImprSrc(%d) returned length %d, want %d", tt.n, len(result), tt.n)
			}
		})
	}

	t.Run("uniqueness", func(t *testing.T) {
		a := RandStringBytesMaskImprSrc(32)
		b := RandStringBytesMaskImprSrc(32)
		if a == b {
			t.Errorf("Two random strings should differ, both got %q", a)
		}
	})

	t.Run("hex_chars_only", func(t *testing.T) {
		result := RandStringBytesMaskImprSrc(40)
		for i, c := range result {
			if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
				t.Errorf("character at index %d is %q, expected hex char", i, c)
			}
		}
	})
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
