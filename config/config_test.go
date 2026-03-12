package config

import (
	"os"
	"testing"
)

func TestGetEnvAsInt(t *testing.T) {
	tests := []struct {
		name         string
		envKey       string
		envValue     string
		defaultValue int
		expected     int
	}{
		{
			name:         "valid positive integer",
			envKey:       "TEST_PORT",
			envValue:     "5432",
			defaultValue: 3000,
			expected:     5432,
		},
		{
			name:         "valid zero value",
			envKey:       "TEST_ZERO",
			envValue:     "0",
			defaultValue: 3000,
			expected:     0,
		},
		{
			name:         "invalid non-numeric value",
			envKey:       "TEST_INVALID",
			envValue:     "not_a_number",
			defaultValue: 3000,
			expected:     3000,
		},
		{
			name:         "missing environment variable",
			envKey:       "TEST_MISSING_VAR",
			envValue:     "",
			defaultValue: 3000,
			expected:     3000,
		},
		{
			name:         "negative integer",
			envKey:       "TEST_NEGATIVE",
			envValue:     "-1",
			defaultValue: 3000,
			expected:     -1,
		},
		{
			name:         "large integer",
			envKey:       "TEST_LARGE",
			envValue:     "999999",
			defaultValue: 3000,
			expected:     999999,
		},
		{
			name:         "integer with spaces",
			envKey:       "TEST_SPACES",
			envValue:     "  999999  ",
			defaultValue: 3000,
			expected:     999999,
		},
		{
			name:         "blank string value",
			envKey:       "TEST_BLANK",
			envValue:     "   ",
			defaultValue: 3000,
			expected:     3000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv(tt.envKey, tt.envValue)
			defer os.Unsetenv(tt.envKey)

			result := getEnvAsInt(tt.envKey, tt.defaultValue)

			if result != tt.expected {
				t.Errorf("getEnvAsInt(%q) = %d, want %d", tt.envKey, result, tt.expected)
			}
		})
	}
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		envKey       string
		envValue     string
		defaultValue string
		expected     string
	}{
		{
			name:         "valid value",
			envKey:       "TEST_PORT",
			envValue:     "localhost",
			defaultValue: "default_host",
			expected:     "localhost",
		},
		{
			name:         "missing environment variable",
			envKey:       "TEST_MISSING_VAR",
			envValue:     "",
			defaultValue: "default_host",
			expected:     "default_host",
		},
		{
			name:         "value with spaces",
			envKey:       "TEST_SPACES",
			envValue:     "   ",
			defaultValue: "default_host",
			expected:     "default_host",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv(tt.envKey, tt.envValue)
			defer os.Unsetenv(tt.envKey)

			result := getEnv(tt.envKey, tt.defaultValue)

			if result != tt.expected {
				t.Errorf("getEnv(%q) = %q, want %q", tt.envKey, result, tt.expected)
			}
		})
	}
}
