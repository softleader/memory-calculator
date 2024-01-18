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

func TestLoadedClassCount_Contribute_ZeroWithoutJavaHome(t *testing.T) {
	os.Unsetenv(envJavaHome)
	defer os.Unsetenv(envJavaHome)

	lcc := LoadedClassCount(0)
	err := lcc.Contribute()
	if err == nil {
		t.Fatalf("Contribute should return an error when JAVA_HOME is not set")
	}
}

func TestLoadedClassCount_Contribute_NonZero(t *testing.T) {
	testValue := 100
	lcc := LoadedClassCount(testValue)
	err := lcc.Contribute()
	defer os.Unsetenv(EnvLoadedClassCount)
	if err != nil {
		t.Fatalf("Contribute returned an error: %v", err)
	}
}
