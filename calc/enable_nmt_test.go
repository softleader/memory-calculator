package calc

import (
	"os"
	"strconv"
	"testing"
)

func TestNewEnableNmt_NoEnvVar(t *testing.T) {
	os.Unsetenv(EnvEnableNmt)
	defer os.Unsetenv(EnvEnableNmt)

	Nmt := NewEnableNmt()
	if *Nmt != DefaultEnableNmt {
		t.Errorf("Expected default value %v, got %v", DefaultEnableNmt, *Nmt)
	}
}

func TestNewEnableNmt_EnvVarTrue(t *testing.T) {
	os.Setenv(EnvEnableNmt, "true")
	defer os.Unsetenv(EnvEnableNmt)

	Nmt := NewEnableNmt()
	if *Nmt != EnableNmt(true) {
		t.Errorf("Expected true, got %v", *Nmt)
	}
}

func TestNewEnableNmt_EnvVarFalse(t *testing.T) {
	os.Setenv(EnvEnableNmt, "false")
	defer os.Unsetenv(EnvEnableNmt)

	Nmt := NewEnableNmt()
	if *Nmt != EnableNmt(false) {
		t.Errorf("Expected false, got %v", *Nmt)
	}
}

func TestEnableNmt_ContributeTrue(t *testing.T) {
	Nmt := EnableNmt(true)
	err := Nmt.Contribute()
	defer os.Unsetenv(EnvEnableNmt)
	if err != nil {
		t.Fatalf("Contribute returned an error: %v", err)
	}

	verifyNmtEnvVar(t, EnvEnableNmt, "true")
	verifyNmtEnvVar(t, envBplDebugPort, strconv.Itoa(defaultDebugPort))
}

func TestEnableNmt_ContributeFalse(t *testing.T) {
	Nmt := EnableNmt(false)
	err := Nmt.Contribute()
	defer os.Unsetenv(EnvEnableNmt)
	if err != nil {
		t.Fatalf("Contribute returned an error: %v", err)
	}

	verifyNmtEnvVar(t, EnvEnableNmt, "false")
	verifyNmtEnvVar(t, envBplDebugPort, strconv.Itoa(defaultDebugPort))
}

func verifyNmtEnvVar(t *testing.T, envVar, expectedValue string) {
	value, exists := os.LookupEnv(envVar)
	if !exists {
		t.Fatalf("Environment variable %s not set", envVar)
	}
	if value != expectedValue {
		t.Errorf("Expected environment variable %s to be '%s', got '%s'", envVar, expectedValue, value)
	}
}
