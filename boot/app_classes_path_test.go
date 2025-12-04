package boot_test

import (
	"os"
	"testing"

	"github.com/softleader/memory-calculator/boot"
)

func TestNewAppClassesPath_NoEnvVar(t *testing.T) {
	os.Unsetenv(boot.EnvAppClassesPath)

	acp := boot.NewAppClassesPath()
	if *acp != boot.DefaultAppClassesPath {
		t.Errorf("Expected default value '%s', got '%s'", boot.DefaultAppClassesPath, *acp)
	}
}

func TestNewAppClassesPath_EnvVarSet(t *testing.T) {
	testValue := "custom/path"
	os.Setenv(boot.EnvAppClassesPath, testValue)
	defer os.Unsetenv(boot.EnvAppClassesPath)

	acp := boot.NewAppClassesPath()
	if *acp != boot.AppClassesPath(testValue) {
		t.Errorf("Expected value '%s', got '%s'", testValue, *acp)
	}
}

func TestAppClassesPath_Contribute(t *testing.T) {
	testValue := "test/path"
	acp := boot.AppClassesPath(testValue)
	err := acp.Contribute()
	defer os.Unsetenv(boot.EnvAppClassesPath)
	if err != nil {
		t.Errorf("Contribute returned an error: %v", err)
	}

	envValue, exists := os.LookupEnv(boot.EnvAppClassesPath)
	if !exists {
		t.Fatalf("Environment variable %s not set", boot.EnvAppClassesPath)
	}
	if envValue != testValue {
		t.Errorf("Expected environment variable %s to be '%s', got '%s'", boot.EnvAppClassesPath, testValue, envValue)
	}
}
