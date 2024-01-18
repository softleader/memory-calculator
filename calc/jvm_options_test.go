package calc

import (
	"os"
	"testing"
)

func TestNewJVMOptions_NoEnvVar(t *testing.T) {
	os.Unsetenv(EnvJVMOptions)
	defer os.Unsetenv(EnvJVMOptions)

	j := NewJVMOptions()
	if *j != DefaultJVMOptions {
		t.Errorf("Expected default value '%s', got '%s'", DefaultJVMOptions, *j)
	}
}

func TestNewJVMOptions_EnvVarSet(t *testing.T) {
	testValue := "-Xmx1G"
	os.Setenv(EnvJVMOptions, testValue)
	defer os.Unsetenv(EnvJVMOptions)

	j := NewJVMOptions()
	if *j != JVMOptions(testValue) {
		t.Errorf("Expected value '%s', got '%s'", testValue, *j)
	}
}

func TestJVMOptions_Contribute_NonEmpty(t *testing.T) {
	testValue := "-Xmx1G"
	j := JVMOptions(testValue)
	err := j.Contribute()
	defer os.Unsetenv(EnvJVMOptions)
	if err != nil {
		t.Fatalf("Contribute returned an error: %v", err)
	}

	envValue, exists := os.LookupEnv(EnvJVMOptions)
	if !exists {
		t.Fatalf("Environment variable %s not set", EnvJVMOptions)
	}
	if envValue != testValue {
		t.Errorf("Expected environment variable %s to be '%s', got '%s'", EnvJVMOptions, testValue, envValue)
	}
}

func TestJVMOptions_Contribute_Empty(t *testing.T) {
	j := JVMOptions("")
	err := j.Contribute()
	defer os.Unsetenv(EnvJVMOptions)
	if err != nil {
		t.Fatalf("Contribute returned an error: %v", err)
	}

	_, exists := os.LookupEnv(EnvJVMOptions)
	if exists {
		t.Errorf("Environment variable %s should not be set", EnvJVMOptions)
	}
}
