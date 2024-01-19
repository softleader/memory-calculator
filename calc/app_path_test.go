package calc

import (
	"os"
	"testing"
)

func TestNewAppPath_NoEnvVar(t *testing.T) {
	os.Unsetenv(EnvAppPath)
	defer os.Unsetenv(EnvAppPath)

	j := NewAppPath()
	if *j != DefaultAppPath {
		t.Errorf("Expected default value '%s', got '%s'", DefaultAppPath, *j)
	}
}

func TestNewAppPath_EnvVarSet(t *testing.T) {
	testValue := "/test/path"
	os.Setenv(EnvAppPath, testValue)
	defer os.Unsetenv(EnvAppPath)

	j := NewAppPath()
	if *j != AppPath(testValue) {
		t.Errorf("Expected value '%s', got '%s'", testValue, *j)
	}
}

func TestAppPath_Contribute(t *testing.T) {
	testValue := "/test/path"
	appPath := AppPath(testValue)
	err := appPath.Contribute()
	defer os.Unsetenv(EnvAppPath)
	if err != nil {
		t.Errorf("Contribute returned an error: %v", err)
	}

	envValue, exists := os.LookupEnv(EnvAppPath)
	if !exists {
		t.Fatalf("Environment variable %s not set", EnvAppPath)
	}
	if envValue != testValue {
		t.Errorf("Expected environment variable %s to be '%s', got '%s'", EnvAppPath, testValue, envValue)
	}
}
