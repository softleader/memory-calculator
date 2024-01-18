package calc

import (
	"os"
	"strconv"
	"testing"
)

func TestNewEnableJdwp_NoEnvVar(t *testing.T) {
	os.Unsetenv(EnvEnableJdwp)
	defer os.Unsetenv(EnvEnableJdwp)

	jdwp := NewEnableJdwp()
	if *jdwp != DefaultEnableJdwp {
		t.Errorf("Expected default value %v, got %v", DefaultEnableJdwp, *jdwp)
	}
}

func TestNewEnableJdwp_EnvVarTrue(t *testing.T) {
	os.Setenv(EnvEnableJdwp, "true")
	defer os.Unsetenv(EnvEnableJdwp)

	jdwp := NewEnableJdwp()
	if *jdwp != EnableJdwp(true) {
		t.Errorf("Expected true, got %v", *jdwp)
	}
}

func TestNewEnableJdwp_EnvVarFalse(t *testing.T) {
	os.Setenv(EnvEnableJdwp, "false")
	defer os.Unsetenv(EnvEnableJdwp)

	jdwp := NewEnableJdwp()
	if *jdwp != EnableJdwp(false) {
		t.Errorf("Expected false, got %v", *jdwp)
	}
}

func TestEnableJdwp_ContributeTrue(t *testing.T) {
	jdwp := EnableJdwp(true)
	err := jdwp.Contribute()
	defer os.Unsetenv(EnvEnableJdwp)
	if err != nil {
		t.Fatalf("Contribute returned an error: %v", err)
	}

	verifyJdwpEnvVar(t, EnvEnableJdwp, "true")
	verifyJdwpEnvVar(t, envBplDebugPort, strconv.Itoa(defaultDebugPort))
}

func TestEnableJdwp_ContributeFalse(t *testing.T) {
	jdwp := EnableJdwp(false)
	err := jdwp.Contribute()
	defer os.Unsetenv(EnvEnableJdwp)
	if err != nil {
		t.Fatalf("Contribute returned an error: %v", err)
	}

	verifyJdwpEnvVar(t, EnvEnableJdwp, "false")
	verifyJdwpEnvVar(t, envBplDebugPort, strconv.Itoa(defaultDebugPort))
}

func verifyJdwpEnvVar(t *testing.T, envVar, expectedValue string) {
	value, exists := os.LookupEnv(envVar)
	if !exists {
		t.Fatalf("Environment variable %s not set", envVar)
	}
	if value != expectedValue {
		t.Errorf("Expected environment variable %s to be '%s', got '%s'", envVar, expectedValue, value)
	}
}
