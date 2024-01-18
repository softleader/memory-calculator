package calc

import (
	"os"
	"testing"
)

func TestNewVerbose_NoEnvVar(t *testing.T) {
	os.Unsetenv(EnvVerbose)
	defer os.Unsetenv(EnvVerbose)

	v := NewVerbose()
	if *v != DefaultVerbose {
		t.Errorf("Expected default value %v, got %v", DefaultVerbose, *v)
	}
}

func TestNewVerbose_EnvVarSetDebug(t *testing.T) {
	os.Setenv(EnvVerbose, levelDebug)
	defer os.Unsetenv(EnvVerbose)

	v := NewVerbose()
	if *v {
		t.Errorf("Expected true, got %v", *v)
	}
}

func TestNewVerbose_EnvVarSetOther(t *testing.T) {
	os.Setenv(EnvVerbose, "INFO")
	defer os.Unsetenv(EnvVerbose)

	v := NewVerbose()
	if *v {
		t.Errorf("Expected false, got %v", *v)
	}
}

func TestVerbose_ContributeTrue(t *testing.T) {
	v := Verbose(true)
	err := v.Contribute()
	defer os.Unsetenv(EnvVerbose)
	if err != nil {
		t.Fatalf("Contribute returned an error: %v", err)
	}

	envValue, exists := os.LookupEnv(EnvVerbose)
	if !exists || envValue != levelDebug {
		t.Errorf("Expected environment variable %s to be set to '%s', got '%s'", EnvVerbose, levelDebug, envValue)
	}
}

func TestVerbose_ContributeFalse(t *testing.T) {
	v := Verbose(false)
	err := v.Contribute()
	defer os.Unsetenv(EnvVerbose)
	if err != nil {
		t.Fatalf("Contribute returned an error: %v", err)
	}

	_, exists := os.LookupEnv(EnvVerbose)
	if exists {
		t.Errorf("Environment variable %s should not be set", EnvVerbose)
	}
}
