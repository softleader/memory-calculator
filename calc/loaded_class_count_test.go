package calc

import (
	"os"
	"strconv"
	"testing"
)

func TestNewLoadedClassCount_NoEnvVar(t *testing.T) {
	os.Unsetenv(EnvLoadedClassCount)
	defer os.Unsetenv(EnvLoadedClassCount)

	lcc := NewLoadedClassCount()
	if *lcc != DefaultLoadedClassCount {
		t.Errorf("Expected default value %v, got %v", DefaultLoadedClassCount, *lcc)
	}
}

func TestNewLoadedClassCount_EnvVarSet(t *testing.T) {
	testValue := "100"
	os.Setenv(EnvLoadedClassCount, testValue)
	defer os.Unsetenv(EnvLoadedClassCount)

	lcc := NewLoadedClassCount()
	expectedValue, _ := strconv.Atoi(testValue)
	if int(*lcc) != expectedValue {
		t.Errorf("Expected value %v, got %v", expectedValue, *lcc)
	}
}

func TestLoadedClassCount_Contribute(t *testing.T) {
	testValue := 100
	lcc := LoadedClassCount(testValue)
	err := lcc.Contribute()
	defer os.Unsetenv(EnvLoadedClassCount)
	if err != nil {
		t.Fatalf("Contribute returned an error: %v", err)
	}

	envValue, exists := os.LookupEnv(EnvLoadedClassCount)
	if !exists {
		t.Fatalf("Environment variable %s not set", EnvLoadedClassCount)
	}
	expectedEnvValue := strconv.Itoa(testValue)
	if envValue != expectedEnvValue {
		t.Errorf("Expected environment variable %s to be '%s', got '%s'", EnvLoadedClassCount, expectedEnvValue, envValue)
	}
}
