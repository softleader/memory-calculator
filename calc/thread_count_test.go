package calc

import (
	"os"
	"strconv"
	"testing"
)

func TestNewThreadCount_NoEnvVar(t *testing.T) {
	os.Unsetenv(EnvThreadCount)
	defer os.Unsetenv(EnvThreadCount)

	tc := NewThreadCount()
	if *tc != DefaultThreadCount {
		t.Errorf("Expected default value %v, got %v", DefaultThreadCount, *tc)
	}
}

func TestNewThreadCount_EnvVarSet(t *testing.T) {
	testValue := "250"
	os.Setenv(EnvThreadCount, testValue)
	defer os.Unsetenv(EnvThreadCount)

	tc := NewThreadCount()
	expectedValue, _ := strconv.Atoi(testValue)
	if int(*tc) != expectedValue {
		t.Errorf("Expected value %v, got %v", expectedValue, *tc)
	}
}

func TestThreadCount_Contribute(t *testing.T) {
	testValue := 250
	tc := ThreadCount(testValue)
	err := tc.Contribute()
	defer os.Unsetenv(EnvThreadCount)
	if err != nil {
		t.Fatalf("Contribute returned an error: %v", err)
	}

	envValue, exists := os.LookupEnv(EnvThreadCount)
	if !exists {
		t.Fatalf("Environment variable %s not set", EnvThreadCount)
	}
	expectedEnvValue := strconv.Itoa(testValue)
	if envValue != expectedEnvValue {
		t.Errorf("Expected environment variable %s to be '%s', got '%s'", EnvThreadCount, expectedEnvValue, envValue)
	}
}
