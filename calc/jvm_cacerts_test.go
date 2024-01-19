package calc

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewJVMCacerts_NoEnvVar(t *testing.T) {
	os.Unsetenv(EnvJVMCacerts)
	j := NewJVMCacerts()
	if *j != DefaultJVMCacerts {
		t.Errorf("Expected default JVMCacerts, got %v", *j)
	}
}

func TestNewJVMCacerts_WithEnvVar(t *testing.T) {
	testVal := "/path/to/cacerts"
	os.Setenv(EnvJVMCacerts, testVal)
	defer os.Unsetenv(EnvJVMCacerts)

	j := NewJVMCacerts()
	if *j != JVMCacerts(testVal) {
		t.Errorf("Expected JVMCacerts %v, got %v", testVal, *j)
	}
}

func TestContribute_JVMCacertsEmptyAndNoJavaHome(t *testing.T) {
	j := NewJVMCacerts()
	os.Unsetenv(envJavaHome)

	err := j.Contribute()
	defer os.Unsetenv(EnvJVMCacerts)
	if err != nil {
		t.Errorf("Contribute returned error: %v", err)
	}

	val, ok := os.LookupEnv(EnvJVMCacerts)
	if ok {
		t.Errorf("Expected environment variable %v not to be set", EnvJVMCacerts)
		if val != j.String() {
			t.Errorf("Expected empty JVMCacerts, got %v", val)
		}
	}
}

func TestContribute_JVMCacertsEmptyAndJavaHomeSet(t *testing.T) {
	javaHomePath, err := os.MkdirTemp("", "fake-java-home")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(javaHomePath)
	cacert := filepath.Join(javaHomePath, subPathCacerts)
	if err = os.MkdirAll(filepath.Dir(cacert), 0755); err != nil {
		t.Fatalf("Failed to create cacert dir: %v", err)
	}
	if _, err = os.Create(cacert); err != nil {
		t.Fatalf("Failed to create cacert file: %v", err)
	}

	j := NewJVMCacerts()
	os.Setenv(envJavaHome, javaHomePath)
	defer os.Unsetenv(envJavaHome)

	err = j.Contribute()
	defer os.Unsetenv(EnvJVMCacerts)
	if err != nil {
		t.Errorf("Contribute returned error: %v", err)
	}

	val, ok := os.LookupEnv(EnvJVMCacerts)
	if !ok || val != cacert {
		t.Errorf("Expected JVMCacerts to be set to %v, got %v", cacert, val)
	}
}

func TestContribute_JVMCacertsNotEmpty(t *testing.T) {
	javaHomePath, err := os.MkdirTemp("", "fake-java-home")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(javaHomePath)
	cacert := filepath.Join(javaHomePath, subPathCacerts)
	if err = os.MkdirAll(filepath.Dir(cacert), 0755); err != nil {
		t.Fatalf("Failed to create cacert dir: %v", err)
	}
	if _, err = os.Create(cacert); err != nil {
		t.Fatalf("Failed to create cacert file: %v", err)
	}

	j := JVMCacerts(cacert)
	err = j.Contribute()
	defer os.Unsetenv(EnvJVMCacerts)
	if err != nil {
		t.Errorf("Contribute returned error: %v", err)
	}

	val, ok := os.LookupEnv(EnvJVMCacerts)
	if !ok || val != cacert {
		t.Errorf("Expected JVMCacerts to be set to %v, got %v", cacert, val)
	}
}
