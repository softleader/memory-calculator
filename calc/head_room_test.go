package calc

import (
	"os"
	"strconv"
	"testing"
)

func TestNewHeadRoom_NoEnvVar(t *testing.T) {
	os.Unsetenv(EnvHeadRoom)
	defer os.Unsetenv(EnvHeadRoom)

	hr := NewHeadRoom()
	if *hr != DefaultHeadRoom {
		t.Errorf("Expected default value %v, got %v", DefaultHeadRoom, *hr)
	}
}

func TestNewHeadRoom_EnvVarSet(t *testing.T) {
	testValue := "20"
	os.Setenv(EnvHeadRoom, testValue)
	defer os.Unsetenv(EnvHeadRoom)

	hr := NewHeadRoom()
	expectedValue, _ := strconv.Atoi(testValue)
	if int(*hr) != expectedValue {
		t.Errorf("Expected value %v, got %v", expectedValue, *hr)
	}
}
func TestHeadRoom_Contribute(t *testing.T) {
	testValue := 25
	hr := HeadRoom(testValue)
	err := hr.Contribute()
	defer os.Unsetenv(EnvHeadRoom)
	if err != nil {
		t.Fatalf("Contribute returned an error: %v", err)
	}

	envValue, exists := os.LookupEnv(EnvHeadRoom)
	if !exists {
		t.Fatalf("Environment variable %s not set", EnvHeadRoom)
	}
	expectedEnvValue := strconv.Itoa(testValue)
	if envValue != expectedEnvValue {
		t.Errorf("Expected environment variable %s to be '%s', got '%s'", EnvHeadRoom, expectedEnvValue, envValue)
	}
}
