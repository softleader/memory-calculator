package boot_test

import (
	"os"
	"testing"

	"github.com/softleader/memory-calculator/boot"
)

func TestNewAppLibPath_NoEnvVar(t *testing.T) {
	os.Unsetenv(boot.EnvAppLibPath)

	alp := boot.NewAppLibPath()
	if *alp != boot.DefaultAppLibPath {
		t.Errorf("Expected default value '%s', got '%s'", boot.DefaultAppLibPath, *alp)
	}
}

func TestNewAppLibPath_EnvVarSet(t *testing.T) {
	testValue := "custom/libs"
	os.Setenv(boot.EnvAppLibPath, testValue)
	defer os.Unsetenv(boot.EnvAppLibPath)

	alp := boot.NewAppLibPath()
	if *alp != boot.AppLibPath(testValue) {
		t.Errorf("Expected value '%s', got '%s'", testValue, *alp)
	}
}

func TestAppLibPath_Contribute(t *testing.T) {
	testValue := "test/libs"
	alp := boot.AppLibPath(testValue)
	err := alp.Contribute()
	defer os.Unsetenv(boot.EnvAppLibPath)
	if err != nil {
		t.Errorf("Contribute returned an error: %v", err)
	}

	envValue, exists := os.LookupEnv(boot.EnvAppLibPath)
	if !exists {
		t.Fatalf("Environment variable %s not set", boot.EnvAppLibPath)
	}
	if envValue != testValue {
		t.Errorf("Expected environment variable %s to be '%s', got '%s'", boot.EnvAppLibPath, testValue, envValue)
	}
}
