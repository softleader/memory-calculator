package calc

import (
	"os"
	"testing"
)

func TestBuildJavaToolOptions_NoEnvVar(t *testing.T) {
	os.Unsetenv(EnvJavaToolOptions)
	j := BuildJavaToolOptions()
	expected := "-XX:+ExitOnOutOfMemoryError"
	if j.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, j.String())
	}
}

func TestBuildJavaToolOptions_EnvVarWithoutContributeOptions(t *testing.T) {
	testOption := "-Xmx1G"
	os.Setenv(EnvJavaToolOptions, testOption)
	defer os.Unsetenv(EnvJavaToolOptions)

	j := BuildJavaToolOptions()
	expected := testOption + " " + "-XX:+ExitOnOutOfMemoryError"
	if j.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, j.String())
	}
}

func TestBuildJavaToolOptions_EnvVarWithContributeOptions(t *testing.T) {
	os.Setenv(EnvJavaToolOptions, "-XX:+ExitOnOutOfMemoryError")
	defer os.Unsetenv(EnvJavaToolOptions)

	j := BuildJavaToolOptions()
	expected := "-XX:+ExitOnOutOfMemoryError"
	if j.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, j.String())
	}
}
