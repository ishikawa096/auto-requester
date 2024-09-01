package scheduler

import (
	"os"
	"testing"
)

func TestGetStrEnv(t *testing.T) {
	tests := []struct {
		name         string
		envKey       string
		envValue     string
		defaultValue string
		expected     string
	}{
		{
			name:         "Environment variable is valid",
			envKey:       "TEST_STR_ENV",
			envValue:     "test",
			defaultValue: "default",
			expected:     "test",
		},
		{
			name:         "Environment variable is empty",
			envKey:       "TEST_STR_ENV",
			envValue:     "",
			defaultValue: "default",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.envKey, tt.envValue)
			} else {
				os.Unsetenv(tt.envKey)
			}

			result := getStrEnv(tt.envKey, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("getStrEnv(%s, %s) = %s; want %s", tt.envKey, tt.defaultValue, result, tt.expected)
			}

			os.Unsetenv(tt.envKey)
		})
	}
}

func TestGetIntEnv(t *testing.T) {
	tests := []struct {
		name         string
		envKey       string
		envValue     string
		defaultValue int
		expected     int
	}{
		{
			name:         "Environment variable is valid",
			envKey:       "TEST_INT_ENV",
			envValue:     "100",
			defaultValue: 0,
			expected:     100,
		},
		{
			name:         "Environment variable is invalid",
			envKey:       "TEST_INT_ENV",
			envValue:     "invalid",
			defaultValue: 0,
			expected:     0,
		},
		{
			name:         "Environment variable is not set",
			envKey:       "TEST_INT_ENV",
			envValue:     "",
			defaultValue: 100,
			expected:     100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.envKey, tt.envValue)
			} else {
				os.Unsetenv(tt.envKey)
			}

			result := getIntEnv(tt.envKey, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("getIntEnv(%s, %v) = %v; want %v", tt.envKey, tt.defaultValue, result, tt.expected)
			}

			os.Unsetenv(tt.envKey)
		})
	}
}

func TestGetBoolEnv(t *testing.T) {
	tests := []struct {
		name         string
		envKey       string
		envValue     string
		defaultValue bool
		expected     bool
	}{
		{
			name:         "Environment variable is true",
			envKey:       "TEST_BOOL_ENV",
			envValue:     "true",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "Environment variable is false",
			envKey:       "TEST_BOOL_ENV",
			envValue:     "false",
			defaultValue: true,
			expected:     false,
		},
		{
			name:         "Environment variable is invalid",
			envKey:       "TEST_BOOL_ENV",
			envValue:     "invalid",
			defaultValue: true,
			expected:     true,
		},
		{
			name:         "Environment variable is not set",
			envKey:       "TEST_BOOL_ENV",
			envValue:     "",
			defaultValue: true,
			expected:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.envKey, tt.envValue)
			} else {
				os.Unsetenv(tt.envKey)
			}

			result := getBoolEnv(tt.envKey, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("getBoolEnv(%s, %v) = %v; want %v", tt.envKey, tt.defaultValue, result, tt.expected)
			}

			os.Unsetenv(tt.envKey)
		})
	}
}
