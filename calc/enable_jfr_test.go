package calc

import (
	"os"
	"strconv"
	"testing"
)

func TestNewEnableJfr_NoEnvVar(t *testing.T) {
	os.Unsetenv(EnvEnableJfr)
	defer os.Unsetenv(EnvEnableJfr)

	Jfr := NewEnableJfr()
	if *Jfr != DefaultEnableJfr {
		t.Errorf("Expected default value %v, got %v", DefaultEnableJfr, *Jfr)
	}
}

func TestNewEnableJfr_EnvVarTrue(t *testing.T) {
	os.Setenv(EnvEnableJfr, "true")
	defer os.Unsetenv(EnvEnableJfr)

	Jfr := NewEnableJfr()
	if *Jfr != EnableJfr(true) {
		t.Errorf("Expected true, got %v", *Jfr)
	}
}

func TestNewEnableJfr_EnvVarFalse(t *testing.T) {
	os.Setenv(EnvEnableJfr, "false")
	defer os.Unsetenv(EnvEnableJfr)

	Jfr := NewEnableJfr()
	if *Jfr != EnableJfr(false) {
		t.Errorf("Expected false, got %v", *Jfr)
	}
}

func TestEnableJfr_ContributeTrue(t *testing.T) {
	Jfr := EnableJfr(true)
	err := Jfr.Contribute()
	defer os.Unsetenv(EnvEnableJfr)
	if err != nil {
		t.Fatalf("Contribute returned an error: %v", err)
	}

	verifyJfrEnvVar(t, EnvEnableJfr, "true")
	verifyJfrEnvVar(t, envBplDebugPort, strconv.Itoa(defaultDebugPort))
}

func TestEnableJfr_ContributeFalse(t *testing.T) {
	Jfr := EnableJfr(false)
	err := Jfr.Contribute()
	defer os.Unsetenv(EnvEnableJfr)
	if err != nil {
		t.Fatalf("Contribute returned an error: %v", err)
	}

	verifyJfrEnvVar(t, EnvEnableJfr, "false")
	verifyJfrEnvVar(t, envBplDebugPort, strconv.Itoa(defaultDebugPort))
}

func verifyJfrEnvVar(t *testing.T, envVar, expectedValue string) {
	value, exists := os.LookupEnv(envVar)
	if !exists {
		t.Fatalf("Environment variable %s not set", envVar)
	}
	if value != expectedValue {
		t.Errorf("Expected environment variable %s to be '%s', got '%s'", envVar, expectedValue, value)
	}
}
