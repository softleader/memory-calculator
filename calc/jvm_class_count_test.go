package calc

import (
	"os"
	"testing"
)

func TestNewJVMClassCount_NoEnvVar(t *testing.T) {
	os.Unsetenv(EnvJVMClassCount)
	lcc := NewJVMClassCount()
	if *lcc != DefaultJVMClassCount {
		t.Errorf("Expected default JVMClassCount, got %v", *lcc)
	}
}

func TestNewJVMClassCount_WithEnvVar(t *testing.T) {
	testVal := "10"
	os.Setenv(EnvJVMClassCount, testVal)
	defer os.Unsetenv(EnvJVMClassCount)

	lcc := NewJVMClassCount()
	if lcc.String() != testVal {
		t.Errorf("Expected JVMClassCount %v, got %v", testVal, lcc.String())
	}
}

func TestContribute_NonZeroValue(t *testing.T) {
	lcc := JVMClassCount(5)
	err := lcc.Contribute()
	defer os.Unsetenv(EnvJVMClassCount)
	if err != nil {
		t.Errorf("Contribute should not return error for non-zero JVMClassCount")
	}
}

func TestContribute_ZeroValueNoJavaHome(t *testing.T) {
	lcc := JVMClassCount(0)
	os.Unsetenv(envJavaHome)

	err := lcc.Contribute()
	defer os.Unsetenv(EnvJVMClassCount)
	if err == nil {
		t.Errorf("Contribute should return error when JAVA_HOME is not set")
	}
}
