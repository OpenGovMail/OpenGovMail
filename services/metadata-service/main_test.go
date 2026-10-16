package main

import (
	"os"
	"testing"
)

func TestGetEnvAsBool(t *testing.T) {
	const key = "TEST_GET_ENV_AS_BOOL"

	tests := []struct {
		name  string
		value string
		set   bool
		def   bool
		want  bool
	}{
		{"unset uses default true", "", false, true, true},
		{"unset uses default false", "", false, false, false},
		{"empty uses default", "", true, true, true},
		{"whitespace only uses default", "   ", true, true, true},
		{"true", "true", true, false, true},
		{"True", "True", true, false, true},
		{"TRUE", "TRUE", true, false, true},
		{"1", "1", true, false, true},
		{"false", "false", true, true, false},
		{"False", "False", true, true, false},
		{"0", "0", true, true, false},
		{"whitespace trimmed", " false ", true, true, false},
		{"yes is invalid, uses default", "yes", true, true, true},
		{"garbage uses default", "banana", true, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.set {
				t.Setenv(key, tt.value)
			} else {
				os.Unsetenv(key)
			}
			if got := getEnvAsBool(key, tt.def); got != tt.want {
				t.Errorf("getEnvAsBool(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}