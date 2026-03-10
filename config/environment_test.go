package config

import (
	"os"
	"testing"
)

func TestGetEnvironmentVariableAsInt(t *testing.T) {
	tests := []struct {
		name     string
		envKey   string
		envValue string
		expected int
	}{
		{
			name:     "valid positive integer",
			envKey:   "TEST_PORT",
			envValue: "5432",
			expected: 5432,
		},
		{
			name:     "valid zero value",
			envKey:   "TEST_ZERO",
			envValue: "0",
			expected: 0,
		},
		{
			name:     "invalid non-numeric value",
			envKey:   "TEST_INVALID",
			envValue: "not_a_number",
			expected: 0,
		},
		{
			name:     "missing environment variable",
			envKey:   "TEST_MISSING_VAR",
			envValue: "",
			expected: 0,
		},
		{
			name:     "negative integer",
			envKey:   "TEST_NEGATIVE",
			envValue: "-1",
			expected: -1,
		},
		{
			name:     "large integer",
			envKey:   "TEST_LARGE",
			envValue: "999999",
			expected: 999999,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv(tt.envKey, tt.envValue)
			defer os.Unsetenv(tt.envKey)

			result := getEnvironmentVariableAsInt(tt.envKey)

			if result != tt.expected {
				t.Errorf("getEnvironmentVariableAsInt(%q) = %d, want %d", tt.envKey, result, tt.expected)
			}
		})
	}
}
