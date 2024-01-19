package calc

import (
	"os"
	"testing"
)

func TestNewJVMClassAdj_NoEnvVar(t *testing.T) {
	os.Unsetenv(EnvJVMClassAdj)
	jca := NewJVMClassAdj()
	if *jca != DefaultJVMClassAdj {
		t.Errorf("Expected default JVMClassAdj, got %v", *jca)
	}
}

func TestNewJVMClassAdj_WithEnvVar(t *testing.T) {
	testVal := "10"
	os.Setenv(EnvJVMClassAdj, testVal)
	defer os.Unsetenv(EnvJVMClassAdj)

	jca := NewJVMClassAdj()
	if jca.String() != testVal {
		t.Errorf("Expected JVMClassAdj %v, got %v", testVal, jca.String())
	}
}

func TestContribute_PositiveValue(t *testing.T) {
	jca := JVMClassAdj("5")
	err := jca.Contribute()
	defer os.Unsetenv(EnvJVMClassAdj)
	if err != nil {
		t.Errorf("Contribute should not return error for positive JVMClassAdj")
	}

	val, ok := os.LookupEnv(EnvJVMClassAdj)
	if !ok || val != jca.String() {
		t.Errorf("Expected environment variable %v to be set to %v", EnvJVMClassAdj, jca.String())
	}
}

func TestContribute_NonPositiveValue(t *testing.T) {
	jca := JVMClassAdj("")
	err := jca.Contribute()
	defer os.Unsetenv(EnvJVMClassAdj)
	if err != nil {
		t.Errorf("Contribute should not return error for non-positive JVMClassAdj")
	}

	_, ok := os.LookupEnv(EnvJVMClassAdj)
	if ok {
		t.Errorf("Environment variable %v should not be set for non-positive JVMClassAdj", EnvJVMClassAdj)
	}
}
