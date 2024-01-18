package calc

import (
	"os"
	"testing"
)

func TestAppPath_Set(t *testing.T) {
	appPath := NewAppPath()
	testValue := "/test/path"
	err := appPath.Set(testValue)
	if err != nil {
		t.Errorf("Set returned an error: %v", err)
	}
	if appPath.String() != testValue {
		t.Errorf("Expected AppPath value '%s', got '%s'", testValue, appPath.String())
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
