package calc

import (
	"os"
	"strconv"
	"testing"
)

func TestNewEnableJmx_NoEnvVar(t *testing.T) {
	os.Unsetenv(EnvEnableJmx)
	defer os.Unsetenv(EnvEnableJmx)

	Jmx := NewEnableJmx()
	if *Jmx != DefaultEnableJmx {
		t.Errorf("Expected default value %v, got %v", DefaultEnableJmx, *Jmx)
	}
}

func TestNewEnableJmx_EnvVarTrue(t *testing.T) {
	os.Setenv(EnvEnableJmx, "true")
	defer os.Unsetenv(EnvEnableJmx)

	Jmx := NewEnableJmx()
	if *Jmx != EnableJmx(true) {
		t.Errorf("Expected true, got %v", *Jmx)
	}
}

func TestNewEnableJmx_EnvVarFalse(t *testing.T) {
	os.Setenv(EnvEnableJmx, "false")
	defer os.Unsetenv(EnvEnableJmx)

	Jmx := NewEnableJmx()
	if *Jmx != EnableJmx(false) {
		t.Errorf("Expected false, got %v", *Jmx)
	}
}

func TestEnableJmx_ContributeTrue(t *testing.T) {
	Jmx := EnableJmx(true)
	err := Jmx.Contribute()
	defer os.Unsetenv(EnvEnableJmx)
	if err != nil {
		t.Fatalf("Contribute returned an error: %v", err)
	}

	verifyJmxEnvVar(t, EnvEnableJmx, "true")
	verifyJmxEnvVar(t, envBplDebugPort, strconv.Itoa(defaultDebugPort))
}

func TestEnableJmx_ContributeFalse(t *testing.T) {
	Jmx := EnableJmx(false)
	err := Jmx.Contribute()
	defer os.Unsetenv(EnvEnableJmx)
	if err != nil {
		t.Fatalf("Contribute returned an error: %v", err)
	}

	verifyJmxEnvVar(t, EnvEnableJmx, "false")
	verifyJmxEnvVar(t, envBplDebugPort, strconv.Itoa(defaultDebugPort))
}

func verifyJmxEnvVar(t *testing.T, envVar, expectedValue string) {
	value, exists := os.LookupEnv(envVar)
	if !exists {
		t.Fatalf("Environment variable %s not set", envVar)
	}
	if value != expectedValue {
		t.Errorf("Expected environment variable %s to be '%s', got '%s'", envVar, expectedValue, value)
	}
}
