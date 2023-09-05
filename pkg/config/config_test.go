package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	tests := []struct {
		name        string
		envKey      string
		envValue    string
		defaultVal  interface{}
		expectedVal interface{}
	}{
		{
			name:        "Test GetStringVal with existing environment variable",
			envKey:      "EXISTING_KEY",
			envValue:    "existing_value",
			defaultVal:  "default_value",
			expectedVal: "existing_value",
		},
		{
			name:        "Test GetStringVal with missing environment variable",
			envKey:      "MISSING_KEY",
			envValue:    "",
			defaultVal:  "default_value",
			expectedVal: "default_value",
		},
		{
			name:        "Test GetIntValue with existing environment variable",
			envKey:      "EXISTING_INT_KEY",
			envValue:    "42",
			defaultVal:  24,
			expectedVal: 42,
		},
		{
			name:        "Test GetIntValue with invalid environment variable",
			envKey:      "INVALID_INT_KEY",
			envValue:    "not_an_integer",
			defaultVal:  24,
			expectedVal: 24,
		},
		{
			name:        "Test GetIntValue with missing environment variable",
			envKey:      "MISSING_INT_KEY",
			envValue:    "",
			defaultVal:  24,
			expectedVal: 24,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Set the environment variable for the tc.
			os.Setenv(tc.envKey, tc.envValue)
			defer os.Unsetenv(tc.envKey)

			if tc.defaultVal == "default_value" {
				actual := getStringVal(tc.envKey, tc.defaultVal.(string))
				assert.Equal(t, tc.expectedVal.(string), actual)
			} else {
				actual := getIntValue(tc.envKey, tc.defaultVal.(int))
				assert.Equal(t, tc.expectedVal.(int), actual)
			}
		})
	}
}

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name            string
		envKeyWebPage   string
		envValWebPage   string
		envKeyWorkers   string
		envValWorkers   string
		expectedWebPage string
		expectedWorkers int
	}{
		{
			name:            "Valid environment variables",
			envKeyWebPage:   "WEB_PAGE_URL",
			envValWebPage:   "https://example.com",
			envKeyWorkers:   "WORKER_COUNT",
			envValWorkers:   "5",
			expectedWebPage: "https://example.com",
			expectedWorkers: 5,
		},
		{
			name:            "Missing environment variables",
			envKeyWebPage:   "WEB_PAGE_URL",
			envValWebPage:   "",
			envKeyWorkers:   "WORKER_COUNT",
			envValWorkers:   "",
			expectedWebPage: defaultWebPage,
			expectedWorkers: defaultWorkers,
		},
		{
			name:            "Invalid worker count",
			envKeyWebPage:   "WEB_PAGE_URL",
			envValWebPage:   "https://example.com",
			envKeyWorkers:   "WORKER_COUNT",
			envValWorkers:   "invalid",
			expectedWebPage: "https://example.com",
			expectedWorkers: defaultWorkers,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Set environment variables for the test.
			os.Setenv(keyWebPage, test.envValWebPage)
			os.Setenv(keyWorkers, test.envValWorkers)
			defer os.Unsetenv(keyWebPage)
			defer os.Unsetenv(keyWorkers)

			config := NewConfig().(*config)
			assert.Equal(t, test.expectedWebPage, config.cfg.WepPage)
			assert.Equal(t, test.expectedWorkers, config.cfg.Workers)
		})
	}
}

func TestGetConfig(t *testing.T) {
	tests := []struct {
		name            string
		envKeyWebPage   string
		envValWebPage   string
		envKeyWorkers   string
		envValWorkers   string
		expectedWebPage string
		expectedWorkers int
	}{
		{
			name:            "Valid environment variables",
			envKeyWebPage:   "WEB_PAGE_URL",
			envValWebPage:   "https://example.com",
			envKeyWorkers:   "WORKER_COUNT",
			envValWorkers:   "5",
			expectedWebPage: "https://example.com",
			expectedWorkers: 5,
		},
		{
			name:            "Missing environment variables",
			envKeyWebPage:   "WEB_PAGE_URL",
			envValWebPage:   "",
			envKeyWorkers:   "WORKER_COUNT",
			envValWorkers:   "",
			expectedWebPage: defaultWebPage,
			expectedWorkers: defaultWorkers,
		},
		{
			name:            "Invalid worker count",
			envKeyWebPage:   "WEB_PAGE_URL",
			envValWebPage:   "https://example.com",
			envKeyWorkers:   "WORKER_COUNT",
			envValWorkers:   "invalid",
			expectedWebPage: "https://example.com",
			expectedWorkers: defaultWorkers,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Set environment variables for the test.
			os.Setenv(keyWebPage, test.envValWebPage)
			os.Setenv(keyWorkers, test.envValWorkers)
			defer os.Unsetenv(keyWebPage)
			defer os.Unsetenv(keyWorkers)

			config := NewConfig()
			cfg := config.GetConfig()
			assert.Equal(t, test.expectedWebPage, cfg.WepPage)
			assert.Equal(t, test.expectedWorkers, cfg.Workers)
		})
	}
}
